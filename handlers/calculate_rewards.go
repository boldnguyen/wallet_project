package handlers

import (
	"context"
	"fmt"
	"wallet_project/models"

	"gorm.io/gorm"
)

// CalculateRewardsHandler tính toán phần thưởng và cập nhật kết quả cược
func CalculateRewardsHandler(ctx context.Context, db *gorm.DB, spinID int) error {
	// Truy vấn kết quả vòng quay
	var spinResult models.SpinResult
	if err := db.Table("spin_results").Where("id = ?", spinID).First(&spinResult).Error; err != nil {
		return fmt.Errorf("failed to find spin result for spin_id %d: %v", spinID, err)
	}

	// Truy vấn các cược liên quan đến spin_id
	var bets []models.Bet
	if err := db.Where("spin_id = ?", spinID).Find(&bets).Error; err != nil {
		return fmt.Errorf("failed to find bets for spin_id %d: %v", spinID, err)
	}

	// Duyệt qua tất cả các cược và xác định thắng thua
	for _, bet := range bets {
		var win bool
		switch bet.BetType {
		case "color":
			// Kiểm tra màu sắc (Đỏ/Đen)
			win = (bet.Selection == spinResult.Color)
		case "parity":
			// Kiểm tra chẵn/lẻ
			win = (bet.Selection == spinResult.Parity)
		case "group":
			// Kiểm tra nhóm (1st12, 2nd12, 3rd12)
			win = (bet.Selection == spinResult.Group)
		case "number":
			// Kiểm tra số
			win = (bet.Selection == fmt.Sprintf("%d", spinResult.Number))
		}

		// Cập nhật trạng thái và tính toán phần thưởng
		if win {
			bet.Status = "won"
			// Tính toán phần thưởng (ví dụ tỷ lệ 1:1, 2:1...)
			switch bet.BetType {
			case "color", "parity":
				bet.Payout = bet.Amount // Tỷ lệ 1:1
			case "group":
				bet.Payout = bet.Amount * 2 // Tỷ lệ 2:1
			case "number":
				bet.Payout = bet.Amount * 35 // Tỷ lệ 35:1
			}
		} else {
			bet.Status = "lost"
			bet.Payout = 0
		}

		// Cập nhật cược trong cơ sở dữ liệu
		if err := db.Save(&bet).Error; err != nil {
			return fmt.Errorf("failed to update bet status: %v", err)
		}

		// Cập nhật lại số dư người chơi nếu thắng
		if win {
			var user models.User
			if err := db.Where("wallet_address = ?", bet.PlayerID).First(&user).Error; err != nil {
				return fmt.Errorf("user not found for wallet address %s: %v", bet.PlayerID, err)
			}

			// Cập nhật số dư của người chơi
			user.Balance += bet.Payout
			if err := db.Save(&user).Error; err != nil {
				return fmt.Errorf("failed to update user balance: %v", err)
			}
		}
	}

	return nil
}
