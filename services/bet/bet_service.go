package bet

import (
	"context"
	"fmt"
	"time"
	"wallet_project/models"
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

// PlaceBet places a bet for a player.
func (bs *BetService) PlaceBet(ctx context.Context, playerID string, betType models.BetType, amount float64, selection string) (string, error) {
	// Check if player exists
	wallet, exists := bs.wallets[playerID]
	if !exists {
		return "", fmt.Errorf("player wallet not found")
	}

	// Check if the player has enough balance
	if wallet.Balance < amount {
		return "", fmt.Errorf("insufficient balance")
	}

	// Deduct the bet amount from the player's balance
	wallet.Balance -= amount

	// Create and store the bet
	bet := models.Bet{
		PlayerID:  playerID,
		BetType:   betType,
		Amount:    amount,
		Selection: selection,
		Status:    "placed",
		Timestamp: time.Now().Unix(),
	}
	bs.bets = append(bs.bets, bet)

	// Generate a bet ID (for simplicity, use the timestamp)
	betID := fmt.Sprintf("%d", bet.Timestamp)
	return betID, nil
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
