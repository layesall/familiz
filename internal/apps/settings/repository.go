package settings

import (
	"database/sql"
	"errors"
	"familiz/internal/database"
)

func GetContributionSettings() (*ContributionSettings, error) {
	var s ContributionSettings
	err := database.DB.QueryRow(`
        SELECT amount_single, amount_married, amount_minor, current_year, updated_at
        FROM contribution_settings
        WHERE id = 1
    `).Scan(&s.AmountSingle, &s.AmountMarried, &s.AmountMinor, &s.CurrentYear, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func UpdateContributionSettings(req UpdateContributionRequest) error {
	if req.AmountSingle < 0 || req.AmountMarried < 0 || req.AmountMinor < 0 {
		return errors.New("les montants ne peuvent pas être négatifs")
	}

	result, err := database.DB.Exec(`
        UPDATE contribution_settings
        SET amount_single = ?, amount_married = ?, amount_minor = ?, updated_at = datetime('now')
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

func UpdateEventSetting(eventType string, req UpdateEventSettingRequest) error {
	if req.DefaultAmount < 0 {
		return errors.New("le montant ne peut pas être négatif")
	}

	result, err := database.DB.Exec(`
        UPDATE event_settings
        SET default_amount = ?, updated_at = datetime('now')
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

// --- ARCHIVAGE ---

func GetCurrentYear() (int, error) {
	var year int
	err := database.DB.QueryRow("SELECT current_year FROM contribution_settings WHERE id = 1").Scan(&year)
	if err != nil {
		return 0, err
	}
	return year, nil
}

func ArchiveYear() error {
	currentYear, err := GetCurrentYear()
	if err != nil {
		return err
	}

	// Archiver les transactions
	_, err = database.DB.Exec(`
        UPDATE transactions SET is_archived = 1 WHERE year = ? AND is_archived = 0
    `, currentYear)
	if err != nil {
		return err
	}

	// Archiver les événements (l'année est dans la date)
	_, err = database.DB.Exec(`
        UPDATE events SET is_archived = 1 WHERE strftime('%Y', event_date) = ? AND is_archived = 0
    `, currentYear)
	if err != nil {
		return err
	}

	// Incrémenter l'année
	_, err = database.DB.Exec(`
        UPDATE contribution_settings SET current_year = current_year + 1, updated_at = datetime('now') WHERE id = 1
    `)
	if err != nil {
		return err
	}

	return nil
}
