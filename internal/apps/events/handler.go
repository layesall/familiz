package events

import (
	"encoding/json"
	"net/http"
	"strconv"

	"familiz/internal/database"
	"familiz/internal/services"
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
	if req.EventDate == "" {
		http.Error(w, "event_date est obligatoire (format YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// GESTION DU MONTANT : Auto-calcul ou valeur saisie
	var finalAmount float64
	finalAmount, err := services.CalculateEventAmount(req.Type, req.AmountReceived)
	if err != nil {
		http.Error(w, "Erreur de calcul: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Vérifier que le membre existe
	var exists int
	err = database.DB.QueryRow("SELECT id FROM members WHERE id = ?", req.MemberID).Scan(&exists)
	if err != nil {
		http.Error(w, "Membre introuvable", http.StatusNotFound)
		return
	}

	// Insertion : on utilise finalAmount (CORRECTION)
	eventID, err := CreateEventRepo(req.MemberID, req.Type, finalAmount, req.EventDate)
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
		"amount_received": finalAmount, // CORRECTION : on renvoie le montant final
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

	memberIDStr := r.URL.Query().Get("member_id")
	if memberIDStr != "" {
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

	evts, err := GetAllEvents()
	if err != nil {
		http.Error(w, "Erreur récupération événements: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(evts)
}

// --- UPDATE Handler ---
func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

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

	existing, err := GetEventByID(id)
	if err != nil {
		http.Error(w, "Erreur base de données", http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "Événement introuvable", http.StatusNotFound)
		return
	}

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

// --- DELETE Handler ---
func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

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
