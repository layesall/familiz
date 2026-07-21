package reports

import (
	"fmt"
	"strconv"
	"time"

	"familiz/internal/apps/events"
	"familiz/internal/apps/members"
	"familiz/internal/apps/transactions"

	"github.com/jung-kurt/gofpdf"
)

// generatePDFHeader ajoute l'en-tête commun à tous les PDF
func generatePDFHeader(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, title)
	pdf.Ln(12)

	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, "Generé le "+time.Now().Format("02/01/2006 15:04"))
	pdf.Ln(10)
}

// generateMemberProfilePDF génère le PDF pour un membre spécifique
func generateMemberProfilePDF(member *members.Member, txs []transactions.Transaction, evts []events.Event) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(15, 20, 15)

	// Titre
	generatePDFHeader(pdf, "Profil de "+member.FirstName+" "+member.LastName)

	// Informations personnelles
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Nom complet:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, member.FirstName+" "+member.LastName)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Date de naissance:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, member.BirthDate)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Statut marital:")
	pdf.SetFont("Arial", "", 12)
	status := map[string]string{"single": "Célibataire", "married": "Marié(e)", "minor": "Mineur"}[member.MaritalStatus]
	pdf.Cell(0, 8, status)
	pdf.Ln(12)

	// Transactions
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Historique des transactions")
	pdf.Ln(8)

	if len(txs) == 0 {
		pdf.SetFont("Arial", "I", 11)
		pdf.Cell(0, 8, "Aucune transaction enregistrée")
		pdf.Ln(8)
	} else {
		// Entête tableau
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(20, 8, "Mois")
		pdf.Cell(25, 8, "Année")
		pdf.Cell(45, 8, "Montant (€)")
		pdf.Cell(40, 8, "Payé le")
		pdf.Cell(0, 8, "Note")
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 10)
		for _, t := range txs {
			pdf.Cell(20, 8, strconv.Itoa(t.Month))
			pdf.Cell(25, 8, strconv.Itoa(t.Year))
			pdf.Cell(45, 8, fmt.Sprintf("%.2f", t.Amount))
			pdf.Cell(40, 8, t.PaidAt.Format("02/01/2006"))
			pdf.Cell(0, 8, t.Note)
			pdf.Ln(8)
		}
	}

	pdf.Ln(8)

	// Événements
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Événements")
	pdf.Ln(8)

	if len(evts) == 0 {
		pdf.SetFont("Arial", "I", 11)
		pdf.Cell(0, 8, "Aucun événement enregistré")
		pdf.Ln(8)
	} else {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(45, 8, "Type")
		pdf.Cell(45, 8, "Montant reçu (€)")
		pdf.Cell(40, 8, "Date")
		pdf.Cell(0, 8, "Archivé")
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 10)
		for _, e := range evts {
			typeLabel := map[string]string{"wedding": "Mariage", "baptism": "Baptême"}[e.Type]
			pdf.Cell(45, 8, typeLabel)
			pdf.Cell(45, 8, fmt.Sprintf("%.2f", e.AmountReceived))
			pdf.Cell(40, 8, e.EventDate)
			archived := "Non"
			if e.IsArchived {
				archived = "Oui"
			}
			pdf.Cell(0, 8, archived)
			pdf.Ln(8)
		}
	}

	return pdf, nil
}

// generateGlobalReportPDF génère le rapport global (par année)
func generateGlobalReportPDF(summaries []MemberSummary, year int) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(15, 20, 15)

	generatePDFHeader(pdf, fmt.Sprintf("Rapport annuel %d", year))

	if len(summaries) == 0 {
		pdf.SetFont("Arial", "I", 12)
		pdf.Cell(0, 10, "Aucune donnée pour cette année")
		return pdf, nil
	}

	// Entête tableau
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(30, 8, "ID")
	pdf.Cell(60, 8, "Membre")
	pdf.Cell(50, 8, "Total payé (€)")
	pdf.Cell(45, 8, "Nb événements")
	pdf.Ln(8)

	// Données
	pdf.SetFont("Arial", "", 10)
	for _, s := range summaries {
		pdf.Cell(30, 8, strconv.Itoa(s.ID))
		pdf.Cell(60, 8, s.FirstName+" "+s.LastName)
		pdf.Cell(50, 8, fmt.Sprintf("%.2f", s.TotalPaid))
		pdf.Cell(45, 8, strconv.Itoa(s.EventsCount))
		pdf.Ln(8)
	}

	// Total général
	var totalPaid float64
	for _, s := range summaries {
		totalPaid += s.TotalPaid
	}
	pdf.Ln(8)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Total collecté : %.2f €", totalPaid))

	return pdf, nil
}
