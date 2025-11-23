package usecase

import (
	"context"
	"fmt"

	"github.com/silentmol/avito-backend-trainee/internal/pr/dto"
)

func (u *PRUsecase) GetReview(ctx context.Context,
	request *dto.GetReviewRequest) (*dto.GetReviewResponse, error) {

	prs, err := u.prProvider.GetReview(ctx, request.UserId)
	if err != nil {
		return nil, fmt.Errorf("get pull requests for review in provider: %w", err)
	}

	respPRs := make([]dto.ReviewPullRequest, 0, len(*prs))
	for _, pr := range *prs {
		respPRs = append(respPRs, dto.ReviewPullRequest{
			ID:       pr.ID,
			Name:     pr.Name,
			AuthorId: pr.AuthorId,
			Status:   pr.Status,
		})
	}

	return &dto.GetReviewResponse{
		UserId:       request.UserId,
		PullRequests: respPRs,
	}, nil
}
