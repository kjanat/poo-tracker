package main

import (
	"log"
	"os"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/service"
	"github.com/kjanat/poo-tracker/backend/server"
)

func main() {
	// Get JWT secret from environment variable with fallback for development
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev_secret_change_in_production"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable for production.")
	}

	userRepo := repository.NewMemoryUserRepository()
	auth := &service.JWTAuthService{
		UserRepo: userRepo,
		Secret:   jwtSecret,
		Expiry:   24 * time.Hour,
	}

	bowelRepo := repository.NewMemoryBowelRepo()
	mealRepo := repository.NewMemoryMealRepo()
	app := server.New(bowelRepo, mealRepo, service.AvgBristol{}, auth)

	log.Println("Starting Poo Tracker server...")
	if err := app.Engine.Run(); err != nil {
		log.Printf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
