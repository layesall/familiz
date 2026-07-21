package reports

import (
	"net/http"
	"strconv"

	"familiz/internal/apps/events"
	"familiz/internal/apps/members"
	"familiz/internal/apps/transactions"
	"familiz/internal/utils"
)

// MemberPDFHandler génère le PDF pour un membre
func MemberPDFHandler(w http.ResponseWriter, r *http.Request) {
	// Vérification admin
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	// Récupérer l'ID
	idStr := r.URL.Path[len("/members/"):]
	idStr = idStr[:len(idStr)-4] // enlever le "/pdf"
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// Récupérer les données
	member, err := members.GetMemberByID(id)
	if err != nil {
		http.Error(w, "Erreur récupération membre", http.StatusInternalServerError)
		return
	}
	if member == nil {
		http.Error(w, "Membre introuvable", http.StatusNotFound)
		return
	}

	includeArchived := r.URL.Query().Get("archived") == "true"
	txs, err := transactions.GetTransactionsByMemberID(id, includeArchived)
	if err != nil {
		http.Error(w, "Erreur récupération transactions", http.StatusInternalServerError)
		return
	}
	evts, err := events.GetEventsByMemberID(id, includeArchived)
	if err != nil {
		http.Error(w, "Erreur récupération événements", http.StatusInternalServerError)
		return
	}

	// Générer le PDF
	pdf, err := generateMemberProfilePDF(member, txs, evts)
	if err != nil {
		http.Error(w, "Erreur génération PDF", http.StatusInternalServerError)
		return
	}

	// Envoyer le PDF
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=membre_"+strconv.Itoa(id)+".pdf")
	err = pdf.Output(w)
	if err != nil {
		http.Error(w, "Erreur envoi PDF", http.StatusInternalServerError)
		return
	}
}

// GlobalReportPDFHandler génère le rapport global pour une année
func GlobalReportPDFHandler(w http.ResponseWriter, r *http.Request) {
	// Vérification admin
	role, ok := r.Context().Value(utils.UserRoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Accès refusé : admin requis", http.StatusForbidden)
		return
	}

	// Récupérer l'année
	yearStr := r.URL.Query().Get("year")
	if yearStr == "" {
		http.Error(w, "Paramètre year requis (ex: ?year=2026)", http.StatusBadRequest)
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 {
		http.Error(w, "Année invalide", http.StatusBadRequest)
		return
	}

	// Récupérer les données
	summaries, err := GetGlobalReportData(year)
	if err != nil {
		http.Error(w, "Erreur récupération données", http.StatusInternalServerError)
		return
	}

	// Générer le PDF
	pdf, err := generateGlobalReportPDF(summaries, year)
	if err != nil {
		http.Error(w, "Erreur génération PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=rapport_"+strconv.Itoa(year)+".pdf")
	err = pdf.Output(w)
	if err != nil {
		http.Error(w, "Erreur envoi PDF", http.StatusInternalServerError)
		return
	}
}
