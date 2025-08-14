package main

import (
	"context"
	"log"

	"dimiplan-backend/config"
	"dimiplan-backend/ent"
	"dimiplan-backend/routes"

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

	app := routes.Setup(cfg, client)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
