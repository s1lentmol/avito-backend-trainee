package domain

type User struct {
	ID       string `json:"user_id" validate:"required"`
	Name     string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}
