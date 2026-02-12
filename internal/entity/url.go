package entity

import "time"

type Url struct {
	ID          uint64     `gorm:"primary_key;auto_increment" json:"id"`
	OriginalUrl string     `gorm:"type:varchar(255);not null" json:"original_url"`
	ShortUrl    string     `gorm:"type:varchar(8);not null;unique" json:"short_url"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"default:NULL" json:"deleted_at"`
}
