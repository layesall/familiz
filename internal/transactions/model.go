package transactions

import "time"

// Transaction représente une ligne dans la table transactions
type Transaction struct {
	ID        int       `json:"id"`
	MemberID  int       `json:"member_id"`
	Month     int       `json:"month"`   // 1 à 12
	Year      int       `json:"year"`    // ex: 2025
	Amount    float64   `json:"amount"`  // montant payé
	Note      string    `json:"note"`    // optionnel (ex: "retard")
	PaidAt    time.Time `json:"paid_at"` // date d'enregistrement
	CreatedAt time.Time `json:"created_at"`
}

// CreateTransactionRequest représente le JSON que l'admin envoie pour payer
type CreateTransactionRequest struct {
	MemberID int     `json:"member_id"`
	Month    int     `json:"month"`
	Year     int     `json:"year"`
	Amount   float64 `json:"amount"`
	Note     string  `json:"note"`
}
