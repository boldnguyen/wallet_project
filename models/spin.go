package models

// SpinResult stores the result of a roulette spin.
type SpinResult struct {
	ID        uint   `gorm:"primaryKey"`
	Number    int    `json:"number"`
	Color     string `json:"color"`
	Parity    string `json:"parity"`
	Group     string `json:"group"`
	Timestamp int64  `json:"timestamp"`
}
