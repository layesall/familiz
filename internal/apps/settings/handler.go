package settings

import (
	"encoding/json"
	"net/http"
	"strconv"

	"familiz/internal/database"
	"familiz/internal/utils"
)

func GetContributionSettingsHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	settings, err := GetContributionSettings()
	if err != nil {
		http.Error(w, "Erreur récupération paramètres", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func UpdateContributionSettingsHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	var req UpdateContributionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := UpdateContributionSettings(req)
	if err != nil {
		http.Error(w, "Erreur mise à jour: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Paramètres de cotisation mis à jour avec succès",
	})
}

func GetEventSettingsHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	settings, err := GetAllEventSettings()
	if err != nil {
		http.Error(w, "Erreur récupération paramètres", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func UpdateEventSettingHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	eventType := r.URL.Path[len("/settings/events/"):]
	if eventType == "" {
		http.Error(w, "Type d'événement manquant (wedding ou baptism)", http.StatusBadRequest)
		return
	}
	if eventType != "wedding" && eventType != "baptism" {
		http.Error(w, "Type invalide. Utilisez 'wedding' ou 'baptism'", http.StatusBadRequest)
		return
	}

	var req UpdateEventSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := UpdateEventSetting(eventType, req)
	if err != nil {
		http.Error(w, "Erreur mise à jour: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Paramètre événement mis à jour avec succès",
	})
}

// --- ARCHIVAGE HANDLER ---

func ArchiveYearHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	currentYear, err := GetCurrentYear()
	if err != nil {
		http.Error(w, "Erreur récupération année: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Vérifier s'il y a des données pour cette année
	var count int
	err = database.DB.QueryRow(`
        SELECT COUNT(*) FROM transactions WHERE year = ? AND is_archived = 0
    `, currentYear).Scan(&count)
	if err != nil {
		http.Error(w, "Erreur vérification transactions", http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "Aucune transaction à archiver pour l'année "+strconv.Itoa(currentYear), http.StatusBadRequest)
		return
	}

	err = ArchiveYear()
	if err != nil {
		http.Error(w, "Erreur lors de l'archivage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Année archivée avec succès",
		"archived_year": currentYear,
		"new_year":      currentYear + 1,
	})
}
