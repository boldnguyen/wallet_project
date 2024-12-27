package services

import (
	"wallet_project/models"

	"gorm.io/gorm"
)

// SaveUserToDatabase saves a user to the database
func SaveUserToDatabase(db *gorm.DB, playerID string, walletAddress string, balance float64) error {
	user := models.User{
		PlayerID:      playerID,
		WalletAddress: walletAddress,
		Balance:       balance,
	}
	return db.Create(&user).Error
}
