package entity

import (
	"time"
)

type Wallet struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      string    `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Description *string   `gorm:"type:varchar(255)"`
	Balance     float64   `gorm:"type:decimal(20,8);not null;default:0"`
	CreatedAt   time.Time `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt   time.Time `gorm:"type:timestamp;not null;default:now()"`
}
