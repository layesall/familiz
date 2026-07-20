package auth

// User représente l'utilisateur en base
type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
	MemberID     *int   `json:"member_id"`
}

// RegisterRequest correspond au JSON reçu pour l'inscription
type RegisterRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	BirthDate     string `json:"birth_date"`
	MaritalStatus string `json:"marital_status"`
}

// LoginRequest correspond au JSON reçu pour la connexion
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse correspond au JSON renvoyé après connexion
type LoginResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
