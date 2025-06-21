package app

import (
	"fmt"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/analytics"
	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/database"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/repository/memory"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service"
	"go.uber.org/zap"
)

// Container holds all application dependencies
type Container struct {
	Config   *Config
	Database database.Database

	// Repositories
	UserRepository          user.Repository
	BowelMovementRepository bowelmovement.Repository
	MealRepository          meal.Repository
	MedicationRepository    medication.Repository
	SymptomRepository       symptom.Repository

	// Services
	UserService          user.Service
	BowelMovementService bowelmovement.Service
	MealService          meal.Service
	MedicationService    medication.Service
	SymptomService       symptom.Service
	AnalyticsService     analytics.Service
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

	// Validate database connection
	if sqlDB, err := db.GetDB().DB(); err != nil {
		return nil, fmt.Errorf("failed to get underlying database connection: %w", err)
	} else if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database connection validation failed: %w", err)
	}

	// Run migrations
	if err := db.Migrate(); err != nil {
		return nil, err
	}

	container := &Container{
		Config:   config,
		Database: db,
	}

	// Initialize repositories
	container.UserRepository = memory.NewUserRepository()
	container.BowelMovementRepository = memory.NewBowelMovementRepository()
	container.MealRepository = memory.NewMealRepository()
	container.MedicationRepository = memory.NewMedicationRepository()
	container.SymptomRepository = memory.NewSymptomRepository()

	// Initialize services
	container.UserService = service.NewUserService(container.UserRepository)
	container.BowelMovementService = service.NewBowelMovementService(container.BowelMovementRepository)
	container.MealService = service.NewMealService(container.MealRepository)
	container.MedicationService = service.NewMedicationService(container.MedicationRepository)
	container.SymptomService = service.NewSymptomService(container.SymptomRepository)
	container.AnalyticsService = service.NewAnalyticsService(
		container.BowelMovementRepository,
		container.MealRepository,
		container.SymptomRepository,
		container.MedicationService,
		zap.NewNop(),
		&service.AnalyticsServiceConfig{
			DefaultMedicationLimit: 1000,
			DefaultDataWindow:      30 * 24 * time.Hour,
		},
	)

	return container, nil
}

// Cleanup closes all resources
func (c *Container) Cleanup() error {
	if c.Database != nil {
		return c.Database.Close()
	}
	return nil
}
