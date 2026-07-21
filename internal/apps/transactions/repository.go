package transactions

import (
	"database/sql"
	"errors"
	"familiz/internal/database"
)

// --- CREATE ---
func CreateTransactionRepo(memberID, month, year int, amount float64, note string) (int64, error) {
	result, err := database.DB.Exec(`
        INSERT INTO transactions (member_id, month, year, amount, note, paid_at, created_at, is_archived)
        VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'), 0)
    `, memberID, month, year, amount, note)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// --- READ (by Member ID) ---
func GetTransactionsByMemberID(memberID int, includeArchived bool) ([]Transaction, error) {
	// Construction de la requête avec ou sans le filtre d'archivage
	query := `
        SELECT id, member_id, month, year, amount, note, paid_at, is_archived, created_at
        FROM transactions
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

	var txs []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.MemberID, &t.Month, &t.Year, &t.Amount, &t.Note, &t.PaidAt, &t.IsArchived, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	return txs, nil
}

// --- READ (ALL) ---
func GetAllTransactions(includeArchived bool) ([]Transaction, error) {
	query := `
        SELECT id, member_id, month, year, amount, note, paid_at, is_archived, created_at
        FROM transactions
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

	var txs []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.MemberID, &t.Month, &t.Year, &t.Amount, &t.Note, &t.PaidAt, &t.IsArchived, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	return txs, nil
}

// --- READ (by ID) ---
func GetTransactionByID(id int) (*Transaction, error) {
	var t Transaction
	err := database.DB.QueryRow(`
        SELECT id, member_id, month, year, amount, note, paid_at, is_archived, created_at
        FROM transactions
        WHERE id = ?
    `, id).Scan(&t.ID, &t.MemberID, &t.Month, &t.Year, &t.Amount, &t.Note, &t.PaidAt, &t.IsArchived, &t.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// --- UPDATE ---
func UpdateTransactionRepo(id int, req UpdateTransactionRequest) error {
	// 1. Vérifier que ce n'est pas une archive
	var isArchived bool
	err := database.DB.QueryRow("SELECT is_archived FROM transactions WHERE id = ?", id).Scan(&isArchived)
	if err != nil {
		return err
	}
	if isArchived {
		return errors.New("impossible de modifier une transaction archivée")
	}

	// 2. Mise à jour
	result, err := database.DB.Exec(`
        UPDATE transactions
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
		return errors.New("transaction introuvable ou déjà archivée")
	}
	return nil
}

// --- DELETE ---
func DeleteTransactionRepo(id int) error {
	// 1. Vérifier que ce n'est pas une archive
	var isArchived bool
	err := database.DB.QueryRow("SELECT is_archived FROM transactions WHERE id = ?", id).Scan(&isArchived)
	if err != nil {
		return err
	}
	if isArchived {
		return errors.New("impossible de supprimer une transaction archivée")
	}

	// 2. Suppression
	result, err := database.DB.Exec("DELETE FROM transactions WHERE id = ? AND is_archived = 0", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("transaction introuvable ou déjà archivée")
	}
	return nil
}
