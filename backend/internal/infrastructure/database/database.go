package database

import (
	"fmt"
	"os"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database interface for strategy pattern
type Database interface {
	GetDB() *gorm.DB
	Close() error
	Migrate() error
}

// DatabaseType represents the type of database to use
type DatabaseType string

const (
	SQLite     DatabaseType = "sqlite"
	PostgreSQL DatabaseType = "postgres"
)

// Config holds database configuration
type Config struct {
	Type     DatabaseType
	DSN      string
	LogLevel logger.LogLevel
}

// NewDatabase creates a new database instance based on configuration
func NewDatabase(config Config) (Database, error) {
	switch config.Type {
	case SQLite:
		return NewSQLiteDatabase(config)
	case PostgreSQL:
		return NewPostgreSQLDatabase(config)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}
}

// GetConfigFromEnv creates database config from environment variables
func GetConfigFromEnv() Config {
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite" // default to SQLite for development
	}

	var dsn string
	var logLevel logger.LogLevel

	switch DatabaseType(dbType) {
	case SQLite:
		dsn = os.Getenv("SQLITE_DSN")
		if dsn == "" {
			dsn = "./data/poo-tracker.db"
		}
	case PostgreSQL:
		dsn = os.Getenv("DATABASE_URL")
		if dsn == "" {
			// Fallback to individual env vars
			host := getEnvOrDefault("DB_HOST", "localhost")
			port := getEnvOrDefault("DB_PORT", "5432")
			user := getEnvOrDefault("DB_USER", "postgres")
			password := getEnvOrDefault("DB_PASSWORD", "postgres")
			dbname := getEnvOrDefault("DB_NAME", "poo_tracker")
			sslmode := getEnvOrDefault("DB_SSLMODE", "disable")

			dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
				host, port, user, password, dbname, sslmode)
		}
	}

	// Set log level based on environment
	if os.Getenv("DEBUG") == "true" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Warn
	}

	return Config{
		Type:     DatabaseType(dbType),
		DSN:      dsn,
		LogLevel: logLevel,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
