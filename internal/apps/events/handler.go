package events

import (
	"encoding/json"
	"net/http"
	"strconv"

	"familiz/internal/database"
	"familiz/internal/services"
	"familiz/internal/utils"
)

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

	finalAmount, err := services.CalculateEventAmount(req.Type, req.AmountReceived)
	if err != nil {
		http.Error(w, "Erreur de calcul: "+err.Error(), http.StatusBadRequest)
		return
	}

	var exists int
	err = database.DB.QueryRow("SELECT id FROM members WHERE id = ?", req.MemberID).Scan(&exists)
	if err != nil {
		http.Error(w, "Membre introuvable", http.StatusNotFound)
		return
	}

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
		"amount_received": finalAmount,
		"event_date":      req.EventDate,
	})
}

func ListEventsHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	// Paramètre ?archived=true
	includeArchived := r.URL.Query().Get("archived") == "true"

	memberIDStr := r.URL.Query().Get("member_id")
	if memberIDStr != "" {
		memberID, err := strconv.Atoi(memberIDStr)
		if err != nil {
			http.Error(w, "member_id invalide", http.StatusBadRequest)
			return
		}

		evts, err := GetEventsByMemberID(memberID, includeArchived)
		if err != nil {
			http.Error(w, "Erreur récupération événements: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(evts)
		return
	}

	evts, err := GetAllEvents(includeArchived)
	if err != nil {
		http.Error(w, "Erreur récupération événements: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(evts)
}

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

	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide: "+err.Error(), http.StatusBadRequest)
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

	err = UpdateEventRepo(id, req)
	if err != nil {
		if err.Error() == "événement introuvable ou déjà archivé" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if err.Error() == "impossible de modifier un événement archivé" {
			http.Error(w, err.Error(), http.StatusForbidden)
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
		if err.Error() == "événement introuvable ou déjà archivé" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if err.Error() == "impossible de supprimer un événement archivé" {
			http.Error(w, err.Error(), http.StatusForbidden)
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
