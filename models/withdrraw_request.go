package models

import "time"

// WithdrawRequest represents a withdrawal transaction.
type WithdrawRequest struct {
	ID          uint       `gorm:"primaryKey"`
	PlayerID    string     `json:"player_id"`    // Wallet address
	Amount      float64    `json:"amount"`       // Số tiền rút
	Status      string     `json:"status"`       // pending, completed, failed
	RequestedAt time.Time  `json:"requested_at"` // Thời gian yêu cầu
	ProcessedAt *time.Time `json:"processed_at"` // Thời gian xử lý (nếu có)
}
