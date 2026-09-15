package config

import (
	"io"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/RonIT-401/order-service/internal/app/config/section"
)

type Config struct {
	Repository section.Repository
	Processor  section.Processor
	Monitor    section.Monitor
}

var Root Config

func Load(args LoadArgs) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	if err := envconfig.Process("APP", &Root); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}

type LoadArgs struct {
	Output          io.Writer `json:"-"`
	EnableSimpleLog bool
}
