package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes the given password using bcrypt.
func HashPassword(password string) (string, error) {
	// Generate a hashed password with a default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePasswords compares the given hashed password with the given password.
func ComparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
