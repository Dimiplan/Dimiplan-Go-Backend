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
	client, err := ent.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=dimiplan password=postgres")
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	app := routes.Setup(cfg)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
