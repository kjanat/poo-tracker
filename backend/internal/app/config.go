package app

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	// Server config
	Port string
	Host string

	// Security config
	JWTSecret      string
	BCryptCost     int
	SessionTimeout int // minutes

	// Database config will be handled by database package

	// File upload config
	MaxUploadSize int64 // bytes
	UploadPath    string

	// Environment
	Environment string
	Debug       bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		Port:           getEnvOrDefault("PORT", "8080"),
		Host:           getEnvOrDefault("HOST", "localhost"),
		JWTSecret:      getEnvOrDefault("JWT_SECRET", "your-secret-key-change-this-in-production"),
		BCryptCost:     getEnvAsIntOrDefault("BCRYPT_COST", 12),
		SessionTimeout: getEnvAsIntOrDefault("SESSION_TIMEOUT", 1440),           // 24 hours
		MaxUploadSize:  getEnvAsInt64OrDefault("MAX_UPLOAD_SIZE", 10*1024*1024), // 10MB
		UploadPath:     getEnvOrDefault("UPLOAD_PATH", "./uploads"),
		Environment:    getEnvOrDefault("ENVIRONMENT", "development"),
		Debug:          getEnvOrDefault("DEBUG", "false") == "true",
	}

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64OrDefault(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
