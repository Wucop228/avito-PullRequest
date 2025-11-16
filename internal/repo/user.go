package repo

import (
	"database/sql"
	"errors"

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

func GetUserByID(db *sql.DB, userID string) (*models.User, error) {
	query := "SELECT id, username, team_name, is_active FROM users WHERE id = $1"

	user := &models.User{}
	err := db.QueryRow(query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.TeamName,
		&user.IsActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func GetActiveUsersByTeam(db *sql.DB, teamName string) ([]models.User, error) {
	query := "SELECT id, username, team_name, is_active FROM users WHERE team_name = $1 AND is_active = TRUE"

	rows, err := db.Query(query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
