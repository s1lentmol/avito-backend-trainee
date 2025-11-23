package domain

import (
	"time"

	"github.com/silentmol/avito-backend-trainee/internal/apperr"
)

type PrStatus string

const (
	StatusOpen   PrStatus = "OPEN"
	StatusMerged PrStatus = "MERGED"
)

type PullRequest struct {
	ID                string     `json:"pull_request_id"`
	Name              string     `json:"pull_request_name"`
	AuthorId          string     `json:"author_id"`
	Status            PrStatus   `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         time.Time  `json:"createdAt"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

func (p *PullRequest) IsMerged() bool {
	return p.Status == StatusMerged
}

func (p *PullRequest) CanReassign() error {
	if p.IsMerged() {
		return apperr.ErrPRMerged
	}
	return nil
}

func (p *PullRequest) ReplaceReviewer(oldID, newID string) error {
	if oldID == newID {
		return nil
	}

	replaced := false
	for i, reviewerID := range p.AssignedReviewers {
		if reviewerID == oldID {
			p.AssignedReviewers[i] = newID
			replaced = true
			break
		}
	}

	if !replaced {
		return apperr.ErrNotAssigned
	}

	return nil
}

func (p *PullRequest) Merge() {
	if p.IsMerged() {
		return
	}

	p.Status = StatusMerged
	now := time.Now()
	p.MergedAt = &now
}
