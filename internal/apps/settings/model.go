package settings

// ContributionSettings représente la configuration des montants de cotisation
type ContributionSettings struct {
	AmountSingle  float64 `json:"amount_single"`
	AmountMarried float64 `json:"amount_married"`
	AmountMinor   float64 `json:"amount_minor"`
	UpdatedAt     string  `json:"updated_at"`
}

// EventSettings représente la configuration des montants par défaut pour les événements
type EventSettings struct {
	EventType     string  `json:"event_type"`
	DefaultAmount float64 `json:"default_amount"`
	UpdatedAt     string  `json:"updated_at"`
}

// Requête pour mettre à jour les cotisations
type UpdateContributionRequest struct {
	AmountSingle  float64 `json:"amount_single"`
	AmountMarried float64 `json:"amount_married"`
	AmountMinor   float64 `json:"amount_minor"`
}

// Requête pour mettre à jour un événement
type UpdateEventSettingRequest struct {
	DefaultAmount float64 `json:"default_amount"`
}
