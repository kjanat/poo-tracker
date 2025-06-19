package main

import (
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/service"
	"github.com/kjanat/poo-tracker/backend/server"
)

func main() {
	mem := repository.NewMemory()
	app := server.New(mem, mem, service.AvgBristol{})
	if err := app.Engine.Run(); err != nil {
		panic(err)
	}
}
