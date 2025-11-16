package models

type Teams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type RequestTeamAdd struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}
