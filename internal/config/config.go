package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DatabaseURL string
}

func LoadConfig() *Config {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	return &Config {
		Port: port,
		DatabaseURL: databaseURL,
	}
}