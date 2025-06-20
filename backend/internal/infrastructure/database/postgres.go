package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowel_movement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
)

type PostgreSQLDatabase struct {
	db *gorm.DB
}

func NewPostgreSQLDatabase(config Config) (*PostgreSQLDatabase, error) {
	db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	})
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return &PostgreSQLDatabase{db: db}, nil
}

func (p *PostgreSQLDatabase) GetDB() *gorm.DB {
	return p.db
}

func (p *PostgreSQLDatabase) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (p *PostgreSQLDatabase) Migrate() error {
	// Auto-migrate domain models
	return p.db.AutoMigrate(
		&user.User{},
		&bowel_movement.BowelMovement{},
	)
}
