package members

import (
	"database/sql"
)

// CreateMemberTx insère un membre dans une transaction existante (utilisé par auth)
func CreateMemberTx(tx *sql.Tx, firstName, lastName, birthDate, maritalStatus string) (int64, error) {
	result, err := tx.Exec(`
        INSERT INTO members (first_name, last_name, birth_date, marital_status, created_at, updated_at)
        VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
    `, firstName, lastName, birthDate, maritalStatus)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
