package models

// Wallet struct represents an Ethereum wallet with its address and balance.
type Wallet struct {
	Address string
	Balance float64
}

// ValidateBalance checks if the player has enough balance to place the bet.
func (w *Wallet) ValidateBalance(amount float64) bool {
	return w.Balance >= amount
}
