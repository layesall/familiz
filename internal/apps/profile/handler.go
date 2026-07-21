package profile

import (
	"encoding/json"
	"net/http"
	"strconv"

	"familiz/internal/apps/events"
	"familiz/internal/apps/members"
	"familiz/internal/apps/transactions"
	"familiz/internal/utils"
)

type MemberProfileResponse struct {
	Member       members.Member             `json:"member"`
	Transactions []transactions.Transaction `json:"transactions"`
	Events       []events.Event             `json:"events"`
}

// GetMemberProfileHandler gère GET /profile/{id}
func GetMemberProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Vérification admin (optionnelle mais conseillée)
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	// Extraire l'ID de l'URL
	idStr := r.URL.Path[len("/profile/"):]
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// Récupérer le membre
	member, err := members.GetMemberByID(id)
	if err != nil {
		http.Error(w, "Erreur base de données", http.StatusInternalServerError)
		return
	}
	if member == nil {
		http.Error(w, "Membre introuvable", http.StatusNotFound)
		return
	}

	// Récupérer ses transactions
	txs, err := transactions.GetTransactionsByMemberID(id)
	if err != nil {
		http.Error(w, "Erreur récupération transactions", http.StatusInternalServerError)
		return
	}

	// Récupérer ses événements
	evts, err := events.GetEventsByMemberID(id)
	if err != nil {
		http.Error(w, "Erreur récupération événements", http.StatusInternalServerError)
		return
	}

	// Construire la réponse
	profile := MemberProfileResponse{
		Member:       *member,
		Transactions: txs,
		Events:       evts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
