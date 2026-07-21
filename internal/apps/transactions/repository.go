package transactions

import (
	"database/sql"
	"errors"
	"familiz/internal/database"
)

// --- CREATE ---
func CreateTransactionRepo(memberID, month, year int, amount float64, note string) (int64, error) {
	result, err := database.DB.Exec(`
        INSERT INTO transactions (member_id, month, year, amount, note, paid_at, created_at)
        VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))
    `, memberID, month, year, amount, note)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// --- READ (by Member ID) ---
func GetTransactionsByMemberID(memberID int) ([]Transaction, error) {
	rows, err := database.DB.Query(`
        SELECT id, member_id, month, year, amount, note, paid_at, created_at
        FROM transactions
        WHERE member_id = ?
        ORDER BY year DESC, month DESC
    `, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.MemberID, &t.Month, &t.Year, &t.Amount, &t.Note, &t.PaidAt, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	return txs, nil
}

// --- READ (ALL) NOUVEAU ---
func GetAllTransactions() ([]Transaction, error) {
	rows, err := database.DB.Query(`
        SELECT id, member_id, month, year, amount, note, paid_at, created_at
        FROM transactions
        ORDER BY year DESC, month DESC, id DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.MemberID, &t.Month, &t.Year, &t.Amount, &t.Note, &t.PaidAt, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	return txs, nil
}

// --- READ (by ID) NOUVEAU (pour vérifier l'existence avant update/delete) ---
func GetTransactionByID(id int) (*Transaction, error) {
	var t Transaction
	err := database.DB.QueryRow(`
        SELECT id, member_id, month, year, amount, note, paid_at, created_at
        FROM transactions
        WHERE id = ?
    `, id).Scan(&t.ID, &t.MemberID, &t.Month, &t.Year, &t.Amount, &t.Note, &t.PaidAt, &t.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// --- UPDATE (NOUVEAU) ---
func UpdateTransactionRepo(id int, req UpdateTransactionRequest) error {
	result, err := database.DB.Exec(`
        UPDATE transactions
        SET month = ?,
            year = ?,
            amount = ?,
            note = ?
        WHERE id = ?
    `, req.Month, req.Year, req.Amount, req.Note, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("transaction introuvable")
	}
	return nil
}

// --- DELETE (NOUVEAU) ---
func DeleteTransactionRepo(id int) error {
	result, err := database.DB.Exec("DELETE FROM transactions WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("transaction introuvable")
	}
	return nil
}
