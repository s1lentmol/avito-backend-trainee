package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/silentmol/avito-backend-trainee/internal/pr/dto"
)

func (u *PRUsecase) MergePR(ctx context.Context,
	request *dto.MergePRRequest) (*dto.MergePRResponse, error) {

	merged, err := u.prProvider.MergePR(ctx, request.PrID)
	if err != nil {
		slog.Error("PRUsecase.MergePR: provider error",
			slog.String("pr_id", request.PrID),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("merge pull request in provider: %w", err)
	}

	slog.Info("PRUsecase.MergePR: pull request merged",
		slog.String("pr_id", merged.ID),
		slog.String("status", string(merged.Status)),
	)

	return &dto.MergePRResponse{
		PullRequest: *merged,
	}, nil
}
