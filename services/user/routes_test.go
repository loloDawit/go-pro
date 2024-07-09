package user

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/loloDawit/ecom/config"
	"github.com/loloDawit/ecom/services/auth"
	"github.com/loloDawit/ecom/types"
	"github.com/loloDawit/ecom/utils"
	"gopkg.in/go-playground/validator.v9"
)

type mockUserStore struct {
	db *sql.DB
	GetUserByEmailFunc func(email string) (*types.User, error)
	CreateUserFunc     func(user types.User) error
}


func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(email)
	}
	return nil, sql.ErrNoRows
}

// CreateUser implements types.UserStore.
func (m *mockUserStore) CreateUser(user types.User) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(user)
	}
	return nil
}
// GetUserByID implements types.UserStore.
func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}


func TestCheckUserExists(t *testing.T) {
	tests := []struct {
		name          string
		mockBehavior  func(email string) (*types.User, error)
		expectedError error
	}{
		{
			name: "User exists",
			mockBehavior: func(email string) (*types.User, error) {
				return &types.User{}, nil
			},
			expectedError: fmt.Errorf(utils.ErrUserAlreadyExists),
		},
		{
			name: "User does not exist",
			mockBehavior: func(email string) (*types.User, error) {
				return nil, sql.ErrNoRows
			},
			expectedError: nil,
		},
		{
			name: "Database error",
			mockBehavior: func(email string) (*types.User, error) {
				return nil, fmt.Errorf("some database error")
			},
			expectedError: fmt.Errorf(utils.ErrInternalServerError),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStore := &mockUserStore{
				GetUserByEmailFunc: tc.mockBehavior,
			}
			handler := &Handler{store: mockStore}

			err := handler.checkUserExists("test@example.com")
			if (err != nil && tc.expectedError == nil) || (err == nil && tc.expectedError != nil) || (err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error()) {
				t.Fatalf("expected error %v, got %v", tc.expectedError, err)
			}
		})
	}
}

type mockValidator struct{}

func (v *mockValidator) Struct(s interface{}) error {
	switch payload := s.(type) {
	case types.SignupUserPayload:
		if payload.Email == "invalid-email" {
			return validator.ValidationErrors{}
		}
	case types.LoginUserPayload:
		if payload.Email == "invalid-email" {
			return validator.ValidationErrors{}
		}
	default:
		return fmt.Errorf("unsupported payload type")
	}
	return nil
}

func TestSignUp(t *testing.T) {
	originalValidate := utils.Validate
	utils.Validate = &mockValidator{}
	defer func() { utils.Validate = originalValidate }()

	tests := []struct {
		name             string
		payload          types.SignupUserPayload
		mockStore        *mockUserStore
		expectedStatus   int
		expectedResponse map[string]string
	}{
		{
			name: "Valid signup",
			payload: types.SignupUserPayload{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "password",
			},
			mockStore: &mockUserStore{
				GetUserByEmailFunc: func(email string) (*types.User, error) {
					return nil, sql.ErrNoRows
				},
				CreateUserFunc: func(user types.User) error {
					return nil
				},
			},
			expectedStatus:   http.StatusCreated,
			expectedResponse: map[string]string{"message": "User created successfully"},
		},
		{
			name: "User already exists",
			payload: types.SignupUserPayload{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "password",
			},
			mockStore: &mockUserStore{
				GetUserByEmailFunc: func(email string) (*types.User, error) {
					return &types.User{}, nil
				},
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: map[string]string{"error": "user with this email already exists"},
		},
		{
			name: "Invalid payload",
			payload: types.SignupUserPayload{
				FirstName: "",
				LastName:  "",
				Email:     "invalid-email",
				Password:  "short",
			},
			mockStore:      &mockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]string{
				"error": fmt.Sprintf("%s: %v", utils.ErrInvalidPayload, validator.ValidationErrors{}),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tc.payload)
			req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(payloadBytes))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler := &Handler{store: tc.mockStore}
			handler.signUp(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			var responseBody map[string]string
			err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
			if err != nil {
				t.Fatalf("could not unmarshal response body: %v", err)
			}

			if len(responseBody) != len(tc.expectedResponse) {
				t.Errorf("handler returned unexpected body: got %v want %v", responseBody, tc.expectedResponse)
			}

			for key, value := range tc.expectedResponse {
				if responseBody[key] != value {
					t.Errorf("handler returned unexpected body: got %v want %v", responseBody, tc.expectedResponse)
				}
			}
		})
	}
}

