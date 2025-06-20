package database

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowel_movement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
)

type SQLiteDatabase struct {
	db *gorm.DB
}

func NewSQLiteDatabase(config Config) (*SQLiteDatabase, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(filepath.Dir(config.DSN), 0755); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	})
	if err != nil {
		return nil, err
	}

	return &SQLiteDatabase{db: db}, nil
}

func (s *SQLiteDatabase) GetDB() *gorm.DB {
	return s.db
}

func (s *SQLiteDatabase) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *SQLiteDatabase) Migrate() error {
	// Auto-migrate domain models
	return s.db.AutoMigrate(
		&user.User{},
		&bowel_movement.BowelMovement{},
	)
}
