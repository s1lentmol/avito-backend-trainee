package dto

import (
	"github.com/silentmol/avito-backend-trainee/internal/pr/domain"
)

type MergePRRequest struct {
	PrID string `json:"pull_request_id" validate:"required"`
}

type MergePRResponse struct {
	PullRequest domain.PullRequest `json:"pr"`
}
