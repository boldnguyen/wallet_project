package bet

import (
	"fmt"
	"math/rand"
	"time"
	"wallet_project/models"

	"gorm.io/gorm"
)

// SpinService handles the logic for spinning the roulette.
type SpinService struct {
	db *gorm.DB
}

// NewSpinService creates a new SpinService.
func NewSpinService(db *gorm.DB) *SpinService {
	return &SpinService{db: db}
}

// SpinRoulette generates a random spin result and saves it to the database.
func (ss *SpinService) SpinRoulette() (*models.SpinResult, error) {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 36
	number := rand.Intn(37)

	// Determine color
	color := "black"
	if number == 0 {
		color = "green"
	} else if number%2 == 1 {
		color = "red"
	}

	// Determine parity (odd/even)
	parity := "even"
	if number%2 != 0 {
		parity = "odd"
	}

	// Determine group (1st12, 2nd12, 3rd12)
	group := ""
	switch {
	case number >= 1 && number <= 12:
		group = "1st12"
	case number >= 13 && number <= 24:
		group = "2nd12"
	case number >= 25 && number <= 36:
		group = "3rd12"
	default:
		group = "none"
	}

	// Create the spin result
	spinResult := &models.SpinResult{
		Number:    number,
		Color:     color,
		Parity:    parity,
		Group:     group,
		Timestamp: time.Now().Unix(),
	}

	// Save the spin result to the database
	if err := ss.db.Create(spinResult).Error; err != nil {
		return nil, fmt.Errorf("failed to save spin result: %v", err)
	}

	return spinResult, nil
}
