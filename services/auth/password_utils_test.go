package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Failed to hash password: %v", err)
	}
	if hashedPassword == "" {
		t.Errorf("Hashed password should not be empty")
	}
}

func TestComparePasswords(t *testing.T) {
	password := "mysecretpassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Failed to hash password: %v", err)
	}

	err = ComparePasswords(hashedPassword, password)
	if err != nil {
		t.Errorf("Passwords should match: %v", err)
	}

	wrongPassword := "wrongpassword"
	err = ComparePasswords(hashedPassword, wrongPassword)
	if err == nil {
		t.Errorf("Expected error when passwords do not match, got nil")
	}
}
