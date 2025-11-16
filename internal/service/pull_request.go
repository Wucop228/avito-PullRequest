package service

import (
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/Wucop228/avito-PullRequest/internal/models"
	"github.com/Wucop228/avito-PullRequest/internal/repo"
)

var (
	ErrPRExists            = errors.New("PR id already exists")
	ErrPRNotFound          = errors.New("pull request not found")
	ErrPRMerged            = errors.New("cannot reassign on merged PR")
	ErrReviewerNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate         = errors.New("no active replacement candidate in team")
	ErrAuthorNotFound      = errors.New("author not found")
)

type PullRequestService struct {
	db *sql.DB
}

func NewPullRequestService(db *sql.DB) *PullRequestService {
	rand.Seed(time.Now().UnixNano())
	return &PullRequestService{db: db}
}

func (s *PullRequestService) CreatePullRequest(req *models.RequestPullRequestCreate) (*models.PullRequest, error) {
	existing, err := repo.GetPullRequestWithReviewers(s.db, req.PullRequestID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrPRExists
	}

	author, err := repo.GetUserByID(s.db, req.AuthorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, ErrAuthorNotFound
	}

	teamUsers, err := repo.GetActiveUsersByTeam(s.db, author.TeamName)
	if err != nil {
		return nil, err
	}

	candidateIDs := make([]string, 0)
	for _, u := range teamUsers {
		if u.UserID == author.UserID {
			continue
		}
		candidateIDs = append(candidateIDs, u.UserID)
	}

	selected := selectRandomReviewers(candidateIDs, 2)

	return repo.CreatePullRequest(s.db, req, selected)
}

func (s *PullRequestService) MergePullRequest(prID string) (*models.PullRequest, error) {
	pr, err := repo.GetPullRequestWithReviewers(s.db, prID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, ErrPRNotFound
	}

	if pr.Status == "MERGED" {
		return pr, nil
	}

	mergedAt, err := repo.MarkPullRequestMerged(s.db, prID)
	if err != nil {
		return nil, err
	}

	pr.Status = "MERGED"
	pr.MergedAt = mergedAt

	return pr, nil
}

func (s *PullRequestService) ReassignReviewer(prID, oldUserID string) (*models.PullRequest, string, error) {
	pr, err := repo.GetPullRequestWithReviewers(s.db, prID)
	if err != nil {
		return nil, "", err
	}
	if pr == nil {
		return nil, "", ErrPRNotFound
	}

	if pr.Status == "MERGED" {
		return nil, "", ErrPRMerged
	}

	assigned := false
	for _, id := range pr.AssignedReviewers {
		if id == oldUserID {
			assigned = true
			break
		}
	}
	if !assigned {
		return nil, "", ErrReviewerNotAssigned
	}

	user, err := repo.GetUserByID(s.db, oldUserID)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", ErrUserNotFound
	}

	teamUsers, err := repo.GetActiveUsersByTeam(s.db, user.TeamName)
	if err != nil {
		return nil, "", err
	}

	exclude := make(map[string]struct{})
	exclude[oldUserID] = struct{}{}
	exclude[pr.AuthorID] = struct{}{}
	for _, id := range pr.AssignedReviewers {
		exclude[id] = struct{}{}
	}

	candidates := make([]string, 0)
	for _, u := range teamUsers {
		if _, ok := exclude[u.UserID]; ok {
			continue
		}
		candidates = append(candidates, u.UserID)
	}

	if len(candidates) == 0 {
		return nil, "", ErrNoCandidate
	}

	newReviewerID := candidates[rand.Intn(len(candidates))]

	if err := repo.ReplacePullRequestReviewer(s.db, prID, oldUserID, newReviewerID); err != nil {
		return nil, "", err
	}

	for i, id := range pr.AssignedReviewers {
		if id == oldUserID {
			pr.AssignedReviewers[i] = newReviewerID
			break
		}
	}

	return pr, newReviewerID, nil
}

func (s *PullRequestService) GetUserReviews(userID string) ([]models.PullRequestShort, error) {
	return repo.GetPullRequestsByReviewer(s.db, userID)
}

func selectRandomReviewers(ids []string, maxCount int) []string {
	n := len(ids)
	if n == 0 || maxCount <= 0 {
		return []string{}
	}
	if n <= maxCount {
		res := make([]string, n)
		copy(res, ids)
		return res
	}

	copyIDs := make([]string, n)
	copy(copyIDs, ids)
	rand.Shuffle(n, func(i, j int) {
		copyIDs[i], copyIDs[j] = copyIDs[j], copyIDs[i]
	})

	return copyIDs[:maxCount]
}
