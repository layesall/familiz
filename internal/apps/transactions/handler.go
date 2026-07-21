package transactions

import (
	"encoding/json"
	"net/http"
	"strconv"

	"familiz/internal/apps/settings"
	"familiz/internal/database"
	"familiz/internal/services"
	"familiz/internal/utils"
)

// --- CREATE Handler ---
func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON invalide: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validations
	if req.MemberID <= 0 {
		http.Error(w, "member_id est obligatoire", http.StatusBadRequest)
		return
	}
	if req.Month < 1 || req.Month > 12 {
		http.Error(w, "mois invalide (1-12)", http.StatusBadRequest)
		return
	}

	// FORCER l'année en cours (l'admin ne peut pas créer de transaction pour une année passée)
	currentYear, err := settings.GetCurrentYear()
	if err != nil {
		http.Error(w, "Erreur récupération année en cours: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Year = currentYear

	// Calcul automatique du montant
	finalAmount, err := services.CalculateTransactionAmount(req.MemberID, req.Amount)
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
	transactionID, err := CreateTransactionRepo(req.MemberID, req.Month, req.Year, finalAmount, req.Note)
	if err != nil {
		http.Error(w, "Erreur insertion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "Paiement enregistré avec succès",
		"transaction_id": transactionID,
		"member_id":      req.MemberID,
		"month":          req.Month,
		"year":           req.Year,
		"amount":         finalAmount,
	})
}

// --- LIST Handler ---
func ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	// Lire le paramètre ?archived=true (défaut : false)
	includeArchived := r.URL.Query().Get("archived") == "true"

	// Lire le paramètre ?member_id=...
	memberIDStr := r.URL.Query().Get("member_id")
	if memberIDStr != "" {
		memberID, err := strconv.Atoi(memberIDStr)
		if err != nil {
			http.Error(w, "member_id invalide", http.StatusBadRequest)
			return
		}

		txs, err := GetTransactionsByMemberID(memberID, includeArchived)
		if err != nil {
			http.Error(w, "Erreur récupération transactions: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(txs)
		return
	}

	// Mode GLOBAL
	txs, err := GetAllTransactions(includeArchived)
	if err != nil {
		http.Error(w, "Erreur récupération transactions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txs)
}

// --- UPDATE Handler ---
func UpdateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	idStr := r.URL.Path[len("/transactions/"):]
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var req UpdateTransactionRequest
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

	err = UpdateTransactionRepo(id, req)
	if err != nil {
		if err.Error() == "transaction introuvable ou déjà archivée" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if err.Error() == "impossible de modifier une transaction archivée" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, "Erreur mise à jour: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Transaction mise à jour avec succès",
	})
}

// --- DELETE Handler ---
func DeleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	idStr := r.URL.Path[len("/transactions/"):]
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	err = DeleteTransactionRepo(id)
	if err != nil {
		if err.Error() == "transaction introuvable ou déjà archivée" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if err.Error() == "impossible de supprimer une transaction archivée" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, "Erreur suppression: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Transaction supprimée avec succès",
	})
}