func mockComparePasswords(storedPassword, providedPassword string) error {
	if storedPassword != providedPassword {
		return fmt.Errorf("invalid password")
	}
	return nil
}

func mockGenerateToken(secret []byte, userID int, expiration time.Duration) (string, error) {
	return "mocked-token", nil
}


func TestLogin(t *testing.T) {
	originalValidate := utils.Validate
	utils.Validate = &mockValidator{}
	defer func() { utils.Validate = originalValidate }()

	mockCfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "mock-secret",
			Expiration: 3600, // Example expiration time in seconds
		},
	}

	tests := []struct {
		name             string
		payload          types.LoginUserPayload
		mockStore        *mockUserStore
		expectedStatus   int
		expectedResponse map[string]string
	}{
		{
			name: "Valid login",
			payload: types.LoginUserPayload{
				Email:    "john.doe@example.com",
				Password: "password",
			},
			mockStore: func() *mockUserStore {
				hashedPassword, _ := auth.HashPassword("password")
				return &mockUserStore{
					GetUserByEmailFunc: func(email string) (*types.User, error) {
						return &types.User{ID: 1, Password: hashedPassword}, nil
					},
				}
			}(),
			expectedStatus:   http.StatusOK,
			expectedResponse: map[string]string{"token": "mocked-token"},
		},
		{
			name: "User not found",
			payload: types.LoginUserPayload{
				Email:    "john.doe@example.com",
				Password: "password",
			},
			mockStore: &mockUserStore{
				GetUserByEmailFunc: func(email string) (*types.User, error) {
					return nil, sql.ErrNoRows
				},
			},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: map[string]string{"error": utils.ErrUserNotFound},
		},
		{
			name: "Invalid payload",
			payload: types.LoginUserPayload{
				Email:    "invalid-email",
				Password: "password",
			},
			mockStore:      &mockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]string{
				"error": fmt.Sprintf("%s: %v", utils.ErrInvalidPayload, validator.ValidationErrors{}),
			},
		},
		{
			name: "Invalid password",
			payload: types.LoginUserPayload{
				Email:    "john.doe@example.com",
				Password: "wrong-password",
			},
			mockStore: func() *mockUserStore {
				hashedPassword, _ := auth.HashPassword("password")
				return &mockUserStore{
					GetUserByEmailFunc: func(email string) (*types.User, error) {
						return &types.User{ID: 1, Password: hashedPassword}, nil
					},
				}
			}(),
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: map[string]string{"error": utils.ErrUnauthorized},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tc.payload)
			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payloadBytes))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler := &Handler{
				store:            tc.mockStore,
				cfg:              mockCfg,
				comparePasswords: auth.ComparePasswords,
				generateToken:    mockGenerateToken, // Correctly inject mockGenerateToken
			}
			handler.login(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			var responseBody map[string]string
			err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
			if err != nil {
				t.Fatalf("could not unmarshal response body: %v", err)
			}

			if len(responseBody) != len(tc.expectedResponse) {
				t.Errorf("handler returned unexpected body: got %v want %v", responseBody, tc.expectedResponse)
			}

			for key, value := range tc.expectedResponse {
				if responseBody[key] != value {
					t.Errorf("handler returned unexpected body: got %v want %v", responseBody, tc.expectedResponse)
				}
			}
		})
	}
}