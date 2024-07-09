package auth

import "golang.org/x/crypto/bcrypt"

// HashPasswordFunc defines the type for the hashing function.
type HashPasswordFunc func(password []byte, cost int) ([]byte, error)

// HashPassword hashes the given password using bcrypt.
func HashPassword(password string) (string, error) {
	return hashPasswordWithFunc(password, bcrypt.GenerateFromPassword)
}

// hashPasswordWithFunc is a helper function that allows injecting a custom hashing function.
func hashPasswordWithFunc(password string, hashFunc HashPasswordFunc) (string, error) {
	hashedPassword, err := hashFunc([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePasswords compares the given hashed password with the given password.
func ComparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
