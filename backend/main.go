package main

import (
	"log"
	"os"

	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/service"
	"github.com/kjanat/poo-tracker/backend/server"
)

func main() {
	mem := repository.NewMemory()
	app := server.New(mem, mem, service.AvgBristol{})

	log.Println("Starting Poo Tracker server...")
	if err := app.Engine.Run(); err != nil {
		log.Printf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
