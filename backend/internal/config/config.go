package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Database
	DatabaseURL  string
	DatabaseHost string
	DatabasePort string
	DatabaseUser string
	DatabasePass string
	DatabaseName string

	// Server
	ServerPort string
	ServerHost string

	// Environment
	Environment string

	// Redis
	RedisURL string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		DatabaseHost: getEnv("DATABASE_HOST", "localhost"),
		DatabasePort: getEnv("DATABASE_PORT", "5432"),
		DatabaseUser: getEnv("DATABASE_USER", "postgres"),
		DatabasePass: getEnv("DATABASE_PASSWORD", ""),
		DatabaseName: getEnv("DATABASE_NAME", "spectra"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		ServerHost:   getEnv("SERVER_HOST", "0.0.0.0"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379"),
	}

	// Validate required fields
	if cfg.DatabaseURL == "" && cfg.DatabasePass == "" {
		return nil, fmt.Errorf("DATABASE_URL or DATABASE_PASSWORD is required")
	}

	return cfg, nil
}

// GetDatabaseURL returns the full database connection string
func (c *Config) GetDatabaseURL() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DatabaseUser, c.DatabasePass, c.DatabaseHost, c.DatabasePort, c.DatabaseName)
}

// getEnv reads an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
