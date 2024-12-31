package main

import (
	"log"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/database"
	"url-shortener/internal/handlers"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
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
    mux := http.NewServeMux()
    mux.HandleFunc("/shorten", handlers.ShortenURLHandler(db, rdb, cfg.KoyebURL))
    mux.HandleFunc("/short/", handlers.ExpandURLHandler(db, rdb))

    // Enable CORS using the library
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"}, // Allow all origins
        AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type"},
        AllowCredentials: true,
    })

	// Start the server with CORS middleware
    handler := c.Handler(mux)
    log.Printf("Server running on port %s", cfg.Port)
    log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
}