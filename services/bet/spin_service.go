package bet

import (
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// SpinResult represents the result of a roulette spin
type SpinResult struct {
	Number int    `json:"number"`
	Color  string `json:"color"`
	Parity string `json:"parity"`
	Group  string `json:"group"`
}

// SpinService is responsible for handling roulette spin logic
type SpinService struct {
	db *gorm.DB
}

// NewSpinService creates a new SpinService
func NewSpinService(db *gorm.DB) *SpinService {
	return &SpinService{db: db}
}

// Spin performs the roulette spin and saves the result in the database
func (s *SpinService) Spin() (*SpinResult, error) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 36 (inclusive)
	number := rand.Intn(37)

	// Determine the color of the number
	color := "red" // Default color for simplicity (can be enhanced)
	if number%2 == 0 {
		color = "black"
	}

	// Determine if the number is even or odd
	parity := "even"
	if number%2 != 0 {
		parity = "odd"
	}

	// Determine the group (1st 12, 2nd 12, or 3rd 12)
	group := "1st12"
	if number > 12 && number <= 24 {
		group = "2nd12"
	} else if number > 24 {
		group = "3rd12"
	}

	// Create SpinResult object
	spinResult := &SpinResult{
		Number: number,
		Color:  color,
		Parity: parity,
		Group:  group,
	}

	// Save the result into the database
	if err := s.db.Create(&spinResult).Error; err != nil {
		return nil, err
	}

	return spinResult, nil
}
