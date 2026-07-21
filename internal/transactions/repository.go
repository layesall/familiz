package transactions

import "familiz/internal/database"

// GetTransactionsByMemberID récupère toutes les transactions d'un membre
func GetTransactionsByMemberID(memberID int) ([]Transaction, error) {
	rows, err := database.DB.Query(`
        SELECT id, member_id, month, year, amount, paid_at, note, created_at
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
		err := rows.Scan(&t.ID, &t.MemberID, &t.Month, &t.Year, &t.Amount, &t.PaidAt, &t.Note, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	return txs, nil
}
