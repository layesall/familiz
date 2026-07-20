package events

import (
	"encoding/json"
	"net/http"
	"strconv"

	"familiz/internal/database"
	"familiz/internal/utils"
)

// CreateEvent gère POST /events
func CreateEvent(w http.ResponseWriter, r *http.Request) {
	// 1. Vérification admin
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	// 2. Décodage JSON
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON invalide: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 3. VALIDATIONS
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

	// 4. Vérifier que le membre existe
	var exists int
	err := database.DB.QueryRow("SELECT id FROM members WHERE id = ?", req.MemberID).Scan(&exists)
	if err != nil {
		http.Error(w, "Membre introuvable", http.StatusNotFound)
		return
	}

	// 5. Insertion en base
	result, err := database.DB.Exec(`
        INSERT INTO events (member_id, type, amount_received, event_date, created_at, updated_at)
        VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
    `, req.MemberID, req.Type, req.AmountReceived, req.EventDate)
	if err != nil {
		http.Error(w, "Erreur insertion événement: "+err.Error(), http.StatusInternalServerError)
		return
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Erreur récupération ID", http.StatusInternalServerError)
		return
	}

	// 6. Réponse
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

// GetMemberEvents gère GET /events?member_id=1
func GetMemberEvents(w http.ResponseWriter, r *http.Request) {
	memberIDStr := r.URL.Query().Get("member_id")
	if memberIDStr == "" {
		http.Error(w, "Paramètre 'member_id' requis", http.StatusBadRequest)
		return
	}

	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		http.Error(w, "member_id invalide", http.StatusBadRequest)
		return
	}

	rows, err := database.DB.Query(`
        SELECT id, member_id, type, amount_received, event_date, created_at, updated_at
        FROM events
        WHERE member_id = ?
        ORDER BY event_date DESC
    `, memberID)
	if err != nil {
		http.Error(w, "Erreur base de données", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(
			&e.ID,
			&e.MemberID,
			&e.Type,
			&e.AmountReceived,
			&e.EventDate,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "Erreur lecture données", http.StatusInternalServerError)
			return
		}
		events = append(events, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
