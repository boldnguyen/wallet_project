package handlers

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"
	"wallet_project/models"
	"wallet_project/services/bet"
	"wallet_project/services/connect"

	"github.com/ethereum/go-ethereum/crypto"
	"gorm.io/gorm"
)

// generateFakeWallet creates a fake wallet and returns the address and private key
func generateFakeWallet() (string, *ecdsa.PrivateKey, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", nil, err
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return address, privateKey, nil
}

func ConnectWalletHandler(ctx context.Context, walletService *connect.WalletService, betService *bet.BetService, db *gorm.DB) (string, error) {
	fmt.Println("\nConnecting Wallet...")

	// Tạo ví giả
	walletAddress, _, err := generateFakeWallet()
	if err != nil {
		return "", fmt.Errorf("error generating wallet: %v", err)
	}

	// Kiểm tra nếu địa chỉ ví đã tồn tại
	var existingUser models.User
	if err := db.Where("wallet_address = ?", walletAddress).First(&existingUser).Error; err == nil {
		// Nếu ví đã tồn tại, trả về địa chỉ
		fmt.Printf("Wallet already connected for player: %s\n", existingUser.PlayerID)
		return walletAddress, nil
	} else if err != gorm.ErrRecordNotFound {
		return "", fmt.Errorf("database error: %v", err)
	}

	// Tạo PlayerID động
	playerID := fmt.Sprintf("player_%d", time.Now().UnixNano())

	// Thêm người dùng mới vào cơ sở dữ liệu
	newUser := models.User{
		PlayerID:      playerID,
		WalletAddress: walletAddress,
		Balance:       1000.0,
	}
	if err := db.Create(&newUser).Error; err != nil {
		return "", fmt.Errorf("failed to save user to database: %v", err)
	}

	// Kết nối ví với WalletService
	if err := walletService.ConnectWallet(ctx, walletAddress); err != nil {
		return "", fmt.Errorf("failed to connect wallet: %v", err)
	}
	// Thêm ví vào BetService
	betService.AddPlayerWallet(playerID, newUser.Balance)

	fmt.Printf("New wallet connected with PlayerID: %s\n", playerID)
	return walletAddress, nil
}

// WithdrawFundsHandler handles fund withdrawal
func WithdrawFundsHandler(ctx context.Context, walletService *connect.WalletService, walletAddress string, withdrawAmount float64) error {
	fmt.Printf("Attempting to withdraw %.2f for wallet %s...\n", withdrawAmount, walletAddress)
	err := walletService.WithdrawFunds(ctx, walletAddress, withdrawAmount)
	if err != nil {
		return fmt.Errorf("failed to withdraw funds: %v", err)
	}
	fmt.Printf("Withdrawal successful! %.2f has been deducted from your wallet.\n", withdrawAmount)
	return nil
}
