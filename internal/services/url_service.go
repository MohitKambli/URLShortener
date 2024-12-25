package services

import (
	"url-shortener/internal/models"
	"url-shortener/internal/repositories"
	"url-shortener/pkg/utils"
)

type URLService struct {
	Repo *repositories.URLRepository
}

func (service *URLService) ShortenURL(originalURL string) (models.URL, error) {
	shortURL := utils.GenerateShortURL()
	url := models.URL {
		OriginalURL: originalURL,
		ShortURL: shortURL
	}
	err := service.Repo.SaveURL(url)
	return url, err
}