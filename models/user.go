package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	PlayerID      string `gorm:"unique;not null"`
	WalletAddress string `gorm:"unique;not null"`
	Balance       float64
}
