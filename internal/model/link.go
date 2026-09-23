package model

import "time"

type Link struct {
	ID        uint   `gorm:"primaryKey"`
	Slug      string `gorm:"uniqueIndex;not null"`
	URL       string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
