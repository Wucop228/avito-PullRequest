package service

import (
	"database/sql"
	"errors"

	"github.com/Wucop228/avito-PullRequest/internal/models"
	"github.com/Wucop228/avito-PullRequest/internal/repo"
)

var (
	ErrTeamExists   = errors.New("team_name already exists")
	ErrTeamNotFound = errors.New("team_name not found")
)

type TeamService struct {
	db *sql.DB
}

func NewTeamService(db *sql.DB) *TeamService {
	return &TeamService{db: db}
}

func (s *TeamService) CreateTeamWithMembers(req *models.RequestTeamAdd) error {
	team, err := repo.GetTeamByName(s.db, req.TeamName)
	if err != nil {
		return err
	}

	if team != nil && team.Name != "" {
		return ErrTeamExists
	}

	if err := repo.CreateTeamWithMembers(s.db, req); err != nil {
		return err
	}

	return nil
}

func (s *TeamService) GetTeam(name string) (*models.RequestTeamAdd, error) {
	team, err := repo.GetTeamWithMembers(s.db, name)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, ErrTeamNotFound
	}
	return team, nil
}
