package members

// Member représente un membre de la famille
type Member struct {
	ID            int    `json:"id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	BirthDate     string `json:"birth_date"`
	MaritalStatus string `json:"marital_status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// CreateMemberRequest
type CreateMemberRequest struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	BirthDate     string `json:"birth_date"`
	MaritalStatus string `json:"marital_status"`
}

// UpdateMemberRequest
type UpdateMemberRequest struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	BirthDate     string `json:"birth_date"`
	MaritalStatus string `json:"marital_status"`
}
