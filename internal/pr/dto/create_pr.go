package dto

import (
	"github.com/silentmol/avito-backend-trainee/internal/pr/domain"
)

type CreatePRRequest struct {
	PrID     string `json:"pull_request_id" validate:"required"`
	Name     string `json:"pull_request_name" validate:"required"`
	AuthorId string `json:"author_id" validate:"required"`
}

type CreatePRResponse struct {
	PullRequest domain.PullRequest `json:"pr"`
}
