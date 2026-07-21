package services

import (
	"errors"
	"familiz/internal/apps/members"
	"familiz/internal/apps/settings"
)

// CalculateTransactionAmount retourne le montant final pour une transaction
func CalculateTransactionAmount(memberID int, customAmount float64) (float64, error) {
	if customAmount > 0 {
		return customAmount, nil
	}

	member, err := members.GetMemberByID(memberID)
	if err != nil {
		return 0, err
	}
	if member == nil {
		return 0, errors.New("membre introuvable")
	}

	contribSettings, err := settings.GetContributionSettings()
	if err != nil {
		return 0, err
	}

	switch member.MaritalStatus {
	case "single":
		return contribSettings.AmountSingle, nil
	case "married":
		return contribSettings.AmountMarried, nil
	case "minor":
		return contribSettings.AmountMinor, nil
	default:
		return 0, errors.New("statut marital inconnu")
	}
}

// CalculateEventAmount retourne le montant final pour un événement
func CalculateEventAmount(eventType string, customAmount float64) (float64, error) {
	if customAmount > 0 {
		return customAmount, nil
	}

	evtSetting, err := settings.GetEventSettingByType(eventType)
	if err != nil {
		return 0, err
	}
	if evtSetting == nil {
		return 0, errors.New("type d'événement non configuré")
	}
	return evtSetting.DefaultAmount, nil
}
