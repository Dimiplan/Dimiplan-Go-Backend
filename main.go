package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"dimiplan-backend/config"
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/migrate"
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

	if !fiber.IsChild() {
		if err := client.Schema.Create(context.Background(), migrate.WithDropIndex(true), migrate.WithDropColumn(true)); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
	}

	app, redis := server.Setup(cfg)
	routes.Setup(app, cfg, client)

	log.Infof("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":"+cfg.Port, fiber.ListenConfig{
		EnablePrefork: true,
		CertFile:      "./keys/cert.pem",
		CertKeyFile:   "./keys/key.pem",
	}))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Info("Shutting down server...")
	client.Close()
	redis.Close()
	app.Shutdown()
}
