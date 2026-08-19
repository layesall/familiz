package auth

import (
	"encoding/json"
	"net/http"

	"familiz/internal/apps/members"
	"familiz/internal/database"
	"familiz/internal/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Vérifier si l'email existe déjà
	existing, _ := FindUserByEmail(req.Email)
	if existing != nil {
		http.Error(w, "Cet email est déjà utilisé", http.StatusConflict)
		return
	}

	// Démarrer une transaction
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Erreur base de données", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Créer le membre via le package members
	memberID, err := members.CreateMemberTx(tx, req.FirstName, req.LastName, req.BirthDate, req.MaritalStatus)
	if err != nil {
		http.Error(w, "Erreur création membre: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Créer l'utilisateur (on passe la transaction tx)
	_, err = CreateUser(tx, req, memberID) // <-- MODIFICATION ICI
	if err != nil {
		http.Error(w, "Erreur création utilisateur: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Erreur validation transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Admin créé avec succès",
		"member_id": memberID,
		"email":     req.Email,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	user, err := FindUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, "Erreur génération token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Token: token,
		Email: user.Email,
		Role:  user.Role,
	})
}
