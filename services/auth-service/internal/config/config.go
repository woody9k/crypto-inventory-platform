package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the auth service
type Config struct {
	Port        string
	Environment string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	JWTExpiry   time.Duration
	LogLevel    string
	CORSOrigins []string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENV", "development"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://crypto_user:crypto_pass_dev@localhost:5432/crypto_inventory?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://:redis_pass_dev@localhost:6379/0"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-key-change-in-production"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}

	// Parse JWT expiry
	jwtExpiryStr := getEnv("JWT_EXPIRY", "24h")
	jwtExpiry, err := time.ParseDuration(jwtExpiryStr)
	if err != nil {
		return nil, err
	}
	config.JWTExpiry = jwtExpiry

	// Parse CORS origins
	corsOrigins := getEnv("CORS_ORIGINS", "http://localhost:3000")
	config.CORSOrigins = []string{corsOrigins}

	return config, nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with fallback
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

// getEnvAsBool gets an environment variable as boolean with fallback
func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return fallback
}
