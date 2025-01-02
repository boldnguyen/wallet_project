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

// CalculateAndDistributeRewards calculates rewards for bets and updates user balances.
func (bs *BetService) CalculateAndDistributeRewards(ctx context.Context, db *gorm.DB, spinResult models.SpinResult) error {
	fmt.Println("\nCalculating and distributing rewards...")

	// Lấy tất cả các cược liên quan đến spin_id hiện tại
	var bets []models.Bet
	if err := db.Where("spin_id = ?", spinResult.ID).Find(&bets).Error; err != nil {
		return fmt.Errorf("failed to retrieve bets for spin_id %d: %v", spinResult.ID, err)
	}

	for _, bet := range bets {
		// Tính toán phần thưởng dựa trên loại cược
		var payout float64
		var won bool

		switch bet.BetType {
		case "color":
			if spinResult.Color == bet.Selection {
				payout = bet.Amount * 2 // Tỷ lệ 1:1
				won = true
			}
		case "group":
			if spinResult.Group == bet.Selection {
				payout = bet.Amount * 3 // Tỷ lệ 2:1
				won = true
			}
		case "parity":
			if spinResult.Parity == bet.Selection {
				payout = bet.Amount * 2 // Tỷ lệ 1:1
				won = true
			}
		case "number":
			if fmt.Sprintf("%d", spinResult.Number) == bet.Selection {
				payout = bet.Amount * 36 // Tỷ lệ 35:1
				won = true
			}
		default:
			continue
		}

		// Cập nhật trạng thái cược
		if won {
			bet.Status = "won"
			bet.Payout = payout
		} else {
			bet.Status = "lost"
			bet.Payout = 0
		}

		// Lưu thông tin cược cập nhật
		if err := db.Save(&bet).Error; err != nil {
			return fmt.Errorf("failed to update bet %d: %v", bet.ID, err)
		}

		// Nếu thắng, cập nhật số dư người chơi
		if won {
			var user models.User
			if err := db.Where("player_id = ?", bet.PlayerID).First(&user).Error; err != nil {
				return fmt.Errorf("failed to find user with player_id %s: %v", bet.PlayerID, err)
			}

			user.Balance += payout
			if err := db.Save(&user).Error; err != nil {
				return fmt.Errorf("failed to update balance for player_id %s: %v", bet.PlayerID, err)
			}
		}
	}

	return nil
}

// CancelBet cancels an existing bet if possible and refunds the amount to the user's wallet.
func (bs *BetService) CancelBet(db *gorm.DB, betID uint) error {
	// Retrieve the bet
	var bet models.Bet
	if err := db.First(&bet, betID).Error; err != nil {
		return fmt.Errorf("bet not found: %v", err)
	}

	// Ensure the bet is in a cancellable state
	if bet.Status != "placed" {
		return fmt.Errorf("bet cannot be canceled as it is in '%s' state", bet.Status)
	}

	// Retrieve the user
	var user models.User
	if err := db.Where("player_id = ?", bet.PlayerID).First(&user).Error; err != nil {
		return fmt.Errorf("user not found for player_id %s: %v", bet.PlayerID, err)
	}

	// Refund the bet amount to the user's balance
	user.Balance += bet.Amount
	if err := db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to refund user balance: %v", err)
	}

	// Update the bet status to "canceled"
	bet.Status = "canceled"
	if err := db.Save(&bet).Error; err != nil {
		return fmt.Errorf("failed to update bet status: %v", err)
	}

	return nil
}
