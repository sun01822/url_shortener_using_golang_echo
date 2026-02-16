package entity

import (
	"time"

	"gorm.io/gorm"
)

type Url struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	OriginalUrl string         `gorm:"type:varchar(255);not null" json:"original_url"`
	ShortUrl    string         `gorm:"type:varchar(8);not null;uniqueIndex" json:"short_url"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
