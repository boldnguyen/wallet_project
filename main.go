package main

import (
	"context"
	"fmt"
	"log"
	"wallet_project/handlers"
	"wallet_project/models"
	"wallet_project/services/bet"
	"wallet_project/services/connect"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// Initialize database connection and perform migrations
func initDatabase() {
	var err error
	dsn := "user=postgres password=12345 dbname=postgres port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database object: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	if err := db.AutoMigrate(&models.SpinResult{}, &models.User{}, &models.Bet{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

func main() {
	// Connect to PostgreSQL
	initDatabase()

	// Initialize services
	ctx := context.Background()
	walletService := connect.NewWalletService()
	betService := bet.NewBetService()
	spinService := bet.NewSpinService(walletService)

	router := gin.Default()

	// API endpoint to connect wallet
	router.POST("/connect_wallet", func(c *gin.Context) {
		walletAddress, err := handlers.ConnectWalletHandler(ctx, walletService, db)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Lấy PlayerID động từ `ConnectWalletHandler`
		var user models.User
		if err := db.Where("wallet_address = ?", walletAddress).First(&user).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve user from database"})
			return
		}

		c.JSON(200, gin.H{"wallet_address": walletAddress, "player_id": user.PlayerID, "balance": user.Balance})
	})

	// API endpoint to place bet
	router.POST("/place_bet", func(c *gin.Context) {
		var betRequest struct {
			WalletAddress string         `json:"wallet_address"` // Thêm WalletAddress vào body
			BetType       models.BetType `json:"bet_type"`
			BetAmount     float64        `json:"bet_amount"`
			Selection     string         `json:"selection"`
		}
		if err := c.ShouldBindJSON(&betRequest); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Call PlaceBetHandler với wallet_address từ body
		betID, err := handlers.PlaceBetHandler(ctx, betService, db, betRequest.WalletAddress, betRequest.BetType, betRequest.BetAmount, betRequest.Selection)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"bet_id": betID})
	})

	// API endpoint to spin roulette
	router.GET("/spin_roulette", func(c *gin.Context) {
		spinResult, err := handlers.SpinRouletteHandler(spinService)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"number": spinResult.Number,
			"color":  spinResult.Color,
			"parity": spinResult.Parity,
			"group":  spinResult.Group,
		})
	})

	// Run the server
	router.Run(":8080") // Port 8080
}
