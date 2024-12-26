package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"wallet_project/models"
	"wallet_project/services/bet"
	"wallet_project/services/connect"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/crypto"
)

var db *gorm.DB

// Initialize database connection and perform migrations
func initDatabase() {
	var err error
	// Define your connection string (DSN) for PostgreSQL
	dsn := "user=postgres password=12345 dbname=postgres port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Check the connection to the database
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database object: %v", err)
	}

	// Ping the database to confirm the connection is alive
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	// Log success message
	fmt.Println("Successfully connected to the database!")

	// Automatically migrate the models, including SpinResult
	if err := db.AutoMigrate(&models.SpinResult{}, &models.User{}, &models.Bet{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

// Function to generate a fake wallet and return the wallet address and private key
func generateFakeWallet() (string, *ecdsa.PrivateKey, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", nil, err
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return address, privateKey, nil
}

func displayMenu() {
	fmt.Println("\n=== Menu ===")
	fmt.Println("1. Connect Wallet")
	fmt.Println("2. Place Bet")
	fmt.Println("3. Spin Roulette")
	fmt.Println("4. Withdraw Funds")
	fmt.Println("5. Exit")
	fmt.Print("Choose an option: ")
}
func main() {

	// Connect to PostgreSQL
	initDatabase()

	ctx := context.Background()
	walletService := connect.NewWalletService()
	betService := bet.NewBetService()
	spinService := bet.NewSpinService(walletService)

	var playerID string
	var walletAddress string
	var balance float64 = 1000.0

	for {
		displayMenu()
		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1: // Connect Wallet
			fmt.Println("\nConnecting Wallet...")
			var err error
			walletAddress, _, err = generateFakeWallet()
			if err != nil {
				log.Printf("Error generating wallet: %v", err)
				continue
			}
			if err := walletService.ConnectWallet(ctx, walletAddress); err != nil {
				log.Printf("Failed to connect wallet: %v", err)
				continue
			}
			playerID = "player1" // For simplicity, using a static player ID
			betService.AddPlayerWallet(playerID, balance)
			fmt.Printf("Wallet %s connected successfully! Balance: %.2f\n", walletAddress, balance)

		case 2: // Place Bet
			if playerID == "" {
				fmt.Println("Please connect a wallet first.")
				continue
			}
			fmt.Println("\nEnter bet details:")
			fmt.Print("Bet type (1: Number, 2: Group, 3: Color, 4: Parity): ")
			var betTypeInput int
			fmt.Scan(&betTypeInput)
			var betType models.BetType
			switch betTypeInput {
			case 1:
				betType = models.NumberBet
			case 2:
				betType = models.GroupBet
			case 3:
				betType = models.ColorBet
			case 4:
				betType = models.EvenOddBet
			default:
				fmt.Println("Invalid bet type.")
				continue
			}

			fmt.Print("Bet amount: ")
			var betAmount float64
			fmt.Scan(&betAmount)

			fmt.Print("Selection: ")
			var selection string
			fmt.Scan(&selection)

			betID, err := betService.PlaceBet(ctx, playerID, betType, betAmount, selection)
			if err != nil {
				log.Printf("Failed to place bet: %v", err)
				continue
			}
			fmt.Printf("Bet placed successfully with ID: %s\n", betID)

		case 3: // Spin Roulette
			fmt.Println("\nSpinning Roulette...")
			spinResult, err := spinService.Spin()
			if err != nil {
				log.Printf("Failed to process spin: %v", err)
				continue
			}
			fmt.Printf("Spin Result: Number: %d, Color: %s, Parity: %s, Group: %s\n",
				spinResult.Number, spinResult.Color, spinResult.Parity, spinResult.Group)

			spinID := len(spinService.GetSpins()) - 1
			if err := spinService.ProcessBets(spinID); err != nil {
				log.Printf("Failed to process bets: %v", err)
				continue
			}

			// Get and display bet results
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

			// Show updated balance
			updatedBalance, err := betService.GetPlayerBalance(playerID)
			if err != nil {
				log.Printf("Failed to get updated balance: %v", err)
			} else {
				fmt.Printf("Updated balance for %s: %.2f\n", playerID, updatedBalance)
			}

		case 4: // Withdraw Funds
			if playerID == "" {
				fmt.Println("Please connect a wallet first.")
				continue
			}
			fmt.Print("Enter amount to withdraw: ")
			var withdrawAmount float64
			fmt.Scan(&withdrawAmount)

			fmt.Printf("Attempting to withdraw %.2f for player %s...\n", withdrawAmount, playerID)
			err := walletService.WithdrawFunds(ctx, walletAddress, withdrawAmount)
			if err != nil {
				log.Printf("Failed to withdraw funds: %v", err)
			} else {
				fmt.Printf("Withdrawal successful! %.2f has been deducted from your wallet.\n", withdrawAmount)
			}

		case 5: // Exit
			fmt.Println("Exiting program. Goodbye!")
			return

		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
