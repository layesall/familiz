package events

import (
	"database/sql"
	"errors"
	"familiz/internal/database"
)

// --- CREATE ---
func CreateEventRepo(memberID int, eventType string, amountReceived float64, eventDate string) (int64, error) {
	result, err := database.DB.Exec(`
        INSERT INTO events (member_id, type, amount_received, event_date, created_at, updated_at)
        VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
    `, memberID, eventType, amountReceived, eventDate)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// --- READ (by Member ID) ---
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

// --- READ (ALL) NOUVEAU ---
func GetAllEvents() ([]Event, error) {
	rows, err := database.DB.Query(`
        SELECT id, member_id, type, amount_received, event_date, created_at, updated_at
        FROM events
        ORDER BY event_date DESC, id DESC
    `)
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

// --- READ (by ID) NOUVEAU (pour vérifier l'existence) ---
func GetEventByID(id int) (*Event, error) {
	var e Event
	err := database.DB.QueryRow(`
        SELECT id, member_id, type, amount_received, event_date, created_at, updated_at
        FROM events
        WHERE id = ?
    `, id).Scan(&e.ID, &e.MemberID, &e.Type, &e.AmountReceived, &e.EventDate, &e.CreatedAt, &e.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// --- UPDATE (NOUVEAU) ---
func UpdateEventRepo(id int, req UpdateEventRequest) error {
	result, err := database.DB.Exec(`
        UPDATE events
        SET type = ?,
            amount_received = ?,
            event_date = ?,
            updated_at = datetime('now')
        WHERE id = ?
    `, req.Type, req.AmountReceived, req.EventDate, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("événement introuvable")
	}
	return nil
}

// --- DELETE (NOUVEAU) ---
func DeleteEventRepo(id int) error {
	result, err := database.DB.Exec("DELETE FROM events WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("événement introuvable")
	}
	return nil
}
