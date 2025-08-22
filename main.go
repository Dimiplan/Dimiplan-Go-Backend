package main

import (
	"context"

	"dimiplan-backend/config"
	"dimiplan-backend/ent"
	"dimiplan-backend/routes"
	"dimiplan-backend/server"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	client, err := ent.Open("postgres", cfg.DatabaseString)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	app := server.Setup(cfg)
	routes.Setup(app, cfg, client)

	log.Infof("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":"+cfg.Port, fiber.ListenConfig{
		// EnablePrefork: true,
		CertFile:    "./keys/cert.pem",
		CertKeyFile: "./keys/key.pem",
	}))
}
