package service

import (
	"crypto/sha256"
	"github.com/jxskiss/base62"
	"url_shortener/internal/entity"
	"url_shortener/internal/models"
	"url_shortener/internal/repository"

	"github.com/labstack/gommon/log"
)

type Service interface {
	// CreateShortUrl creates a new short URL for the given original URL.
	CreateShortUrl(request models.CreateShortUrlRequest) (models.CreateShortUrlResponse, error)

	// GetOriginalUrl retrieves the original URL associated with the given short URL.
	GetOriginalUrl(shortCode string) (*entity.Url, error)

	// DeleteShortUrl deletes the short URL and its associated original URL.
	DeleteShortUrl()

	// UpdateShortUrl updates the original URL associated with the given short URL.
	UpdateShortUrl()

	// ListShortUrls lists all short URLs and their associated original URLs.
	ListShortUrls()
}

type urlService struct {
	repo repository.Repository
}

func NewUrlService(repo repository.Repository) *urlService {
	return &urlService{
		repo: repo,
	}
}

func (s *urlService) CreateShortUrl(request models.CreateShortUrlRequest) (models.CreateShortUrlResponse, error) {
	url := entity.Url{
		OriginalUrl: request.OriginalUrl,
	}

	if request.CustomShortUrl != "" {
		url.ShortUrl = request.CustomShortUrl
	} else {
		url.ShortUrl = generateShortUrl(request.OriginalUrl)
	}

	resp, err := s.repo.Create(&url)
	if err != nil {
		log.Error(err.Error())
		return models.CreateShortUrlResponse{}, err
	}

	return models.CreateShortUrlResponse{
		ShortUrl:  resp.ShortUrl,
		Status:    "success",
		CreatedBy: "system",
		CreatedAt: resp.CreatedAt,
	}, nil
}

func (s *urlService) GetOriginalUrl(shortCode string) (*entity.Url, error) {
	return s.repo.Get(shortCode)
}

func (s *urlService) DeleteShortUrl() {

}

func (s *urlService) UpdateShortUrl() {

}

func (s *urlService) ListShortUrls() {

}

func generateShortUrl(url string) string {
	hash := sha256.Sum256([]byte(url))
	encoded := base62.EncodeToString(hash[:]) // encode hash in Base62
	return encoded[:8]
}
