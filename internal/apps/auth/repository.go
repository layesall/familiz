package auth

import (
	"database/sql"
	"familiz/internal/database"
	"familiz/internal/utils"
)

// CreateUser MODIFIÉ : on passe la transaction en paramètre
func CreateUser(tx *sql.Tx, req RegisterRequest, memberID int64) (int64, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return 0, err
	}

	// Utiliser tx.Exec au lieu de database.DB.Exec
	result, err := tx.Exec(`
        INSERT INTO users (email, password_hash, role, member_id, created_at, updated_at)
        VALUES (?, ?, 'admin', ?, datetime('now'), datetime('now'))
    `, req.Email, hashedPassword, memberID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func FindUserByEmail(email string) (*User, error) {
	var user User
	err := database.DB.QueryRow(`
        SELECT id, email, password_hash, role, member_id 
        FROM users 
        WHERE email = ?
    `, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.MemberID)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
