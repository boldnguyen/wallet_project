package handlers

import (
	"context"
	"fmt"
	"time"
	"wallet_project/models"
	"wallet_project/services/bet"

	"gorm.io/gorm"
)

// PlaceBetHandler handles a new bet request.
func PlaceBetHandler(ctx context.Context, betService *bet.BetService, db *gorm.DB, playerID string, betType models.BetType, amount float64, selection string) (uint, error) {
	fmt.Println("\nPlacing Bet...")

	// Log để kiểm tra giá trị wallet_address
	fmt.Printf("Checking user with wallet address: %s\n", playerID) // playerID có thể là wallet_address nếu cần thiết

	// Kiểm tra sự tồn tại của người chơi
	var user models.User
	if err := db.Where("wallet_address = ?", playerID).First(&user).Error; err != nil {
		return 0, fmt.Errorf("user not found for wallet address %s: %v", playerID, err) // Thêm log chi tiết
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

	// Thực hiện các thao tác bổ sung với betService nếu cần
	betID, err := betService.PlaceBet(ctx, db, playerID, betType, amount, selection)
	if err != nil {
		return 0, fmt.Errorf("failed to place bet: %v", err)
	}

	return betID, nil
}
