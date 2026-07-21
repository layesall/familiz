package reports

import (
	"familiz/internal/database"
)

// MemberSummary pour le rapport global
type MemberSummary struct {
	ID          int     `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	TotalPaid   float64 `json:"total_paid"`
	EventsCount int     `json:"events_count"`
}

// GetGlobalReportData récupère les données pour le rapport global (par année)
func GetGlobalReportData(year int) ([]MemberSummary, error) {
	query := `
        SELECT 
            m.id,
            m.first_name,
            m.last_name,
            COALESCE(SUM(t.amount), 0) as total_paid,
            COALESCE(COUNT(e.id), 0) as events_count
        FROM members m
        LEFT JOIN transactions t ON t.member_id = m.id AND t.year = ? AND t.is_archived = 0
        LEFT JOIN events e ON e.member_id = m.id AND strftime('%Y', e.event_date) = ? AND e.is_archived = 0
        GROUP BY m.id
        ORDER BY m.last_name, m.first_name
    `
	rows, err := database.DB.Query(query, year, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []MemberSummary
	for rows.Next() {
		var s MemberSummary
		err := rows.Scan(&s.ID, &s.FirstName, &s.LastName, &s.TotalPaid, &s.EventsCount)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, nil
}
