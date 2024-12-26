package main

import (
	"log"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/database"
	"url-shortener/internal/handlers"
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

	// Set up routes
	http.HandleFunc("/shorten", handlers.ShortenURLHandler(db))
	http.HandleFunc("/short/", handlers.ExpandURLHandler(db))

	// Start the server
	log.Printf("Server running on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":" + cfg.Port, nil))
}