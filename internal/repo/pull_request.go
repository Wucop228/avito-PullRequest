package repo

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Wucop228/avito-PullRequest/internal/models"
)

func GetPullRequestWithReviewers(db *sql.DB, id string) (*models.PullRequest, error) {
	query := `
		SELECT id, name, author_id, status, created_at, merged_at
		FROM pull_requests
		WHERE id = $1
	`

	var pr models.PullRequest
	var createdAt time.Time
	var mergedAt sql.NullTime

	err := db.QueryRow(query, id).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&createdAt,
		&mergedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	pr.CreatedAt = &createdAt
	if mergedAt.Valid {
		m := mergedAt.Time
		pr.MergedAt = &m
	}

	rows, err := db.Query(
		`SELECT reviewer_id FROM pull_request_reviewers WHERE pull_request_id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviewers := make([]string, 0)
	for rows.Next() {
		var r string
		if err := rows.Scan(&r); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	pr.AssignedReviewers = reviewers

	return &pr, nil
}

func CreatePullRequest(db *sql.DB, req *models.RequestPullRequestCreate, reviewerIDs []string) (*models.PullRequest, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var createdAt time.Time
	insertPR := `
		INSERT INTO pull_requests (id, name, author_id, status)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`
	err = tx.QueryRow(
		insertPR,
		req.PullRequestID,
		req.PullRequestName,
		req.AuthorID,
		"OPEN",
	).Scan(&createdAt)
	if err != nil {
		return nil, err
	}

	if len(reviewerIDs) > 0 {
		insertReviewer := `
			INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id)
			VALUES ($1, $2)
		`
		for _, r := range reviewerIDs {
			if _, err := tx.Exec(insertReviewer, req.PullRequestID, r); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	pr := &models.PullRequest{
		PullRequestID:     req.PullRequestID,
		PullRequestName:   req.PullRequestName,
		AuthorID:          req.AuthorID,
		Status:            "OPEN",
		AssignedReviewers: reviewerIDs,
		CreatedAt:         &createdAt,
	}

	return pr, nil
}

func MarkPullRequestMerged(db *sql.DB, id string) (*time.Time, error) {
	query := `
		UPDATE pull_requests
		SET status = 'MERGED', merged_at = NOW()
		WHERE id = $1
		RETURNING merged_at
	`

	var mergedAt time.Time
	if err := db.QueryRow(query, id).Scan(&mergedAt); err != nil {
		return nil, err
	}

	return &mergedAt, nil
}

func ReplacePullRequestReviewer(db *sql.DB, prID, oldReviewerID, newReviewerID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(
		`DELETE FROM pull_request_reviewers WHERE pull_request_id = $1 AND reviewer_id = $2`,
		prID,
		oldReviewerID,
	); err != nil {
		return err
	}

	if _, err := tx.Exec(
		`INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2)`,
		prID,
		newReviewerID,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func GetPullRequestsByReviewer(db *sql.DB, userID string) ([]models.PullRequestShort, error) {
	query := `
		SELECT pr.id, pr.name, pr.author_id, pr.status
		FROM pull_requests pr
		JOIN pull_request_reviewers r ON pr.id = r.pull_request_id
		WHERE r.reviewer_id = $1
		ORDER BY pr.created_at
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prs := make([]models.PullRequestShort, 0)
	for rows.Next() {
		var pr models.PullRequestShort
		if err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}
