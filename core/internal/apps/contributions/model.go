package contributions

import "time"

type Contribution struct {
	ID         int       `json:"id"`
	MemberID   int       `json:"member_id"`
	Month      int       `json:"month"`
	Year       int       `json:"year"`
	Amount     float64   `json:"amount"`
	Note       string    `json:"note"`
	PaidAt     time.Time `json:"paid_at"`
	IsArchived bool      `json:"is_archived"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateContributionRequest struct {
	MemberID int     `json:"member_id"`
	Month    int     `json:"month"`
	Year     int     `json:"year"`
	Amount   float64 `json:"amount"`
	Note     string  `json:"note"`
}

type UpdateContributionRequest struct {
	Month  int     `json:"month"`
	Year   int     `json:"year"`
	Amount float64 `json:"amount"`
	Note   string  `json:"note"`
}
