package config

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mocking the file reading function
func mockReadFile(fn func(string) ([]byte, error)) {
	readFileFunc = fn
}

func restoreReadFile() {
	readFileFunc = os.ReadFile
}

// Mock environment variable lookup
func mockLookupEnv(fn func(string) (string, bool)) {
	lookupEnvFunc = fn
}

func restoreLookupEnv() {
	lookupEnvFunc = os.LookupEnv
}

func TestLoadConfig(t *testing.T) {
	// Mock environment variables
	mockLookupEnv(func(key string) (string, bool) {
		envVars := map[string]string{
			"DB_USER":     "test_user",
			"DB_PASSWORD": "test_password",
			"DB_HOSTNAME": "test_hostname",
			"DB_NAME":     "test_db",
			"JWT_SECRET":  "test_secret",
		}
		val, ok := envVars[key]
		return val, ok
	})
	defer restoreLookupEnv()

	// Mock YAML file content
	configYAML := `
environment: test
db_user: test_user
db_password: test_password
db_addr: test_hostname
db_name: test_db
jwt:
  expiration: 7200
  secret: test_secret
address: :8080
`

	// Mock readFileFunc function
	mockReadFile(func(filename string) ([]byte, error) {
		if filename == "./config/test.yml" {
			return []byte(configYAML), nil
		}
		return nil, fmt.Errorf("file not found: %s", filename)
	})
	defer restoreReadFile()

	// Test case
	ctx := context.Background()
	cfg := LoadConfig(ctx, "./config", "test", "deployment")

	expectedConfig := &Config{
		Environment: "test",
		DBuser:      "test_user",
		DBpassword:  "test_password",
		DBaddr:      "test_hostname",
		DBname:      "test_db",
		JWT: JWTConfig{
			Expiration: 7200,
			Secret:     "test_secret",
		},
		Address: ":8080",
	}

	assert.Equal(t, expectedConfig, cfg)
}
