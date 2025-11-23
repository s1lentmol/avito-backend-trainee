package dto

import (
	"github.com/silentmol/avito-backend-trainee/internal/pr/domain"
)

type GetReviewRequest struct {
	UserId string `query:"user_id" validate:"required"`
}

type ReviewPullRequest struct {
	ID       string           `json:"pull_request_id"`
	Name     string           `json:"pull_request_name"`
	AuthorId string           `json:"author_id"`
	Status   domain.PrStatus  `json:"status"`
}

type GetReviewResponse struct {
	UserId       string               `json:"user_id"`
	PullRequests []ReviewPullRequest  `json:"pull_requests"`
}
