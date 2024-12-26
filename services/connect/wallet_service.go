package connect

import (
	"context"
	"fmt"
	"wallet_project/models"

	"github.com/ethereum/go-ethereum/common"
)

// WalletService struct represents the service that interacts with the wallet.
type WalletService struct{}

// NewWalletService creates a new instance of WalletService.
func NewWalletService() *WalletService {
	return &WalletService{}
}

// ConnectWallet connects to the wallet by verifying the wallet address.
func (ws *WalletService) ConnectWallet(ctx context.Context, walletAddress string) error {
	if !common.IsHexAddress(walletAddress) {
		return fmt.Errorf("invalid wallet address")
	}
	fmt.Printf("Wallet %s connected successfully!\n", walletAddress)
	return nil
}

// UpdatePlayerBalance updates the balance of the player's wallet.
func (ws *WalletService) UpdatePlayerBalance(playerID string, amount float64) error {
	// In a real implementation, you'd fetch the wallet by playerID from a database.
	// For now, let's assume a mock wallet for simplicity.
	mockWallet := &models.Wallet{
		Address: playerID,
		Balance: 1000.0, // Mock initial balance
	}

	mockWallet.Balance += amount

	if mockWallet.Balance < 0 {
		return fmt.Errorf("insufficient balance for player %s", playerID)
	}

	fmt.Printf("Player %s balance updated to %.2f\n", playerID, mockWallet.Balance)
	return nil
}

// WithdrawFunds handles the withdrawal of funds to a player's wallet.
func (ws *WalletService) WithdrawFunds(ctx context.Context, playerID string, amount float64) error {
	// Simulate fetching the player's wallet.
	mockWallet := &models.Wallet{
		Address: playerID,
		Balance: 1000.0, // Mock initial balance
	}

	// Check if the wallet has enough balance.
	if mockWallet.Balance < amount {
		return fmt.Errorf("insufficient balance to withdraw %.2f for player %s", amount, playerID)
	}

	// Deduct the amount from the player's wallet.
	mockWallet.Balance -= amount

	// Simulate sending a withdrawal transaction to the blockchain.
	fmt.Printf("Initiating withdrawal of %.2f to wallet address: %s\n", amount, mockWallet.Address)
	err := ws.sendWithdrawalTransaction(mockWallet.Address, amount)
	if err != nil {
		// Revert balance deduction on failure.
		mockWallet.Balance += amount
		return fmt.Errorf("failed to withdraw funds: %v", err)
	}

	fmt.Printf("Withdrawal of %.2f successful for player %s. Remaining balance: %.2f\n", amount, playerID, mockWallet.Balance)
	return nil
}

// sendWithdrawalTransaction simulates sending a withdrawal transaction to a smart contract.
func (ws *WalletService) sendWithdrawalTransaction(walletAddress string, amount float64) error {
	// In a real implementation, this function would interact with a smart contract
	// to initiate the withdrawal transaction. Here, we'll simulate success.
	if !common.IsHexAddress(walletAddress) {
		return fmt.Errorf("invalid wallet address")
	}

	// Simulated successful transaction.
	fmt.Printf("Withdrawal transaction of %.2f to %s completed on blockchain.\n", amount, walletAddress)
	return nil
}
