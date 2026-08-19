package contributions

import (
	"database/sql"
	"errors"
	"familiz/internal/database"
)

// --- CREATE ---
func CreateContributionRepo(memberID, month, year int, amount float64, note string) (int64, error) {
	result, err := database.DB.Exec(`
        INSERT INTO contributions (member_id, month, year, amount, note, paid_at, created_at, is_archived)
        VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'), 0)
    `, memberID, month, year, amount, note)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// --- READ (by Member ID) ---
func GetContributionsByMemberID(memberID int, includeArchived bool) ([]Contribution, error) {
	query := `
        SELECT id, member_id, month, year, amount, note, paid_at, is_archived, created_at
        FROM contributions
        WHERE member_id = ?
    `
	if !includeArchived {
		query += " AND is_archived = 0"
	}
	query += " ORDER BY year DESC, month DESC"

	rows, err := database.DB.Query(query, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contributions []Contribution
	for rows.Next() {
		var c Contribution
		err := rows.Scan(&c.ID, &c.MemberID, &c.Month, &c.Year, &c.Amount, &c.Note, &c.PaidAt, &c.IsArchived, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		contributions = append(contributions, c)
	}
	return contributions, nil
}

// --- READ (ALL) ---
func GetAllContributions(includeArchived bool) ([]Contribution, error) {
	query := `
        SELECT id, member_id, month, year, amount, note, paid_at, is_archived, created_at
        FROM contributions
    `
	if !includeArchived {
		query += " WHERE is_archived = 0"
	}
	query += " ORDER BY year DESC, month DESC, id DESC"

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contributions []Contribution
	for rows.Next() {
		var c Contribution
		err := rows.Scan(&c.ID, &c.MemberID, &c.Month, &c.Year, &c.Amount, &c.Note, &c.PaidAt, &c.IsArchived, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		contributions = append(contributions, c)
	}
	return contributions, nil
}

// --- READ (by ID) ---
func GetContributionByID(id int) (*Contribution, error) {
	var c Contribution
	err := database.DB.QueryRow(`
        SELECT id, member_id, month, year, amount, note, paid_at, is_archived, created_at
        FROM contributions
        WHERE id = ?
    `, id).Scan(&c.ID, &c.MemberID, &c.Month, &c.Year, &c.Amount, &c.Note, &c.PaidAt, &c.IsArchived, &c.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// --- UPDATE ---
func UpdateContributionRepo(id int, req UpdateContributionRequest) error {
	var isArchived bool
	err := database.DB.QueryRow("SELECT is_archived FROM contributions WHERE id = ?", id).Scan(&isArchived)
	if err != nil {
		return err
	}
	if isArchived {
		return errors.New("impossible de modifier une contribution archivée")
	}

	result, err := database.DB.Exec(`
        UPDATE contributions
        SET month = ?, year = ?, amount = ?, note = ?
        WHERE id = ? AND is_archived = 0
    `, req.Month, req.Year, req.Amount, req.Note, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("contribution introuvable ou déjà archivée")
	}
	return nil
}

// --- DELETE ---
func DeleteContributionRepo(id int) error {
	var isArchived bool
	err := database.DB.QueryRow("SELECT is_archived FROM contributions WHERE id = ?", id).Scan(&isArchived)
	if err != nil {
		return err
	}
	if isArchived {
		return errors.New("impossible de supprimer une contribution archivée")
	}

	result, err := database.DB.Exec("DELETE FROM contributions WHERE id = ? AND is_archived = 0", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("contribution introuvable ou déjà archivée")
	}
	return nil
}
