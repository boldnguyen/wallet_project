package handlers

import (
	"net/http"
	"wallet_project/services/bet"

	"github.com/gin-gonic/gin"
)

// SpinRouletteHandler handles the roulette spin request.
func SpinRouletteHandler(spinService *bet.SpinService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Thực hiện vòng quay và lưu kết quả vào cơ sở dữ liệu
		spinResult, err := spinService.Spin()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Trả về kết quả spin
		c.JSON(http.StatusOK, gin.H{
			"number": spinResult.Number,
			"color":  spinResult.Color,
			"parity": spinResult.Parity,
			"group":  spinResult.Group,
		})
	}
}
