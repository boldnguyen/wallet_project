package handlers

import (
	"fmt"
	"wallet_project/models" // Thêm import models
	"wallet_project/services/bet"
)

func SpinRouletteHandler(spinService *bet.SpinService) (*models.SpinResult, error) { // Sử dụng models.SpinResult
	fmt.Println("\nSpinning Roulette...")
	spinResult, err := spinService.Spin()
	if err != nil {
		return nil, fmt.Errorf("failed to process spin: %v", err)
	}
	return spinResult, nil
}
