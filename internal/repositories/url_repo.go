package repositories

import (
	"database/sql"
	"url-shortener/internal/models"
)

type URLRepository struct {
	DB *sql.DB
}

func (repo *URLRepository) SaveURL(url models.URL) error {
	// Create the table if it doesn't exist
	_, err := repo.DB.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			original_url TEXT NOT NULL,
			short_url TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		return err
	}
	// Insert the URL record into the table
	_, err = repo.DB.Exec("INSERT INTO urls (original_url, short_url) VALUES ($1, $2)", url.OriginalURL, url.ShortURL)
	return err
}

func (repo *URLRepository) GetOriginalURL(shortURL string) (string, error) {
	var originalURL string
	err := repo.DB.QueryRow("select original_url from urls where short_url = $1", shortURL).Scan(&originalURL)
	return originalURL, err
}