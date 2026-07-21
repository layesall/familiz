package events

import "familiz/internal/database"

// GetEventsByMemberID récupère tous les événements d'un membre
func GetEventsByMemberID(memberID int) ([]Event, error) {
	rows, err := database.DB.Query(`
        SELECT id, member_id, type, amount_received, event_date, created_at, updated_at
        FROM events
        WHERE member_id = ?
        ORDER BY event_date DESC
    `, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evts []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.MemberID, &e.Type, &e.AmountReceived, &e.EventDate, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		evts = append(evts, e)
	}
	return evts, nil
}
