package events

import (
	"database/sql"
	"errors"
	"familiz/internal/database"
)

func CreateEventRepo(memberID int, eventType string, amountReceived float64, eventDate string) (int64, error) {
	result, err := database.DB.Exec(`
        INSERT INTO events (member_id, type, amount_received, event_date, created_at, updated_at, is_archived)
        VALUES (?, ?, ?, ?, datetime('now'), datetime('now'), 0)
    `, memberID, eventType, amountReceived, eventDate)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetEventsByMemberID(memberID int, includeArchived bool) ([]Event, error) {
	query := `
        SELECT id, member_id, type, amount_received, event_date, created_at, updated_at, is_archived
        FROM events
        WHERE member_id = ?
    `
	if !includeArchived {
		query += " AND is_archived = 0"
	}
	query += " ORDER BY event_date DESC"

	rows, err := database.DB.Query(query, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evts []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.MemberID, &e.Type, &e.AmountReceived, &e.EventDate, &e.CreatedAt, &e.UpdatedAt, &e.IsArchived)
		if err != nil {
			return nil, err
		}
		evts = append(evts, e)
	}
	return evts, nil
}

func GetAllEvents(includeArchived bool) ([]Event, error) {
	query := `
        SELECT id, member_id, type, amount_received, event_date, created_at, updated_at, is_archived
        FROM events
    `
	if !includeArchived {
		query += " WHERE is_archived = 0"
	}
	query += " ORDER BY event_date DESC, id DESC"

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evts []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.MemberID, &e.Type, &e.AmountReceived, &e.EventDate, &e.CreatedAt, &e.UpdatedAt, &e.IsArchived)
		if err != nil {
			return nil, err
		}
		evts = append(evts, e)
	}
	return evts, nil
}

func GetEventByID(id int) (*Event, error) {
	var e Event
	err := database.DB.QueryRow(`
        SELECT id, member_id, type, amount_received, event_date, created_at, updated_at, is_archived
        FROM events
        WHERE id = ?
    `, id).Scan(&e.ID, &e.MemberID, &e.Type, &e.AmountReceived, &e.EventDate, &e.CreatedAt, &e.UpdatedAt, &e.IsArchived)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func UpdateEventRepo(id int, req UpdateEventRequest) error {
	var isArchived bool
	err := database.DB.QueryRow("SELECT is_archived FROM events WHERE id = ?", id).Scan(&isArchived)
	if err != nil {
		return err
	}
	if isArchived {
		return errors.New("impossible de modifier un événement archivé")
	}

	result, err := database.DB.Exec(`
        UPDATE events
        SET type = ?, amount_received = ?, event_date = ?, updated_at = datetime('now')
        WHERE id = ? AND is_archived = 0
    `, req.Type, req.AmountReceived, req.EventDate, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("événement introuvable ou déjà archivé")
	}
	return nil
}

func DeleteEventRepo(id int) error {
	var isArchived bool
	err := database.DB.QueryRow("SELECT is_archived FROM events WHERE id = ?", id).Scan(&isArchived)
	if err != nil {
		return err
	}
	if isArchived {
		return errors.New("impossible de supprimer un événement archivé")
	}

	result, err := database.DB.Exec("DELETE FROM events WHERE id = ? AND is_archived = 0", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("événement introuvable ou déjà archivé")
	}
	return nil
}
