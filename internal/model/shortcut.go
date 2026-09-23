package model

import "time"

// Shortcut is a memorable alias that points at a destination URL.
type Shortcut struct {
	ID          uint   `gorm:"primaryKey"`
	Alias       string `gorm:"uniqueIndex;not null"`
	Destination string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Shortcut) TableName() string {
	return "shortcuts"
}
