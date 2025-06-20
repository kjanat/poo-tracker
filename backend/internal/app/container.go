package app

import (
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/database"
)

// Container holds all application dependencies
type Container struct {
	Config   *Config
	Database database.Database

	// Services will be added in later phases
	// UserService     service.UserService
	// MealService     service.MealService
	// etc...

	// Repositories will be added in later phases
	// UserRepository     repository.UserRepository
	// MealRepository     repository.MealRepository
	// etc...
}

// NewContainer creates a new dependency injection container
func NewContainer() (*Container, error) {
	// Load configuration
	config := LoadConfig()

	// Setup database
	dbConfig := database.GetConfigFromEnv()
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := db.Migrate(); err != nil {
		return nil, err
	}

	container := &Container{
		Config:   config,
		Database: db,
	}

	// Initialize repositories and services will be added in later phases

	return container, nil
}

// Cleanup closes all resources
func (c *Container) Cleanup() error {
	if c.Database != nil {
		return c.Database.Close()
	}
	return nil
}
