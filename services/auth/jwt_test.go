package auth

import (
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
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
