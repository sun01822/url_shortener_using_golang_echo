package repository

import (
	"url_shortener/internal/entity"

	"gorm.io/gorm"
)

type Repository interface {
	// Create creates a new short URL for the given original URL.
	Create(url *entity.Url) (*entity.Url, error)

	// Get retrieves the original URL associated with the given short URL.
	Get(shortCode string) (*entity.Url, error)

	// Delete deletes the short URL and its associated original URL.
	Delete()

	// Update updates the original URL associated with the given short URL.
	Update()
}

type urlRepository struct {
	db *gorm.DB
}

func NewUrlRepository(db *gorm.DB) *urlRepository {
	return &urlRepository{
		db: db,
	}
}

func (r *urlRepository) Create(url *entity.Url) (*entity.Url, error) {
	if err := r.db.Create(&url).Error; err != nil {
		return &entity.Url{}, err
	}
	return url, nil
}

func (r *urlRepository) Get(shortCode string) (*entity.Url, error) {
	var url entity.Url
	err := r.db.Where("short_url = ?", shortCode).First(&url).Error
	if err != nil {
		return &entity.Url{}, err
	}
	return &url, nil
}

func (r *urlRepository) Delete() {

}

func (r *urlRepository) Update() {

}
