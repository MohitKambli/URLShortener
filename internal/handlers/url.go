package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"url-shortener/internal/repositories"
	"url-shortener/internal/services"
	"github.com/redis/go-redis/v9"
	"fmt"
)

func ShortenURLHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			OriginalURL string `json:"original_url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}
		if req.OriginalURL == "" {
			http.Error(w, "original_url is required", http.StatusBadRequest)
			return
		}
		repo := &repositories.URLRepository{DB: db}
		service := &services.URLService{Repo: repo}
		shortenedURL, err := service.ShortenURL(req.OriginalURL)
		if err != nil {
			http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
			return
		}

		// Cache in Redis
		ctx := context.Background()
		rdb.Set(ctx, shortenedURL.ShortURL, shortenedURL.OriginalURL, 2*time.Hour)

		fullShortURL := "http://localhost:4000/short/" + shortenedURL.ShortURL
		response := struct {
			ShortURL string `json:"short_url"`
			OriginalURL string `json:original_url`
		}{
			ShortURL:    fullShortURL,
			OriginalURL: shortenedURL.OriginalURL,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func ExpandURLHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// shortURL := r.URL.Query().Get("short_url")
		shortURL := r.URL.Path[len("/short/"):]
		if shortURL == "" {
			http.Error(w, "short_url parameter is required", http.StatusBadRequest)
			return
		}
		ctx := context.Background()
		// Check Redis Cache
		originalURL, err := rdb.Get(ctx, shortURL).Result()
		if err != redis.Nil {
			fmt.Println("Found in cache")
		} else if err == redis.Nil {
			// Not found in cache, fetch from database
			repo := &repositories.URLRepository{DB: db}
			service := &services.URLService{Repo: repo}
			originalURL, err = service.ExpandURL(shortURL)
			if err != nil {
				http.Error(w, "URL not found", http.StatusNotFound)
				return
			}
			// Cache the result in Redis
			rdb.Set(ctx, shortURL, originalURL, 24*time.Hour)
			fmt.Println("Not found in cache")
		} else if err != nil {
			http.Error(w, "Failed to fetch URL", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}