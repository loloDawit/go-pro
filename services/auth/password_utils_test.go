package auth

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		hashFunc    HashPasswordFunc
		expectError bool
	}{
		{
			name:        "Valid password",
			password:    "password123",
			hashFunc:    bcrypt.GenerateFromPassword,
			expectError: false,
		},
		{
			name:        "Empty password",
			password:    "",
			hashFunc:    bcrypt.GenerateFromPassword,
			expectError: false,
		},
		{
			name:        "Long password",
			password:    "aVeryLongPasswordWithMultipleCharactersAndSymbols123!@#",
			hashFunc:    bcrypt.GenerateFromPassword,
			expectError: false,
		},
		{
			name:     "Simulated error",
			password: "password123",
			hashFunc: func(password []byte, cost int) ([]byte, error) {
				return nil, errors.New("simulated error")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := hashPasswordWithFunc(tt.password, tt.hashFunc)
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, hashedPassword)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)
				err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(tt.password))
				assert.NoError(t, err)
			}
		})
	}
}

func TestComparePasswords(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		hashedPassword string
		expectError    bool
	}{
		{
			name:           "Matching passwords",
			password:       "password123",
			hashedPassword: func() string { h, _ := HashPassword("password123"); return h }(),
			expectError:    false,
		},
		{
			name:           "Non-matching passwords",
			password:       "password123",
			hashedPassword: func() string { h, _ := HashPassword("password321"); return h }(),
			expectError:    true,
		},
		{
			name:           "Empty password",
			password:       "",
			hashedPassword: func() string { h, _ := HashPassword(""); return h }(),
			expectError:    false,
		},
		{
			name:           "Invalid hashed password",
			password:       "password123",
			hashedPassword: "invalidhashedpassword",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ComparePasswords(tt.hashedPassword, tt.password)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
