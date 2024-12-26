package handlers

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"wallet_project/services/connect"

	"github.com/ethereum/go-ethereum/crypto"
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

// ConnectWalletHandler handles wallet connection
func ConnectWalletHandler(ctx context.Context, walletService *connect.WalletService) (string, error) {
	fmt.Println("\nConnecting Wallet...")
	walletAddress, _, err := generateFakeWallet()
	if err != nil {
		return "", fmt.Errorf("error generating wallet: %v", err)
	}

	if err := walletService.ConnectWallet(ctx, walletAddress); err != nil {
		return "", fmt.Errorf("failed to connect wallet: %v", err)
	}
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
