package transactions

import (
	"encoding/json"
	"net/http"
	"strconv"

	"familiz/internal/database"
	"familiz/internal/utils"
)

// CreateTransaction gère POST /transactions
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	// 1. Vérifier que l'utilisateur connecté est bien admin (rôle)
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	// 2. Lire le corps de la requête (JSON)
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON invalide: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 3. VALIDATIONS (obligatoire pour éviter les erreurs en base)
	if req.MemberID <= 0 {
		http.Error(w, "member_id est obligatoire et doit être > 0", http.StatusBadRequest)
		return
	}
	if req.Month < 1 || req.Month > 12 {
		http.Error(w, "mois invalide (doit être entre 1 et 12)", http.StatusBadRequest)
		return
	}
	if req.Year < 2000 || req.Year > 2100 {
		http.Error(w, "année invalide", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		http.Error(w, "le montant doit être supérieur à 0", http.StatusBadRequest)
		return
	}

	// 4. Vérifier que le membre existe vraiment (sécurité)
	var exists int
	err := database.DB.QueryRow("SELECT id FROM members WHERE id = ?", req.MemberID).Scan(&exists)
	if err != nil {
		http.Error(w, "Membre introuvable avec cet ID", http.StatusNotFound)
		return
	}

	// 5. Insérer la transaction en base
	result, err := database.DB.Exec(`
		INSERT INTO transactions (member_id, month, year, amount, note, paid_at, created_at)
		VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`, req.MemberID, req.Month, req.Year, req.Amount, req.Note)
	if err != nil {
		http.Error(w, "Erreur base de données lors de l'insertion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 6. Récupérer l'ID généré pour le renvoyer à l'admin
	transactionID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération de l'ID", http.StatusInternalServerError)
		return
	}

	// 7. Réponse de succès (format JSON)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "Paiement enregistré avec succès",
		"transaction_id": transactionID,
		"member_id":      req.MemberID,
		"month":          req.Month,
		"year":           req.Year,
		"amount":         req.Amount,
	})
}

// GetMemberTransactions gère GET /transactions?member_id=1
func GetMemberTransactions(w http.ResponseWriter, r *http.Request) {
	// 1. Lire le paramètre dans l'URL
	memberIDStr := r.URL.Query().Get("member_id")
	if memberIDStr == "" {
		http.Error(w, "Paramètre 'member_id' requis dans l'URL (ex: ?member_id=1)", http.StatusBadRequest)
		return
	}

	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		http.Error(w, "member_id doit être un nombre valide", http.StatusBadRequest)
		return
	}

	// 2. Optionnel : on peut aussi vérifier le rôle ici si on veut
	// (mais la route est déjà protégée par le middleware)

	// 3. Requête SQL pour récupérer toutes les transactions de ce membre
	rows, err := database.DB.Query(`
		SELECT id, member_id, month, year, amount, paid_at, note, created_at
		FROM transactions
		WHERE member_id = ?
		ORDER BY year DESC, month DESC
	`, memberID)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture des transactions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// 4. Transformer les lignes SQL en slice de structs
	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(
			&t.ID,
			&t.MemberID,
			&t.Month,
			&t.Year,
			&t.Amount,
			&t.PaidAt,
			&t.Note,
			&t.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Erreur lors du parsing des données", http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, t)
	}

	// 5. Renvoyer la liste en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
