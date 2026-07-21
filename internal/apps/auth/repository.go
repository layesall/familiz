package auth

import (
	"database/sql"
	// "log"
	"familiz/internal/database"
	"familiz/internal/utils"
)

// CreateUser insère un nouvel utilisateur et son membre associé dans une transaction
func CreateUser(req RegisterRequest, memberID int64) (int64, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return 0, err
	}

	result, err := database.DB.Exec(`
        INSERT INTO users (email, password_hash, role, member_id, created_at, updated_at)
        VALUES (?, ?, 'admin', ?, datetime('now'), datetime('now'))
    `, req.Email, hashedPassword, memberID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// FindUserByEmail cherche un utilisateur par son email
func FindUserByEmail(email string) (*User, error) {
	var user User
	err := database.DB.QueryRow(`
        SELECT id, email, password_hash, role, member_id 
        FROM users 
        WHERE email = ?
    `, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.MemberID)

	if err == sql.ErrNoRows {
		return nil, nil // Aucun utilisateur trouvé, mais pas d'erreur critique
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
