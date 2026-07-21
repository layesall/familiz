package members

import (
	"database/sql"
	"errors"
	"familiz/internal/database"
)

// CreateMemberTx est utilisé par le package auth pour créer un membre lors de l'inscription
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

// GetMemberByID récupère un membre par son ID
func GetMemberByID(id int) (*Member, error) {
	var m Member
	err := database.DB.QueryRow(`
        SELECT id, first_name, last_name, birth_date, marital_status, created_at, updated_at
        FROM members
        WHERE id = ?
    `, id).Scan(&m.ID, &m.FirstName, &m.LastName, &m.BirthDate, &m.MaritalStatus, &m.CreatedAt, &m.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // Aucun membre trouvé, mais pas d'erreur critique
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GetAllMembers récupère tous les membres (pour la liste)
func GetAllMembers() ([]Member, error) {
	rows, err := database.DB.Query(`
        SELECT id, first_name, last_name, birth_date, marital_status, created_at, updated_at
        FROM members
        ORDER BY id DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var m Member
		err := rows.Scan(&m.ID, &m.FirstName, &m.LastName, &m.BirthDate, &m.MaritalStatus, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

// UpdateMember met à jour les champs d'un membre
func UpdateMember(id int, req UpdateMemberRequest) error {
	result, err := database.DB.Exec(`
        UPDATE members
        SET first_name = ?,
            last_name = ?,
            birth_date = ?,
            marital_status = ?,
            updated_at = datetime('now')
        WHERE id = ?
    `, req.FirstName, req.LastName, req.BirthDate, req.MaritalStatus, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("membre introuvable")
	}
	return nil
}

// DeleteMember supprime un membre par son ID
// (Les transactions et événements seront supprimés automatiquement grâce à ON DELETE CASCADE)
func DeleteMember(id int) error {
	result, err := database.DB.Exec("DELETE FROM members WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("membre introuvable")
	}
	return nil
}
