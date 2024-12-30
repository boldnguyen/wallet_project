package bet

import (
	"context"
	"fmt"
	"wallet_project/models"

	"gorm.io/gorm"
)

// BetService handles betting logic and maintains player wallets.
type BetService struct {
	wallets map[string]*models.Wallet // Map of player ID to Wallet
	bets    []models.Bet              // Slice to store all bets
}

// NewBetService creates a new instance of BetService.
func NewBetService() *BetService {
	return &BetService{
		wallets: make(map[string]*models.Wallet),
		bets:    []models.Bet{},
	}
}

// AddPlayerWallet adds or initializes a player's wallet with a given balance.
func (bs *BetService) AddPlayerWallet(playerID string, initialBalance float64) {
	bs.wallets[playerID] = &models.Wallet{
		Address: playerID,
		Balance: initialBalance,
	}
}

// GetPlayerBalance returns the current balance of the player's wallet.
func (bs *BetService) GetPlayerBalance(playerID string) (float64, error) {
	wallet, exists := bs.wallets[playerID]
	if !exists {
		return 0, fmt.Errorf("player wallet not found")
	}
	return wallet.Balance, nil
}

// PlaceBet places a bet for a player (only handles the business logic, no DB saving)
func (bs *BetService) PlaceBet(ctx context.Context, db *gorm.DB, playerID string, betType models.BetType, amount float64, selection string) (uint, error) {
	// Kiểm tra số dư ví của người chơi từ cơ sở dữ liệu
	var user models.User
	if err := db.Where("player_id = ?", playerID).First(&user).Error; err != nil {
		return 0, fmt.Errorf("player wallet not found: %v", err)
	}

	// Kiểm tra xem người chơi có đủ số dư để đặt cược không
	if user.Balance < amount {
		return 0, fmt.Errorf("insufficient balance")
	}

	// Giảm số dư của người chơi sau khi đặt cược
	user.Balance -= amount
	if err := db.Save(&user).Error; err != nil {
		return 0, fmt.Errorf("failed to update user balance: %v", err)
	}

	// Trả về ID của cược (hiện tại bạn không cần lưu cược ở đây)
	return 0, nil // Không lưu cược vào DB ở đây, chỉ cần trả về ID từ bet_handler.go
}

// GetPlayerBets retrieves all bets for a specific player for the latest spin ID.
func (bs *BetService) GetPlayerBets(playerID string, spinID int) []models.Bet {
	var playerBets []models.Bet
	for _, bet := range bs.bets {
		if bet.PlayerID == playerID {
			playerBets = append(playerBets, bet)
		}
	}
	return playerBets
}
