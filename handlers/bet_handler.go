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
