package events

import (
	"errors"
	"familiz/internal/database"
)

// CalculateEventAmount calcule le montant par défaut selon le type d'événement
func CalculateEventAmount(eventType string, requestedAmount float64) (float64, error) {
	if requestedAmount > 0 {
		return requestedAmount, nil
	}

	var defaultAmount float64
	err := database.DB.QueryRow(
		"SELECT default_amount FROM event_settings WHERE event_type = ?",
		eventType,
	).Scan(&defaultAmount)
	if err != nil {
		return 0, errors.New("paramètre événement non trouvé")
	}
	return defaultAmount, nil
}
