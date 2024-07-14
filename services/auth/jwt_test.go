package auth

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/loloDawit/ecom/types"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	// Define test cases
	tests := []struct {
		name       string
		secret     []byte
		userID     int
		expiration time.Duration
		wantErr    bool
	}{
		{
			name:       "Valid token generation",
			secret:     []byte("secret"),
			userID:     123,
			expiration: time.Second * 60,
			wantErr:    false,
		},
		{
			name:       "Empty secret",
			secret:     []byte(""),
			userID:     123,
			expiration: time.Second * 60,
			wantErr:    true,
		},
		{
			name:       "Negative userID",
			secret:     []byte("secret"),
			userID:     -1,
			expiration: time.Second * 60,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateToken(tt.secret, tt.userID, tt.expiration)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				token, _ := jwt.Parse(got, func(token *jwt.Token) (interface{}, error) {
					return tt.secret, nil
				})

				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					assert.Equal(t, strconv.Itoa(tt.userID), claims["userID"])
					assert.WithinDuration(t, time.Now().Add(tt.expiration), time.Unix(int64(claims["exp"].(float64)), 0), time.Second*5)
				} else {
					t.Errorf("Token is invalid or claims are not as expected")
				}
			}
		})
	}
}

func generateToken(secret []byte, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
	})
	return token.SignedString(secret)
}

func TestJWTMiddleware(t *testing.T) {
	secret := []byte("my_secret_key")

	tests := []struct {
		name              string
		authHeader        string
		bypassUserID      bool
		expectedStatus    int
		expectedResponse  string
		expectUserIDInCtx bool
	}{
		{
			name:              "Missing Authorization header",
			authHeader:        "",
			bypassUserID:      false,
			expectedStatus:    http.StatusUnauthorized,
			expectedResponse:  `{"error":"Authorization header is missing"}`,
			expectUserIDInCtx: false,
		},
		{
			name:              "Missing token",
			authHeader:        "Bearer ",
			bypassUserID:      false,
			expectedStatus:    http.StatusUnauthorized,
			expectedResponse:  `{"error":"Token is missing"}`,
			expectUserIDInCtx: false,
		},
		{
			name:              "Invalid token",
			authHeader:        "Bearer invalid_token",
			bypassUserID:      false,
			expectedStatus:    http.StatusUnauthorized,
			expectedResponse:  `{"error":"Invalid token"}`,
			expectUserIDInCtx: false,
		},
		{
			name: "Valid token",
			authHeader: func() string {
				token, _ := generateToken(secret, "12345")
				return "Bearer " + token
			}(),
			bypassUserID:      false,
			expectedStatus:    http.StatusOK,
			expectedResponse:  "",
			expectUserIDInCtx: true,
		},
		{
			name: "Bypass User ID",
			authHeader: func() string {
				token, _ := generateToken(secret, "12345")
				return "Bearer " + token
			}(),
			bypassUserID:      true,
			expectedStatus:    http.StatusOK,
			expectedResponse:  "",
			expectUserIDInCtx: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := JWTMiddleware(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID := r.Context().Value(types.UserIDKey)
				if tt.expectUserIDInCtx {
					assert.Equal(t, "12345", userID)
				} else {
					assert.Nil(t, userID)
				}
				w.WriteHeader(http.StatusOK)
			}))

			req, err := http.NewRequest("GET", "/", nil)
			assert.NoError(t, err)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			if tt.bypassUserID {
				req.Header.Set("X-Bypass-UserID", "true")
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedResponse != "" {
				assert.JSONEq(t, tt.expectedResponse, rr.Body.String())
			}
		})
	}
}