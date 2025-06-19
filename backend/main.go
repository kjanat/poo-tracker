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
	userRepo := repository.NewMemoryUserRepository()
	auth := &service.JWTAuthService{
		UserRepo: userRepo,
		Secret:   "dev_secret", // TODO: use env var
		Expiry:   24 * time.Hour,
	}
	server.AuthService = auth

	mem := repository.NewMemory()
	app := server.New(mem, mem, service.AvgBristol{})

	log.Println("Starting Poo Tracker server...")
	if err := app.Engine.Run(); err != nil {
		log.Printf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
