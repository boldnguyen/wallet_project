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

// PlaceBet chỉ xử lý logic mà không lưu vào cơ sở dữ liệu
func (bs *BetService) PlaceBet(ctx context.Context, db *gorm.DB, walletAddress string, betType models.BetType, amount float64, selection string) (uint, error) {
	// Kiểm tra số dư ví của người chơi từ cơ sở dữ liệu
	var user models.User
	if err := db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
		return 0, fmt.Errorf("user not found for wallet address %s: %v", walletAddress, err)
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

	return 0, nil // Trả về ID giả, vì không lưu trong hàm này
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
