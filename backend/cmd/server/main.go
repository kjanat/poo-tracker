package main

import (
	"log"

	"github.com/kjanat/poo-tracker/backend/internal/app"
)

func main() {
	// Create and run the application
	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
