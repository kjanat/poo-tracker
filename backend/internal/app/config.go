package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration
type Config struct {
	// Server config
	Port            string        `json:"port"`
	Host            string        `json:"host"`
	ReadTimeout     time.Duration `json:"readTimeout"`
	WriteTimeout    time.Duration `json:"writeTimeout"`
	ShutdownTimeout time.Duration `json:"shutdownTimeout"`

	// Security config
	JWTSecret      string        `json:"-"` // Never serialize secrets
	JWTExpiry      time.Duration `json:"jwtExpiry"`
	BCryptCost     int           `json:"bcryptCost"`
	SessionTimeout int           `json:"sessionTimeout"` // minutes

	// Database config
	DatabaseURL      string        `json:"-"` // Never serialize connection strings
	DatabaseMaxConns int           `json:"databaseMaxConns"`
	DatabaseTimeout  time.Duration `json:"databaseTimeout"`

	// Security & Rate Limiting
	RateLimitRequests int           `json:"rateLimitRequests"`
	RateLimitWindow   time.Duration `json:"rateLimitWindow"`
	CORSOrigins       []string      `json:"corsOrigins"`

	// File upload config
	MaxUploadSize int64  `json:"maxUploadSize"` // bytes
	UploadPath    string `json:"uploadPath"`

	// Logging config
	LogLevel  string `json:"logLevel"`
	LogFormat string `json:"logFormat"` // json, text

	// Feature flags
	EnableMetrics     bool `json:"enableMetrics"`
	EnableProfiling   bool `json:"enableProfiling"`
	EnableSwagger     bool `json:"enableSwagger"`
	EnableHealthCheck bool `json:"enableHealthCheck"`

	// Environment
	Environment string `json:"environment"`
	Debug       bool   `json:"debug"`
}

// LoadConfig loads configuration from environment variables with validation
func LoadConfig() *Config {
	config := &Config{
		// Server config
		Port:            getEnvOrDefault("PORT", "8080"),
		Host:            getEnvOrDefault("HOST", "localhost"),
		ReadTimeout:     getEnvAsDurationOrDefault("READ_TIMEOUT", 15*time.Second),
		WriteTimeout:    getEnvAsDurationOrDefault("WRITE_TIMEOUT", 15*time.Second),
		ShutdownTimeout: getEnvAsDurationOrDefault("SHUTDOWN_TIMEOUT", 10*time.Second),

		// Security config
		JWTSecret:      getEnvOrDefault("JWT_SECRET", "your-secret-key-change-this-in-production"),
		JWTExpiry:      getEnvAsDurationOrDefault("JWT_EXPIRY", 24*time.Hour),
		BCryptCost:     getEnvAsIntOrDefault("BCRYPT_COST", 12),
		SessionTimeout: getEnvAsIntOrDefault("SESSION_TIMEOUT", 1440), // 24 hours

		// Database config
		DatabaseURL:      getEnvOrDefault("DATABASE_URL", ""),
		DatabaseMaxConns: getEnvAsIntOrDefault("DATABASE_MAX_CONNS", 25),
		DatabaseTimeout:  getEnvAsDurationOrDefault("DATABASE_TIMEOUT", 30*time.Second),

		// Security & Rate Limiting
		RateLimitRequests: getEnvAsIntOrDefault("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvAsDurationOrDefault("RATE_LIMIT_WINDOW", time.Minute),
		CORSOrigins:       getEnvAsSliceOrDefault("CORS_ORIGINS", []string{"http://localhost:3000"}),

		// File upload config
		MaxUploadSize: getEnvAsInt64OrDefault("MAX_UPLOAD_SIZE", 10*1024*1024), // 10MB
		UploadPath:    getEnvOrDefault("UPLOAD_PATH", "./uploads"),

		// Logging config
		LogLevel:  getEnvOrDefault("LOG_LEVEL", "info"),
		LogFormat: getEnvOrDefault("LOG_FORMAT", "json"),

		// Feature flags
		EnableMetrics:     getEnvAsBoolOrDefault("ENABLE_METRICS", true),
		EnableProfiling:   getEnvAsBoolOrDefault("ENABLE_PROFILING", false),
		EnableSwagger:     getEnvAsBoolOrDefault("ENABLE_SWAGGER", true),
		EnableHealthCheck: getEnvAsBoolOrDefault("ENABLE_HEALTH_CHECK", true),

		// Environment
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
		Debug:       getEnvAsBoolOrDefault("DEBUG", false),
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		panic(fmt.Sprintf("Configuration validation failed: %v", err))
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

func getEnvAsDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsSliceOrDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma, trim spaces
		result := make([]string, 0)
		for _, item := range strings.Split(value, ",") {
			if trimmed := strings.TrimSpace(item); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

// Validate validates the configuration values
func (c *Config) Validate() error {
	// Validate JWT secret in production
	if c.Environment == "production" && c.JWTSecret == "your-secret-key-change-this-in-production" {
		return fmt.Errorf("JWT_SECRET must be set in production environment")
	}

	// Validate JWT secret minimum length
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}

	// Validate bcrypt cost
	if c.BCryptCost < 4 || c.BCryptCost > 31 {
		return fmt.Errorf("BCRYPT_COST must be between 4 and 31")
	}

	// Validate timeouts
	if c.ReadTimeout <= 0 {
		return fmt.Errorf("READ_TIMEOUT must be positive")
	}
	if c.WriteTimeout <= 0 {
		return fmt.Errorf("WRITE_TIMEOUT must be positive")
	}

	// Validate log level
	validLogLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLogLevels, c.LogLevel) {
		return fmt.Errorf("LOG_LEVEL must be one of: %v", validLogLevels)
	}

	// Validate log format
	validLogFormats := []string{"json", "text"}
	if !contains(validLogFormats, c.LogFormat) {
		return fmt.Errorf("LOG_FORMAT must be one of: %v", validLogFormats)
	}

	// Validate environment
	validEnvironments := []string{"development", "staging", "production"}
	if !contains(validEnvironments, c.Environment) {
		return fmt.Errorf("ENVIRONMENT must be one of: %v", validEnvironments)
	}

	return nil
}

// IsProduction returns true if running in production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment returns true if running in development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// GetServerAddress returns the server address
func (c *Config) GetServerAddress() string {
	return c.Host + ":" + c.Port
}

func contains(validValues []string, value string) bool {
	for _, v := range validValues {
		if v == value {
			return true
		}
	}
	return false
}
