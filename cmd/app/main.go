package main

import (
	"log"

	"github.com/hong195/web-server/config"
	"github.com/hong195/web-server/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
