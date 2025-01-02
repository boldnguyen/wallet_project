package bet

import (
	"errors"
	"fmt"
	"time"
	"wallet_project/models"

	"gorm.io/gorm"
)

// WithdrawService handles withdrawal operations.
type WithdrawService struct {
	db *gorm.DB
}

// NewWithdrawService creates a new WithdrawService.
func NewWithdrawService(db *gorm.DB) *WithdrawService {
	return &WithdrawService{db: db}
}

// ProcessWithdraw processes a withdrawal request.
func (ws *WithdrawService) ProcessWithdraw(playerID string, amount float64) (*models.WithdrawRequest, error) {
	// Tìm người chơi
	var user models.User
	if err := ws.db.Where("wallet_address = ?", playerID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("player not found: %v", err)
	}

	// Kiểm tra số dư
	if user.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	// Tạo yêu cầu rút tiền
	withdrawRequest := &models.WithdrawRequest{
		PlayerID:    playerID,
		Amount:      amount,
		Status:      "pending",
		RequestedAt: time.Now(),
	}

	// Lưu vào cơ sở dữ liệu
	if err := ws.db.Create(withdrawRequest).Error; err != nil {
		return nil, fmt.Errorf("failed to create withdrawal request: %v", err)
	}

	// Trừ số dư
	user.Balance -= amount
	if err := ws.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user balance: %v", err)
	}

	return withdrawRequest, nil
}
