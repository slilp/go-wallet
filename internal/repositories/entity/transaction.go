package entity

import (
	"time"
)

type Transaction struct {
	ID        string    `gorm:"type:varchar(20);primaryKey"`
	From      *string   `gorm:"type:uuid;index"`
	To        *string   `gorm:"type:uuid;index"`
	Amount    float64   `gorm:"type:decimal(20,2);not null"`
	Type      string    `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:now()"`
}
