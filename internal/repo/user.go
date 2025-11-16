package repo

import (
	"database/sql"

	"github.com/Wucop228/avito-PullRequest/internal/models"
)

func UpdateUserIsActive(db *sql.DB, userID string, isActive bool) (*models.User, error) {
	query := "UPDATE users SET is_active = $2 WHERE id = $1 RETURNING id, username, team_name, is_active"

	user := &models.User{}
	err := db.QueryRow(query, userID, isActive).Scan(
		&user.UserID,
		&user.Username,
		&user.TeamName,
		&user.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
