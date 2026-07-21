package settings

import (
	"encoding/json"
	"net/http"

	"familiz/internal/utils"
)

// --- GET /settings/contributions ---
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

// --- PUT /settings/contributions ---
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

// --- GET /settings/events ---
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

// --- PUT /settings/events/{type} ---
func UpdateEventSettingHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	// Extraire le type depuis l'URL
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
