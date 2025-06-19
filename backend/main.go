package main

import (
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/service"
	"github.com/kjanat/poo-tracker/backend/server"
)

func main() {
	app := server.New(repository.NewMemory(), service.AvgBristol{})
	app.Engine.Run() // defaults to :8080
}
