package handlers

import (
	"net/http"
	"wallet_project/services/bet"

	"github.com/gin-gonic/gin"
)

// WithdrawHandler handles withdrawal requests.
func WithdrawHandler(withdrawService *bet.WithdrawService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			PlayerID string  `json:"player_id"`
			Amount   float64 `json:"amount"`
		}

		// Parse JSON body
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		// Process withdrawal
		withdrawRequest, err := withdrawService.ProcessWithdraw(req.PlayerID, req.Amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Respond with success
		c.JSON(http.StatusOK, gin.H{
			"message": "Withdrawal request created successfully",
			"data":    withdrawRequest,
		})
	}
}
