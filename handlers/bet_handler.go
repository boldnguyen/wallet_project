package handlers

import (
	"context"
	"fmt"
	"wallet_project/models"
	"wallet_project/services/bet"
)

func PlaceBetHandler(ctx context.Context, betService *bet.BetService, playerID string, betType models.BetType, betAmount float64, selection string) (string, error) {
	fmt.Println("\nPlacing Bet...")
	betID, err := betService.PlaceBet(ctx, playerID, betType, betAmount, selection)
	if err != nil {
		return "", fmt.Errorf("failed to place bet: %v", err)
	}
	return betID, nil
}

func ProcessBetsHandler(ctx context.Context, spinService *bet.SpinService, betService *bet.BetService, playerID string, spinID int) error {
	// Process bets for the current spin
	spinResult, err := spinService.Spin()
	if err != nil {
		return fmt.Errorf("failed to process spin: %v", err)
	}
	fmt.Printf("Spin Result: Number: %d, Color: %s, Parity: %s, Group: %s\n",
		spinResult.Number, spinResult.Color, spinResult.Parity, spinResult.Group)

	if err := spinService.ProcessBets(spinID); err != nil {
		return fmt.Errorf("failed to process bets: %v", err)
	}

	// Get and display bet results for the player
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
	return nil
}
