package main

import (
	"context"
	"log"

	"github.com/RonIT-401/order-service/internal/app/config"
	rhealth "github.com/RonIT-401/order-service/internal/app/handler/http/health"
	rprocessor "github.com/RonIT-401/order-service/internal/app/processor/http"
	rcpostgres "github.com/RonIT-401/order-service/internal/app/repository/conn/postgres"
)

func main() {
	config.Load(config.LoadArgs{})

	cfg := config.Root

	db, err := rcpostgres.NewClient(context.Background(), cfg.Repository.Postgres)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	log.Printf("Connected to database: %s", cfg.Repository.Postgres.Name)
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	hHealth := rhealth.NewHandler()

	httpProc := rprocessor.NewHTTP(hHealth, cfg.Processor.WebServer)

	if err := httpProc.Serve(); err != nil {
		//nolint:gocritic
		log.Fatalf("HTTP server stopped: %v", err)
	}
}
