package dto

import (
	"github.com/silentmol/avito-backend-trainee/internal/pr/domain"
)

type ReassignPRRequest struct {
	PrID          string `json:"pull_request_id" validate:"required"`
	OldReviewerId string `json:"old_reviewer_id" validate:"required"`
}

type ReassignPRResponse struct {
	PullRequest domain.PullRequest `json:"pr"`
	ReplacedBy  string             `json:"replaced_by"`
}
