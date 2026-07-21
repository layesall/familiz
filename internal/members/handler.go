package members

import (
	"encoding/json"
	"net/http"
	"strconv"

	"familiz/internal/database"
	"familiz/internal/events"
	"familiz/internal/transactions"
	"familiz/internal/utils"
)

// --- 1. CREATE (existante, avec vérification admin) ---
func CreateMemberHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	var req CreateMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.FirstName == "" || req.LastName == "" {
		http.Error(w, "Prénom et nom sont obligatoires", http.StatusBadRequest)
		return
	}

	// Utiliser le repository (mais on n'a pas de fonction Create standard, on va la faire directement ici pour l'instant)
	result, err := database.DB.Exec(`
        INSERT INTO members (first_name, last_name, birth_date, marital_status, created_at, updated_at)
        VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
    `, req.FirstName, req.LastName, req.BirthDate, req.MaritalStatus)
	if err != nil {
		http.Error(w, "Erreur création membre: "+err.Error(), http.StatusInternalServerError)
		return
	}

	memberID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Erreur récupération ID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Membre créé avec succès",
		"member_id":  memberID,
		"first_name": req.FirstName,
		"last_name":  req.LastName,
	})
}

// --- 2. LIST ALL (existante, on la garde) ---
func ListMembersHandler(w http.ResponseWriter, r *http.Request) {
	members, err := GetAllMembers()
	if err != nil {
		http.Error(w, "Erreur base de données", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

// --- 3. GET BY ID (NOUVEAU : Profil complet) ---
type MemberProfileResponse struct {
	Member       Member                     `json:"member"`
	Transactions []transactions.Transaction `json:"transactions"`
	Events       []events.Event             `json:"events"`
}

func GetMemberHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extraire l'ID de l'URL (ex: /members/5)
	idStr := r.URL.Path[len("/members/"):]
	if idStr == "" {
		http.Error(w, "ID manquant dans l'URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// 2. Récupérer le membre
	member, err := GetMemberByID(id)
	if err != nil {
		http.Error(w, "Erreur base de données", http.StatusInternalServerError)
		return
	}
	if member == nil {
		http.Error(w, "Membre introuvable", http.StatusNotFound)
		return
	}

	// 3. Récupérer ses transactions
	txs, err := transactions.GetTransactionsByMemberID(id)
	if err != nil {
		http.Error(w, "Erreur récupération transactions", http.StatusInternalServerError)
		return
	}

	// 4. Récupérer ses événements
	evts, err := events.GetEventsByMemberID(id)
	if err != nil {
		http.Error(w, "Erreur récupération événements", http.StatusInternalServerError)
		return
	}

	// 5. Construire la réponse
	profile := MemberProfileResponse{
		Member:       *member,
		Transactions: txs,
		Events:       evts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// --- 4. UPDATE (NOUVEAU) ---
func UpdateMemberHandler(w http.ResponseWriter, r *http.Request) {
	// Vérification admin
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	// Extraire ID
	idStr := r.URL.Path[len("/members/"):]
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// Décoder la requête
	var req UpdateMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validation simple
	if req.FirstName == "" || req.LastName == "" {
		http.Error(w, "Prénom et nom sont obligatoires", http.StatusBadRequest)
		return
	}

	// Mettre à jour
	err = UpdateMember(id, req)
	if err != nil {
		if err.Error() == "membre introuvable" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Erreur mise à jour: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Membre mis à jour avec succès",
	})
}

// --- 5. DELETE (NOUVEAU) ---
func DeleteMemberHandler(w http.ResponseWriter, r *http.Request) {
	// Vérification admin
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	// Extraire ID
	idStr := r.URL.Path[len("/members/"):]
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// Supprimer (les transactions et events sont supprimés en cascade)
	err = DeleteMember(id)
	if err != nil {
		if err.Error() == "membre introuvable" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Erreur suppression: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Membre supprimé avec succès (ainsi que ses transactions et événements)",
	})
}
