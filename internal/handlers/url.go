package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"url-shortener/internal/repositories"
	"url-shortener/internal/services"
)

func ShortenURLHandler(db *sql.DB) http.HandlerFunc {
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
		repo := &repositories.URLRepository{DB: db}
		service := &services.URLService{Repo: repo}
		url, err := service.ShortenURL(req.OriginalURL)
		if err != nil {
			http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
			return
		}
		json.NewDecoder(w).Encode(url)
	}
}

func ExpandURLHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		shortURL := r.URL.Query().Get("short_url")
		if shortURL == "" {
			http.Error(w, "short_url parameter is required", http.StatusBadRequest)
			return
		}
		repo := &repositories.URLRepository{DB: db}
		service := &services.URLService{Repo: repo}
		originalURL, err := service.ExpandURL(shortURL)
		if err != nill {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		http.Redirect(w, r, originalUrl, http.StatusFound)
	}
}