package transactions

import "time"

type Transaction struct {
	ID         int       `json:"id"`
	MemberID   int       `json:"member_id"`
	Month      int       `json:"month"`
	Year       int       `json:"year"`
	Amount     float64   `json:"amount"`
	Note       string    `json:"note"`
	PaidAt     time.Time `json:"paid_at"`
	IsArchived bool      `json:"is_archived"`
	CreatedAt  time.Time `json:"created_at"`
}

// Requête pour CRÉER un paiement
type CreateTransactionRequest struct {
	MemberID int     `json:"member_id"`
	Month    int     `json:"month"`
	Year     int     `json:"year"`
	Amount   float64 `json:"amount"`
	Note     string  `json:"note"`
}

// Requête pour METTRE À JOUR un paiement (NOUVEAU)
type UpdateTransactionRequest struct {
	Month  int     `json:"month"`
	Year   int     `json:"year"`
	Amount float64 `json:"amount"`
	Note   string  `json:"note"`
}
