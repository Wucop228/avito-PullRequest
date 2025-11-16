package models

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type RequestSetIsActive struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}
