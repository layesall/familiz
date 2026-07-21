package events

import "time"

type Event struct {
	ID             int       `json:"id"`
	MemberID       int       `json:"member_id"`
	Type           string    `json:"type"`
	AmountReceived float64   `json:"amount_received"`
	EventDate      string    `json:"event_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Requête pour CRÉER un événement
type CreateEventRequest struct {
	MemberID       int     `json:"member_id"`
	Type           string  `json:"type"`
	AmountReceived float64 `json:"amount_received"`
	EventDate      string  `json:"event_date"`
}

// Requête pour METTRE À JOUR un événement (NOUVEAU)
type UpdateEventRequest struct {
	Type           string  `json:"type"`
	AmountReceived float64 `json:"amount_received"`
	EventDate      string  `json:"event_date"`
}
