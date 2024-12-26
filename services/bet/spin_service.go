package bet

import (
	"fmt"
	"math/rand"
	"time"
	"wallet_project/models"
	"wallet_project/services/connect"
)

// SpinService handles the roulette spin logic.
type SpinService struct {
	spins []models.SpinResult    // Store spin results
	bets  map[int][]models.Bet   // Map spin ID to bets
	ws    *connect.WalletService // Reference to WalletService for balance updates
}

// NewSpinService creates a new instance of SpinService.
func NewSpinService(walletService *connect.WalletService) *SpinService {
	return &SpinService{
		spins: []models.SpinResult{},
		bets:  make(map[int][]models.Bet),
		ws:    walletService,
	}
}

// Spin simulates a roulette spin and returns the result.
func (ss *SpinService) Spin() (*models.SpinResult, error) {
	// Simulate a spin (number from 0 to 36)
	rand.Seed(time.Now().UnixNano())
	spinNumber := rand.Intn(37)

	// Determine color (red or black) based on roulette color pattern
	color := "black"
	if spinNumber == 0 {
		color = "green" // Zero is usually green on a roulette wheel
	} else if (spinNumber%2 == 0 && spinNumber <= 10) || (spinNumber%2 == 1 && spinNumber > 10 && spinNumber <= 18) || (spinNumber%2 == 0 && spinNumber > 18 && spinNumber <= 28) || (spinNumber%2 == 1 && spinNumber > 28 && spinNumber <= 36) {
		color = "red"
	}

	// Determine parity (even or odd)
	parity := "even"
	if spinNumber%2 != 0 {
		parity = "odd"
	}

	// Determine the group (1st12, 2nd12, or 3rd12)
	group := ""
	if spinNumber >= 1 && spinNumber <= 12 {
		group = "1st12"
	} else if spinNumber >= 13 && spinNumber <= 24 {
		group = "2nd12"
	} else if spinNumber >= 25 && spinNumber <= 36 {
		group = "3rd12"
	}

	// Create the spin result
	spinResult := &models.SpinResult{
		Number:    spinNumber,
		Color:     color,
		Parity:    parity,
		Group:     group,
		Timestamp: time.Now().Unix(),
	}

	// Store the result
	ss.spins = append(ss.spins, *spinResult)

	return spinResult, nil
}

// GetSpins returns all the spin results.
func (ss *SpinService) GetSpins() []models.SpinResult {
	return ss.spins
}

// ProcessBets calculates winnings for all bets related to a spin result.
func (ss *SpinService) ProcessBets(spinID int) error {
	if spinID >= len(ss.spins) {
		return fmt.Errorf("invalid spin ID")
	}

	// Get the spin result
	spinResult := ss.spins[spinID]
	bets, exists := ss.bets[spinID]
	if !exists || len(bets) == 0 {
		fmt.Printf("No bets found for spin ID %d\n", spinID)
		return nil
	}

	fmt.Printf("Processing bets for spin ID %d, Result: %+v\n", spinID, spinResult)

	// Iterate over each bet and determine win/loss
	for i, bet := range bets {
		won := false
		payout := 0.0

		// Match the bet type with spin result
		switch bet.BetType {
		case models.NumberBet:
			if bet.Selection == fmt.Sprintf("%d", spinResult.Number) {
				won = true
				payout = bet.Amount * 35 // Payout for number bet is 35:1
			}
		case models.ColorBet:
			if bet.Selection == spinResult.Color {
				won = true
				payout = bet.Amount * 2 // Payout for color bet is 1:1
			}
		case models.GroupBet:
			if bet.Selection == spinResult.Group {
				won = true
				payout = bet.Amount * 3 // Payout for group bet is 2:1
			}
		case models.EvenOddBet:
			if bet.Selection == spinResult.Parity {
				won = true
				payout = bet.Amount * 2 // Payout for parity bet is 1:1
			}
		}

		// Update bet status and reward the player
		if won {
			bets[i].Status = "won"
			if err := ss.ws.UpdatePlayerBalance(bet.PlayerID, payout); err != nil {
				return fmt.Errorf("failed to update balance for player %s: %v", bet.PlayerID, err)
			}
			fmt.Printf("Bet %d won! Player %s rewarded with %.2f\n", i, bet.PlayerID, payout)
		} else {
			bets[i].Status = "lost"
			fmt.Printf("Bet %d lost. Player %s lost %.2f\n", i, bet.PlayerID, bet.Amount)
		}
	}

	// Save updated bets back
	ss.bets[spinID] = bets
	return nil
}
