package repositories

import (
	"database/sql"
	"url-shortener/internal.models"
)

type URLRepository struct {
	DB *sql.DB
}

func (repo *URLRepository) SaveURL(url models.URL) error {
	_, err := repo.DB.Exec("insert into urls (original_url, short_url) values ($1, $2)", url.OriginalURL, url.ShortURL)
	return err
}

func (repo *URLRepository) GetOriginalURL(shortURL string) (string, error) {
	var originalURL string
	err := repo.DB.QueryRow("select original_url from urls where short_url = $1", shortURL).Scan(&originalURL)
	return originalURL, err
}