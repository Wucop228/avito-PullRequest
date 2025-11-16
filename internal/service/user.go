package service

import (
	"database/sql"
	"errors"

	"github.com/Wucop228/avito-PullRequest/internal/models"
	"github.com/Wucop228/avito-PullRequest/internal/repo"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) SetIsActive(userID string, isActive bool) (*models.User, error) {
	user, err := repo.UpdateUserIsActive(s.db, userID, isActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
