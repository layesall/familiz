package members

import (
	"encoding/json"
	"net/http"

	// "strconv"

	"familiz/internal/database"
	"familiz/internal/utils"
)

type CreateMemberRequest struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	BirthDate     string `json:"birth_date"`
	MaritalStatus string `json:"marital_status"`
}

// CreateMember est un handler protégé par JWT
func CreateMember(w http.ResponseWriter, r *http.Request) {
	// 1. Vérifier que l'utilisateur est admin (optionnel pour l'instant)
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	// 2. Décoder la requête
	var req CreateMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Validations simples
	if req.FirstName == "" || req.LastName == "" {
		http.Error(w, "Prénom et nom sont obligatoires", http.StatusBadRequest)
		return
	}

	// 4. Insérer en base
	result, err := database.DB.Exec(`
        INSERT INTO members (first_name, last_name, birth_date, marital_status, created_at, updated_at)
        VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
    `, req.FirstName, req.LastName, req.BirthDate, req.MaritalStatus)
	if err != nil {
		http.Error(w, "Erreur lors de la création du membre: "+err.Error(), http.StatusInternalServerError)
		return
	}

	memberID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération de l'ID", http.StatusInternalServerError)
		return
	}

	// 5. Réponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Membre créé avec succès",
		"member_id":  memberID,
		"first_name": req.FirstName,
		"last_name":  req.LastName,
	})
}

// ListMembers retourne tous les membres (pour tester)
func ListMembers(w http.ResponseWriter, r *http.Request) {

	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé: admin requis", http.StatusForbidden)
	}

	rows, err := database.DB.Query(`
        SELECT id, first_name, last_name, birth_date, marital_status 
        FROM members 
        ORDER BY id DESC
    `)
	if err != nil {
		http.Error(w, "Erreur base de données", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var m Member
		err := rows.Scan(&m.ID, &m.FirstName, &m.LastName, &m.BirthDate, &m.MaritalStatus)
		if err != nil {
			http.Error(w, "Erreur lecture données", http.StatusInternalServerError)
			return
		}
		members = append(members, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}
