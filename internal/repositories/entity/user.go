package entity

import (
	"time"
)

type User struct {
	ID          string    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email       string    `gorm:"type:varchar(255);not null"`
	Password    string    `gorm:"type:varchar(255);not null"`
	DisplayName string    `gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Wallets     []Wallet  `gorm:"foreignKey:UserID"`
}
