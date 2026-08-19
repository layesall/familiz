package contributions

import (
	"encoding/json"
	"net/http"
	"strconv"

	"familiz/internal/database"
	"familiz/internal/utils"
)

// CREATE
func CreateContributionHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	var req CreateContributionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON invalide: "+err.Error(), http.StatusBadRequest)
		return
	}

	// ✅ Validations
	if req.MemberID <= 0 {
		http.Error(w, "member_id est obligatoire", http.StatusBadRequest)
		return
	}
	if req.Month < 1 || req.Month > 12 {
		http.Error(w, "mois invalide (1-12)", http.StatusBadRequest)
		return
	}
	if req.Year < 2000 || req.Year > 2100 {
		http.Error(w, "année invalide (2000-2100)", http.StatusBadRequest)
		return
	}

	// Calcul auto du montant (service)
	finalAmount, err := CalculateContributionAmount(req.MemberID, req.Amount)
	if err != nil {
		http.Error(w, "Erreur de calcul: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Vérifier l'existence du membre
	var exists int
	err = database.DB.QueryRow("SELECT id FROM members WHERE id = ?", req.MemberID).Scan(&exists)
	if err != nil {
		http.Error(w, "Membre introuvable", http.StatusNotFound)
		return
	}

	// Insertion
	contributionID, err := CreateContributionRepo(req.MemberID, req.Month, req.Year, finalAmount, req.Note)
	if err != nil {
		http.Error(w, "Erreur insertion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":         "Contribution enregistrée avec succès",
		"contribution_id": contributionID,
		"member_id":       req.MemberID,
		"month":           req.Month,
		"year":            req.Year,
		"amount":          finalAmount,
	})
}

// --- LIST Handler ---
func ListContributionsHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	includeArchived := r.URL.Query().Get("archived") == "true"
	memberIDStr := r.URL.Query().Get("member_id")

	if memberIDStr != "" {
		memberID, err := strconv.Atoi(memberIDStr)
		if err != nil {
			http.Error(w, "member_id invalide", http.StatusBadRequest)
			return
		}

		contributions, err := GetContributionsByMemberID(memberID, includeArchived)
		if err != nil {
			http.Error(w, "Erreur récupération contributions: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(contributions)
		return
	}

	// Mode GLOBAL
	contributions, err := GetAllContributions(includeArchived)
	if err != nil {
		http.Error(w, "Erreur récupération contributions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contributions)
}

// --- UPDATE Handler ---
func UpdateContributionHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	idStr := r.URL.Path[len("/contributions/"):]
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var req UpdateContributionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Month < 1 || req.Month > 12 {
		http.Error(w, "mois invalide (1-12)", http.StatusBadRequest)
		return
	}
	if req.Year < 2000 {
		http.Error(w, "année invalide", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		http.Error(w, "montant doit être > 0", http.StatusBadRequest)
		return
	}

	err = UpdateContributionRepo(id, req)
	if err != nil {
		if err.Error() == "contribution introuvable ou déjà archivée" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if err.Error() == "impossible de modifier une contribution archivée" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, "Erreur mise à jour: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Contribution mise à jour avec succès",
	})
}

// --- DELETE Handler ---
func DeleteContributionHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	idStr := r.URL.Path[len("/contributions/"):]
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	err = DeleteContributionRepo(id)
	if err != nil {
		if err.Error() == "contribution introuvable ou déjà archivée" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if err.Error() == "impossible de supprimer une contribution archivée" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, "Erreur suppression: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Contribution supprimée avec succès",
	})
}
