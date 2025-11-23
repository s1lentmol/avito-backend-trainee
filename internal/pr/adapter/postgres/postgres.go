package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/silentmol/avito-backend-trainee/internal/apperr"
	"github.com/silentmol/avito-backend-trainee/internal/pr/domain"
)

type PRRepository struct {
	conn *pgxpool.Pool
}

func NewPRRepository(conn *pgxpool.Pool) *PRRepository {
	return &PRRepository{conn: conn}
}

func (p *PRRepository) GetPR(ctx context.Context, id string) (*domain.PullRequest, error) {
	query := `
		SELECT id, name, author_id, status, reviewer1_id, reviewer2_id, created_at, merged_at
		FROM pull_requests
		WHERE id = $1
	`

	var pr domain.PullRequest
	var reviewer1, reviewer2 *string

	if err := p.conn.QueryRow(ctx, query, id).Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorId,
		&pr.Status,
		&reviewer1,
		&reviewer2,
		&pr.CreatedAt,
		&pr.MergedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("db: failed to get pull request: %w", err)
	}

	pr.AssignedReviewers = make([]string, 0, 2)
	if reviewer1 != nil {
		pr.AssignedReviewers = append(pr.AssignedReviewers, *reviewer1)
	}
	if reviewer2 != nil {
		pr.AssignedReviewers = append(pr.AssignedReviewers, *reviewer2)
	}

	return &pr, nil
}

func (p *PRRepository) UpdatePR(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
	var reviewer1, reviewer2 *string
	if len(pr.AssignedReviewers) > 0 {
		reviewer1 = &pr.AssignedReviewers[0]
	}
	if len(pr.AssignedReviewers) > 1 {
		reviewer2 = &pr.AssignedReviewers[1]
	}

	query := `
		UPDATE pull_requests
		SET name = $2,
		    status = $3,
		    reviewer1_id = $4,
		    reviewer2_id = $5,
		    merged_at = $6
		WHERE id = $1
		RETURNING id, name, author_id, status, reviewer1_id, reviewer2_id, created_at, merged_at
	`

	var updated domain.PullRequest
	var dbReviewer1, dbReviewer2 *string

	if err := p.conn.QueryRow(
		ctx,
		query,
		pr.ID,
		pr.Name,
		pr.Status,
		reviewer1,
		reviewer2,
		pr.MergedAt,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.AuthorId,
		&updated.Status,
		&dbReviewer1,
		&dbReviewer2,
		&updated.CreatedAt,
		&updated.MergedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("db: failed to update pull request: %w", err)
	}

	updated.AssignedReviewers = make([]string, 0, 2)
	if dbReviewer1 != nil {
		updated.AssignedReviewers = append(updated.AssignedReviewers, *dbReviewer1)
	}
	if dbReviewer2 != nil {
		updated.AssignedReviewers = append(updated.AssignedReviewers, *dbReviewer2)
	}

	return &updated, nil
}

func (p *PRRepository) CreatePR(ctx context.Context, pullRequest *domain.PullRequest) (*domain.PullRequest, error) {
	query := `
		INSERT INTO pull_requests (id, name, author_id, status, reviewer1_id, reviewer2_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, author_id, status, reviewer1_id, reviewer2_id, created_at, merged_at
	`

	var reviewer1, reviewer2 *string
	var createdPR domain.PullRequest

	var r1, r2 *string
	if len(pullRequest.AssignedReviewers) > 0 {
		r1 = &pullRequest.AssignedReviewers[0]
	}
	if len(pullRequest.AssignedReviewers) > 1 {
		r2 = &pullRequest.AssignedReviewers[1]
	}

	if err := p.conn.QueryRow(
		ctx,
		query,
		pullRequest.ID,
		pullRequest.Name,
		pullRequest.AuthorId,
		domain.StatusOpen,
		r1,
		r2,
	).Scan(
		&createdPR.ID,
		&createdPR.Name,
		&createdPR.AuthorId,
		&createdPR.Status,
		&reviewer1,
		&reviewer2,
		&createdPR.CreatedAt,
		&createdPR.MergedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, apperr.ErrPRExists
		}
		return nil, fmt.Errorf("db: failed to create pull request: %w", err)
	}

	createdPR.AssignedReviewers = make([]string, 0, 2)
	if reviewer1 != nil {
		createdPR.AssignedReviewers = append(createdPR.AssignedReviewers, *reviewer1)
	}
	if reviewer2 != nil {
		createdPR.AssignedReviewers = append(createdPR.AssignedReviewers, *reviewer2)
	}

	return &createdPR, nil
}

func (p *PRRepository) MergePR(ctx context.Context, id string) (*domain.PullRequest, error) {
	pr, err := p.GetPR(ctx, id)
	if err != nil {
		return nil, err
	}

	pr.Merge()

	updated, err := p.UpdatePR(ctx, pr)
	if err != nil {
		return nil, fmt.Errorf("db: failed to merge pull request: %w", err)
	}

	return updated, nil
}

func (p *PRRepository) ReassignPR(ctx context.Context, prId string, oldReviewerId string) (*domain.PullRequest, string, error) {
	pr, err := p.GetPR(ctx, prId)
	if err != nil {
		return nil, "", err
	}

	return pr, "", nil
}

func (p *PRRepository) GetReview(ctx context.Context, userId string) (*[]domain.PullRequest, error) {
	query := `
		SELECT id, name, author_id, status, reviewer1_id, reviewer2_id, created_at, merged_at
		FROM pull_requests
		WHERE reviewer1_id = $1 OR reviewer2_id = $1
	`

	rows, err := p.conn.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("db: failed to get pull requests for review: %w", err)
	}
	defer rows.Close()

	pullRequests := make([]domain.PullRequest, 0)

	for rows.Next() {
		var pr domain.PullRequest
		var reviewer1, reviewer2 *string

		if err := rows.Scan(
			&pr.ID,
			&pr.Name,
			&pr.AuthorId,
			&pr.Status,
			&reviewer1,
			&reviewer2,
			&pr.CreatedAt,
			&pr.MergedAt,
		); err != nil {
			return nil, fmt.Errorf("db: failed to scan pull request: %w", err)
		}

		pr.AssignedReviewers = make([]string, 0, 2)
		if reviewer1 != nil {
			pr.AssignedReviewers = append(pr.AssignedReviewers, *reviewer1)
		}
		if reviewer2 != nil {
			pr.AssignedReviewers = append(pr.AssignedReviewers, *reviewer2)
		}

		pullRequests = append(pullRequests, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: rows error: %w", err)
	}

	return &pullRequests, nil
}
