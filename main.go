package main

import (
	"dimiplan-backend/config"
	"dimiplan-backend/routes"
	"log"
)

func main() {
	cfg := config.Load()
	app := routes.Setup(cfg)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
