package main

import (
	"log"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/database"
	"url-shortener/internal/handlers"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load Config
	cfg := config.LoadConfig()

	// Initializing Database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis Client
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Username: cfg.RedisUsername,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer rdb.Close()

	// Set up routes
	http.HandleFunc("/shorten", handlers.ShortenURLHandler(db, rdb))
	http.HandleFunc("/short/", handlers.ExpandURLHandler(db, rdb))

	// Start the server
	log.Printf("Server running on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":" + cfg.Port, nil))
}