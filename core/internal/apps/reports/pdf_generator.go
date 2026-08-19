package reports

import (
	"fmt"
	"strconv"
	"time"

	"familiz/internal/apps/contributions"
	"familiz/internal/apps/events"
	"familiz/internal/apps/members"

	"github.com/jung-kurt/gofpdf"
)

// initFonts charge la police UTF-8 (DejaVu) pour le PDF.
// Placez les fichiers .ttf dans un dossier "fonts/" à la racine du projet.
func initFonts(pdf *gofpdf.Fpdf) {
	// Vous pouvez aussi utiliser AddUTF8FontFromBytes avec //go:embed
	pdf.AddUTF8Font("DejaVu", "", "fonts/DejaVuSans.ttf")
	pdf.AddUTF8Font("DejaVu", "B", "fonts/DejaVuSans-Bold.ttf")
	pdf.AddUTF8Font("DejaVu", "I", "fonts/DejaVuSans-Oblique.ttf")
}

// generatePDFHeader ajoute un en-tête élégant.
func generatePDFHeader(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFillColor(41, 128, 185) // bleu
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("DejaVu", "B", 18)
	pdf.CellFormat(0, 15, title, "", 1, "C", true, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(4)

	pdf.SetFont("DejaVu", "I", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.Cell(0, 8, "Généré le "+time.Now().Format("02/01/2006 à 15:04"))
	pdf.Ln(10)
}

// generateFooter ajoute un pied de page avec le numéro de page.
func generateFooter(pdf *gofpdf.Fpdf) {
	pdf.SetY(-15)
	pdf.SetFont("DejaVu", "I", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()), "", 0, "C", false, 0, "")
}

// generateMemberProfilePDF génère le PDF pour un membre.
func generateMemberProfilePDF(member *members.Member, txs []contributions.Contribution, evts []events.Event) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	initFonts(pdf)

	pdf.AddPage()
	generatePDFHeader(pdf, "Profil de "+member.FirstName+" "+member.LastName)

	// Informations personnelles
	pdf.SetFont("DejaVu", "B", 12)
	pdf.Cell(40, 8, "Nom complet :")
	pdf.SetFont("DejaVu", "", 12)
	pdf.Cell(0, 8, member.FirstName+" "+member.LastName)
	pdf.Ln(8)

	pdf.SetFont("DejaVu", "B", 12)
	pdf.Cell(40, 8, "Date de naissance :")
	pdf.SetFont("DejaVu", "", 12)
	pdf.Cell(0, 8, member.BirthDate)
	pdf.Ln(8)

	pdf.SetFont("DejaVu", "B", 12)
	pdf.Cell(40, 8, "Statut marital :")
	pdf.SetFont("DejaVu", "", 12)
	status := map[string]string{"single": "Célibataire", "married": "Marié(e)", "minor": "Mineur"}[member.MaritalStatus]
	pdf.Cell(0, 8, status)
	pdf.Ln(12)

	// Contributions
	pdf.SetFont("DejaVu", "B", 14)
	pdf.SetFillColor(230, 240, 255)
	pdf.CellFormat(0, 10, "Historique des contributions", "", 1, "L", true, 0, "")
	pdf.Ln(4)

	if len(txs) == 0 {
		pdf.SetFont("DejaVu", "I", 11)
		pdf.Cell(0, 8, "Aucune contribution enregistrée")
		pdf.Ln(8)
	} else {
		headers := []string{"Mois", "Année", "Montant (€)", "Payé le", "Note"}
		colWidths := []float64{20, 25, 45, 40, 0}
		pdf.SetFont("DejaVu", "B", 10)
		pdf.SetFillColor(200, 220, 240)
		for i, h := range headers {
			pdf.CellFormat(colWidths[i], 8, h, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(8)

		pdf.SetFont("DejaVu", "", 10)
		alternate := false
		for _, t := range txs {
			if alternate {
				pdf.SetFillColor(245, 245, 245)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}
			alternate = !alternate
			pdf.CellFormat(colWidths[0], 8, strconv.Itoa(t.Month), "1", 0, "C", true, 0, "")
			pdf.CellFormat(colWidths[1], 8, strconv.Itoa(t.Year), "1", 0, "C", true, 0, "")
			pdf.CellFormat(colWidths[2], 8, fmt.Sprintf("%.2f", t.Amount), "1", 0, "R", true, 0, "")
			pdf.CellFormat(colWidths[3], 8, t.PaidAt.Format("02/01/2006"), "1", 0, "C", true, 0, "")
			pdf.CellFormat(colWidths[4], 8, t.Note, "1", 0, "L", true, 0, "")
			pdf.Ln(8)
		}
	}

	pdf.Ln(8)

	// Événements
	pdf.SetFont("DejaVu", "B", 14)
	pdf.SetFillColor(230, 240, 255)
	pdf.CellFormat(0, 10, "Événements", "", 1, "L", true, 0, "")
	pdf.Ln(4)

	if len(evts) == 0 {
		pdf.SetFont("DejaVu", "I", 11)
		pdf.Cell(0, 8, "Aucun événement enregistré")
		pdf.Ln(8)
	} else {
		headers := []string{"Type", "Montant reçu (€)", "Date", "Archivé"}
		colWidths := []float64{45, 50, 40, 0}
		pdf.SetFont("DejaVu", "B", 10)
		pdf.SetFillColor(200, 220, 240)
		for i, h := range headers {
			pdf.CellFormat(colWidths[i], 8, h, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(8)

		pdf.SetFont("DejaVu", "", 10)
		alternate := false
		for _, e := range evts {
			if alternate {
				pdf.SetFillColor(245, 245, 245)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}
			alternate = !alternate
			typeLabel := map[string]string{"wedding": "Mariage", "baptism": "Baptême"}[e.Type]
			pdf.CellFormat(colWidths[0], 8, typeLabel, "1", 0, "L", true, 0, "")
			pdf.CellFormat(colWidths[1], 8, fmt.Sprintf("%.2f", e.AmountReceived), "1", 0, "R", true, 0, "")
			pdf.CellFormat(colWidths[2], 8, e.EventDate, "1", 0, "C", true, 0, "")
			archived := "Non"
			if e.IsArchived {
				archived = "Oui"
			}
			pdf.CellFormat(colWidths[3], 8, archived, "1", 0, "C", true, 0, "")
			pdf.Ln(8)
		}
	}

	generateFooter(pdf)
	return pdf, nil
}

// generateGlobalReportPDF génère le rapport global par année.
func generateGlobalReportPDF(summaries []MemberSummary, year int) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	initFonts(pdf)

	pdf.AddPage()
	generatePDFHeader(pdf, fmt.Sprintf("Rapport annuel %d", year))

	if len(summaries) == 0 {
		pdf.SetFont("DejaVu", "I", 12)
		pdf.Cell(0, 10, "Aucune donnée pour cette année")
		generateFooter(pdf)
		return pdf, nil
	}

	headers := []string{"ID", "Membre", "Total payé (€)", "Nb événements"}
	colWidths := []float64{20, 70, 50, 0}
	pdf.SetFont("DejaVu", "B", 10)
	pdf.SetFillColor(200, 220, 240)
	for i, h := range headers {
		pdf.CellFormat(colWidths[i], 8, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(8)

	pdf.SetFont("DejaVu", "", 10)
	var totalPaid float64
	alternate := false
	for _, s := range summaries {
		if alternate {
			pdf.SetFillColor(245, 245, 245)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		alternate = !alternate
		pdf.CellFormat(colWidths[0], 8, strconv.Itoa(s.ID), "1", 0, "C", true, 0, "")
		pdf.CellFormat(colWidths[1], 8, s.FirstName+" "+s.LastName, "1", 0, "L", true, 0, "")
		pdf.CellFormat(colWidths[2], 8, fmt.Sprintf("%.2f", s.TotalPaid), "1", 0, "R", true, 0, "")
		pdf.CellFormat(colWidths[3], 8, strconv.Itoa(s.EventsCount), "1", 0, "C", true, 0, "")
		pdf.Ln(8)
		totalPaid += s.TotalPaid
	}

	pdf.Ln(6)
	pdf.SetFont("DejaVu", "B", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Total collecté : %.2f €", totalPaid))

	generateFooter(pdf)
	return pdf, nil
}
