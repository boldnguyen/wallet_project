package models

// BetType defines the type of the bet (number, group, color, even/odd).
type BetType string

const (
	NumberBet  BetType = "number"
	GroupBet   BetType = "group"
	ColorBet   BetType = "color"
	EvenOddBet BetType = "even_odd"
)

type Bet struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	PlayerID  string  `json:"player_id"`
	BetType   BetType `json:"bet_type"`
	Amount    float64 `json:"amount"`
	Selection string  `json:"selection"`
	Status    string  `json:"status"`  // "placed", "won", "lost"
	SpinID    int     `json:"spin_id"` // ID of the spin this bet is associated with
	Timestamp int64   `json:"timestamp"`
}
