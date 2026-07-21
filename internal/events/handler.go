package events

import (
	"encoding/json"
	"net/http"
	"strconv"

	"familiz/internal/database"
	"familiz/internal/utils"
)

// --- CREATE Handler ---
func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON invalide: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validations
	if req.MemberID <= 0 {
		http.Error(w, "member_id est obligatoire", http.StatusBadRequest)
		return
	}
	if req.Type != "wedding" && req.Type != "baptism" {
		http.Error(w, "type doit être 'wedding' ou 'baptism'", http.StatusBadRequest)
		return
	}
	if req.AmountReceived < 0 {
		http.Error(w, "amount_received ne peut pas être négatif", http.StatusBadRequest)
		return
	}
	if req.EventDate == "" {
		http.Error(w, "event_date est obligatoire (format YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// Vérifier que le membre existe
	var exists int
	err := database.DB.QueryRow("SELECT id FROM members WHERE id = ?", req.MemberID).Scan(&exists)
	if err != nil {
		http.Error(w, "Membre introuvable", http.StatusNotFound)
		return
	}

	// Insertion
	eventID, err := CreateEventRepo(req.MemberID, req.Type, req.AmountReceived, req.EventDate)
	if err != nil {
		http.Error(w, "Erreur insertion événement: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":         "Événement enregistré avec succès",
		"event_id":        eventID,
		"member_id":       req.MemberID,
		"type":            req.Type,
		"amount_received": req.AmountReceived,
		"event_date":      req.EventDate,
	})
}

// --- LIST (avec ou sans filtre) Handler ---
func ListEventsHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	// Vérifier si le paramètre member_id est présent
	memberIDStr := r.URL.Query().Get("member_id")
	if memberIDStr != "" {
		// Mode FILTRE : on récupère par membre
		memberID, err := strconv.Atoi(memberIDStr)
		if err != nil {
			http.Error(w, "member_id invalide", http.StatusBadRequest)
			return
		}

		evts, err := GetEventsByMemberID(memberID)
		if err != nil {
			http.Error(w, "Erreur récupération événements: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(evts)
		return
	}

	// Mode GLOBAL : tous les événements
	evts, err := GetAllEvents()
	if err != nil {
		http.Error(w, "Erreur récupération événements: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(evts)
}

// --- UPDATE Handler (NOUVEAU) ---
func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	// Extraire l'ID de l'URL
	idStr := r.URL.Path[len("/events/"):]
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// Vérifier que l'événement existe
	existing, err := GetEventByID(id)
	if err != nil {
		http.Error(w, "Erreur base de données", http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "Événement introuvable", http.StatusNotFound)
		return
	}

	// Décoder la requête
	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validations
	if req.Type != "wedding" && req.Type != "baptism" {
		http.Error(w, "type doit être 'wedding' ou 'baptism'", http.StatusBadRequest)
		return
	}
	if req.AmountReceived < 0 {
		http.Error(w, "amount_received ne peut pas être négatif", http.StatusBadRequest)
		return
	}
	if req.EventDate == "" {
		http.Error(w, "event_date est obligatoire (format YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// Mise à jour
	err = UpdateEventRepo(id, req)
	if err != nil {
		if err.Error() == "événement introuvable" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Erreur mise à jour: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Événement mis à jour avec succès",
	})
}

// --- DELETE Handler (NOUVEAU) ---
func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	// Extraire l'ID
	idStr := r.URL.Path[len("/events/"):]
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// Supprimer
	err = DeleteEventRepo(id)
	if err != nil {
		if err.Error() == "événement introuvable" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Erreur suppression: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Événement supprimé avec succès",
	})
}
