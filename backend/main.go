package main

import (
	"github.com/kjanat/poo-tracker/backend/server"
)

func main() {
	r := server.New()
	r.Run() // defaults to :8080
}
