package contributions

import (
	"errors"
	"familiz/internal/database"
)

// CalculateContributionAmount calcule le montant automatique si demandé = 0
func CalculateContributionAmount(memberID int, requestedAmount float64) (float64, error) {
	// Si l'admin a mis un montant > 0, on le garde
	if requestedAmount > 0 {
		return requestedAmount, nil
	}

	// Récupérer le statut marital du membre
	var maritalStatus string
	err := database.DB.QueryRow("SELECT marital_status FROM members WHERE id = ?", memberID).Scan(&maritalStatus)
	if err != nil {
		return 0, errors.New("membre introuvable")
	}

	// Récupérer le montant correspondant dans contribution_settings
	var amount float64
	query := `
		SELECT CASE ? 
			WHEN 'single' THEN amount_single
			WHEN 'married' THEN amount_married
			WHEN 'minor' THEN amount_minor
		END FROM contribution_settings WHERE id = 1
	`
	err = database.DB.QueryRow(query, maritalStatus).Scan(&amount)
	if err != nil {
		return 0, errors.New("paramètres de cotisation non trouvés")
	}

	return amount, nil
}
