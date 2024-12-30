package handlers

import (
	"context"
	"fmt"
	"time"
	"wallet_project/models"
	"wallet_project/services/bet"

	"gorm.io/gorm"
)

func PlaceBetHandler(ctx context.Context, betService *bet.BetService, db *gorm.DB, walletAddress string, betType models.BetType, amount float64, selection string) (uint, error) {
	fmt.Println("\nPlacing Bet...")
	fmt.Println("Wallet Address: ", walletAddress)

	// Tìm người chơi từ wallet_address
	var user models.User
	if err := db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		fmt.Println("Database error:", err)
		return 0, fmt.Errorf("player wallet not found: %v", err)
	}

	// Kiểm tra số dư ví của người chơi
	if user.Balance < amount {
		fmt.Println("Insufficient balance:", user.Balance)
		return 0, fmt.Errorf("insufficient balance")
	}

	// Cập nhật lại số dư ví sau khi đặt cược
	user.Balance -= amount
	if err := db.Save(&user).Error; err != nil {
		return 0, fmt.Errorf("failed to update user balance: %v", err)
	}

	// Tạo thông tin cược
	newBet := models.Bet{
		PlayerID:  user.PlayerID,
		BetType:   betType,
		Amount:    amount,
		Selection: selection,
		Status:    "placed",
		Timestamp: time.Now().Unix(),
	}

	// Lưu cược vào cơ sở dữ liệu
	if err := db.Create(&newBet).Error; err != nil {
		fmt.Println("Error saving bet:", err)
		return 0, fmt.Errorf("failed to save bet to database: %v", err)
	}

	// Trả về ID của cược đã lưu vào cơ sở dữ liệu
	return newBet.ID, nil
}

func ProcessBetsHandler(ctx context.Context, spinService *bet.SpinService, betService *bet.BetService, playerID string, spinID int) error {
	// Process bets for the current spin
	spinResult, err := spinService.Spin()
	if err != nil {
		return fmt.Errorf("failed to process spin: %v", err)
	}
	fmt.Printf("Spin Result: Number: %d, Color: %s, Parity: %s, Group: %s\n",
		spinResult.Number, spinResult.Color, spinResult.Parity, spinResult.Group)

	if err := spinService.ProcessBets(spinID); err != nil {
		return fmt.Errorf("failed to process bets: %v", err)
	}

	// Get and display bet results for the player
	playerBets := betService.GetPlayerBets(playerID, spinID)
	if len(playerBets) == 0 {
		fmt.Println("No bets found for this spin.")
	} else {
		fmt.Printf("Results for player %s:\n", playerID)
		for _, bet := range playerBets {
			if bet.Status == "won" {
				fmt.Printf("Bet Type: %s, Selection: %s, Amount: %.2f - You WON!\n", bet.BetType, bet.Selection, bet.Amount)
			} else if bet.Status == "lost" {
				fmt.Printf("Bet Type: %s, Selection: %s, Amount: %.2f - You LOST.\n", bet.BetType, bet.Selection, bet.Amount)
			}
		}
	}
	return nil
}
