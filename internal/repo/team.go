package repo

import (
	"database/sql"
	"errors"

	"github.com/Wucop228/avito-PullRequest/internal/models"
)

func GetTeamByName(db *sql.DB, name string) (*models.Teams, error) {
	query := "SELECT id, name FROM teams WHERE name=$1"

	team := &models.Teams{}
	err := db.QueryRow(query, name).Scan(&team.ID, &team.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return team, nil
}

func CreateTeamWithMembers(db *sql.DB, team *models.RequestTeamAdd) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "INSERT INTO teams (name) VALUES ($1)"
	_, err = tx.Exec(query, team.TeamName)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO users (id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET username = EXCLUDED.username,
			team_name = EXCLUDED.team_name,
			is_active = EXCLUDED.is_active
	`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, member := range team.Members {
		if _, err := stmt.Exec(member.UserID, member.Username, team.TeamName, member.IsActive); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func GetTeamWithMembers(db *sql.DB, name string) (*models.RequestTeamAdd, error) {
	team, err := GetTeamByName(db, name)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, nil
	}

	query := "SELECT id, username, is_active FROM users WHERE team_name=$1"
	rows, err := db.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]models.TeamMember, 0)
	for rows.Next() {
		var m models.TeamMember
		if err := rows.Scan(&m.UserID, &m.Username, &m.IsActive); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.RequestTeamAdd{
		TeamName: name,
		Members:  members,
	}, nil
}
