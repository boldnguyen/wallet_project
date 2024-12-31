package handlers

import (
	"net/http"
	"wallet_project/services/bet"

	"github.com/gin-gonic/gin"
)

// SpinRouletteHandler handles the spin roulette request.
func SpinRouletteHandler(spinService *bet.SpinService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Spin the roulette and get the result
		spinResult, err := spinService.SpinRoulette()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Respond with the spin result
		c.JSON(http.StatusOK, gin.H{
			"number": spinResult.Number,
			"color":  spinResult.Color,
			"parity": spinResult.Parity,
			"group":  spinResult.Group,
		})
	}
}
