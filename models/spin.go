package models

import "gorm.io/gorm"

// SpinResult represents the result of a roulette spin.
type SpinResult struct {
	gorm.Model        // Includes ID, CreatedAt, UpdatedAt, DeletedAt fields
	Number     int    `json:"number" gorm:"not null"`    // The number from 0 to 36
	Color      string `json:"color" gorm:"not null"`     // Red or Black
	Parity     string `json:"parity" gorm:"not null"`    // Even or Odd
	Group      string `json:"group" gorm:"not null"`     // 1st12, 2nd12, or 3rd12
	Timestamp  int64  `json:"timestamp" gorm:"not null"` // Timestamp of the spin
}
