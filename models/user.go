package models

// User represents a player in the system
type User struct {
	ID       uint   `gorm:"primaryKey"`
	PlayerID string `gorm:"unique;not null"`
	Balance  float64
}
