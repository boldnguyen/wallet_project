package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"wallet_project/models"
	"wallet_project/services/bet"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PlaceBetHandler handles a new bet request.
func PlaceBetHandler(ctx context.Context, betService *bet.BetService, db *gorm.DB, playerID string, betType models.BetType, amount float64, selection string) (uint, error) {
	fmt.Println("\nPlacing Bet...")

	// Kiểm tra sự tồn tại của người chơi
	var user models.User
	if err := db.Where("wallet_address = ?", playerID).First(&user).Error; err != nil {
		return 0, fmt.Errorf("user not found for wallet address %s: %v", playerID, err)
	}

	// Tạo thông tin đặt cược
	newBet := models.Bet{
		PlayerID:  playerID,
		BetType:   betType,
		Amount:    amount,
		Selection: selection,
		Status:    "placed",
		Timestamp: time.Now().Unix(), // Chuyển đổi time.Time thành Unix timestamp (int64)
	}

	// Lưu vào cơ sở dữ liệu
	if err := db.Create(&newBet).Error; err != nil {
		return 0, fmt.Errorf("failed to save bet to database: %v", err)
	}

	// Thực hiện kiểm tra logic bổ sung từ betService
	if _, err := betService.PlaceBet(ctx, db, playerID, betType, amount, selection); err != nil {
		return 0, fmt.Errorf("failed to place bet: %v", err)
	}

	return newBet.ID, nil
}

// CancelBetHandler handles bet cancellation requests.
func CancelBetHandler(betService *bet.BetService, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			BetID uint `json:"bet_id"`
		}

		// Parse the request JSON
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Attempt to cancel the bet
		if err := betService.CancelBet(db, request.BetID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Bet canceled successfully"})
	}
}
