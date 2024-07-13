package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type JWTConfig struct {
	Expiration int64  `yaml:"expiration"`
	Secret     string `yaml:"secret"`
}

type Config struct {
	Environment string    `yaml:"environment"`
	DBuser      string    `yaml:"db_user"`
	DBpassword  string    `yaml:"db_password"`
	DBaddr      string    `yaml:"db_addr"`
	DBname      string    `yaml:"db_name"`
	JWT         JWTConfig `yaml:"jwt"`
	Address     string    `yaml:"address"`
}

// DefaultConfig creates a default config
func DefaultConfig(environment string) *Config {
	return &Config{
		Environment: environment,
		Address:     ":8080",
		JWT:         DefaultJWTConfig(),
	}
}

func DefaultJWTConfig() JWTConfig {
	return JWTConfig{
		Expiration: 3600, // Default to 1 hour
		Secret:     "default_secret",
	}
}


const configFormat = "checkout-api-config-%s.yml"

// LoadConfig creates a new Config instance and populates it with the environment file found in the configuration directory.
func LoadConfig(ctx context.Context, directory string, environment string, deployment string) *Config {
	// Default value if directory is not provided
	if len(directory) < 1 {
		directory = "./config"
	}

	// Default value if environment is not provided
	if len(environment) < 1 {
		environment = "development"
	}

	// Start with the "default" config
	cfg := DefaultConfig(environment)
	log.Printf("Loading config for environment: %s", environment)

	// Load common environment variables
	cfg.DBuser = getEnv("DB_USER")
	cfg.DBpassword = getEnv("DB_PASSWORD")
	cfg.DBaddr = getEnv("DB_HOSTNAME")
	cfg.DBname = getEnv("DB_NAME")
	cfg.JWT.Secret = getEnv("JWT_SECRET")

	// Load YAML configuration based on the environment
	fileName := fmt.Sprintf("%s/%s.yml", directory, environment)
	yamlFile, err := os.ReadFile(fileName)
	if err != nil {
		panic(fmt.Errorf("could not read %s config file: %w", fileName, err))
	}

	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		panic(fmt.Errorf("could not parse %s config file: %w", fileName, err))
	}

	// Load additional environment-specific config if it exists
	secretConfig := fmt.Sprintf(configFormat, environment)
	if data, err := os.ReadFile(directory + "/" + secretConfig); err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			panic(fmt.Errorf("failed to unmarshal configuration file (%s): %v", secretConfig, err))
		}
	}

	return cfg
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return ""
}
