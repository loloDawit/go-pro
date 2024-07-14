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

	// Test cases
	tests := []struct {
		name             string
		directory        string
		environment      string
		mockReadFileFunc func(string) ([]byte, error)
		expectedPanicMsg string
		expectedConfig   *Config
	}{
		{
			name:        "Successful Load",
			directory:   "./config",
			environment: "test",
			mockReadFileFunc: func(filename string) ([]byte, error) {
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
				return []byte(configYAML), nil
			},
			expectedPanicMsg: "",
			expectedConfig: &Config{
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
			},
		},
		{
			name:             "Default Directory and Environment",
			directory:        "",
			environment:      "",
			mockReadFileFunc: func(filename string) ([]byte, error) { return nil, fmt.Errorf("file not found: %s", filename) },
			expectedPanicMsg: "could not read ./config/development.yml config file: file not found: ./config/development.yml",
			expectedConfig:   nil,
		},
		{
			name:        "Error Reading Main Config File",
			directory:   "./config",
			environment: "test",
			mockReadFileFunc: func(filename string) ([]byte, error) {
				return nil, fmt.Errorf("file not found: %s", filename)
			},
			expectedPanicMsg: "could not read ./config/test.yml config file: file not found: ./config/test.yml",
			expectedConfig:   nil,
		},
		{
			name:        "Error Parsing Main Config File",
			directory:   "./config",
			environment: "test",
			mockReadFileFunc: func(filename string) ([]byte, error) {
				invalidYAML := "invalid yaml content"
				return []byte(invalidYAML), nil
			},
			expectedPanicMsg: "could not parse ./config/test.yml config file: yaml: unmarshal errors:",
			expectedConfig:   nil,
		},
		{
			name:        "Error Reading Additional Config File",
			directory:   "./config",
			environment: "test",
			mockReadFileFunc: func(filename string) ([]byte, error) {
				if filename == "./config/checkout-api-config-test.yml" {
					return nil, fmt.Errorf("file not found: %s", filename)
				}
				validYAML := `
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
				return []byte(validYAML), nil
			},
			expectedPanicMsg: "",
			expectedConfig: &Config{
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
			},
		},
		{
			name:        "Error Parsing Additional Config File",
			directory:   "./config",
			environment: "test",
			mockReadFileFunc: func(filename string) ([]byte, error) {
				if filename == "./config/checkout-api-config-test.yml" {
					invalidYAML := "invalid yaml content"
					return []byte(invalidYAML), nil
				}
				validYAML := `
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
				return []byte(validYAML), nil
			},
			expectedPanicMsg: "failed to unmarshal configuration file (checkout-api-config-test.yml): yaml: unmarshal errors:",
			expectedConfig:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock readFileFunc function
			mockReadFile(tt.mockReadFileFunc)
			defer restoreReadFile()

			if tt.expectedPanicMsg != "" {
				// Test for panic
				defer func() {
					if r := recover(); r != nil {
						assert.Contains(t, r.(error).Error(), tt.expectedPanicMsg)
					}
				}()
			}

			// Test case
			ctx := context.Background()
			cfg := LoadConfig(ctx, tt.directory, tt.environment, "deployment")

			// Check results
			if tt.expectedConfig != nil {
				assert.Equal(t, tt.expectedConfig, cfg)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	mockLookupEnv(func(key string) (string, bool) {
		envVars := map[string]string{
			"DB_USER":     "test_user",
			"DB_PASSWORD": "test_password",
		}
		val, ok := envVars[key]
		return val, ok
	})
	defer restoreLookupEnv()

	tests := []struct {
		key      string
		expected string
	}{
		{"DB_USER", "test_user"},
		{"DB_PASSWORD", "test_password"},
		{"DB_HOSTNAME", ""},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			assert.Equal(t, tt.expected, getEnv(tt.key))
		})
	}
}
