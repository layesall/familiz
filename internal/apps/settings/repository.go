package settings

import (
	"database/sql"
	"errors"
	"familiz/internal/database"
)

// --- READ Contributions ---
func GetContributionSettings() (*ContributionSettings, error) {
	var s ContributionSettings
	err := database.DB.QueryRow(`
        SELECT amount_single, amount_married, amount_minor, updated_at
        FROM contribution_settings
        WHERE id = 1
    `).Scan(&s.AmountSingle, &s.AmountMarried, &s.AmountMinor, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// --- UPDATE Contributions ---
func UpdateContributionSettings(req UpdateContributionRequest) error {
	// Validation basique
	if req.AmountSingle < 0 || req.AmountMarried < 0 || req.AmountMinor < 0 {
		return errors.New("les montants ne peuvent pas être négatifs")
	}

	result, err := database.DB.Exec(`
        UPDATE contribution_settings
        SET amount_single = ?,
            amount_married = ?,
            amount_minor = ?,
            updated_at = datetime('now')
        WHERE id = 1
    `, req.AmountSingle, req.AmountMarried, req.AmountMinor)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("paramètres introuvables")
	}
	return nil
}

// --- READ Event Settings (pour un type spécifique) ---
func GetEventSettingByType(eventType string) (*EventSettings, error) {
	var s EventSettings
	err := database.DB.QueryRow(`
        SELECT event_type, default_amount, updated_at
        FROM event_settings
        WHERE event_type = ?
    `, eventType).Scan(&s.EventType, &s.DefaultAmount, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// --- READ All Event Settings ---
func GetAllEventSettings() ([]EventSettings, error) {
	rows, err := database.DB.Query(`
        SELECT event_type, default_amount, updated_at
        FROM event_settings
        ORDER BY event_type
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []EventSettings
	for rows.Next() {
		var s EventSettings
		err := rows.Scan(&s.EventType, &s.DefaultAmount, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, nil
}

// --- UPDATE Event Setting ---
func UpdateEventSetting(eventType string, req UpdateEventSettingRequest) error {
	if req.DefaultAmount < 0 {
		return errors.New("le montant ne peut pas être négatif")
	}

	result, err := database.DB.Exec(`
        UPDATE event_settings
        SET default_amount = ?,
            updated_at = datetime('now')
        WHERE event_type = ?
    `, req.DefaultAmount, eventType)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("type d'événement introuvable")
	}
	return nil
}
