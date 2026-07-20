package events

import "time"

// Event représente une ligne dans la table events
type Event struct {
	ID             int       `json:"id"`
	MemberID       int       `json:"member_id"`
	Type           string    `json:"type"` // "wedding" ou "baptism"
	AmountReceived float64   `json:"amount_received"`
	EventDate      string    `json:"event_date"` // format YYYY-MM-DD
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CreateEventRequest représente le JSON que l'admin envoie
type CreateEventRequest struct {
	MemberID       int     `json:"member_id"`
	Type           string  `json:"type"` // "wedding" ou "baptism"
	AmountReceived float64 `json:"amount_received"`
	EventDate      string  `json:"event_date"` // "2025-06-15"
}
