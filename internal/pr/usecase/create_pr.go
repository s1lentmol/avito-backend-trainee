package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/silentmol/avito-backend-trainee/internal/apperr"
	prdomain "github.com/silentmol/avito-backend-trainee/internal/pr/domain"
	"github.com/silentmol/avito-backend-trainee/internal/pr/dto"
)

func (u *PRUsecase) CreatePR(
	ctx context.Context,
	request *dto.CreatePRRequest,
) (*dto.CreatePRResponse, error) {

	// читаем автора PR и его команду
	author, err := u.userReader.GetUser(ctx, request.AuthorId)
	if err != nil {
		if err == apperr.ErrNotFound {
			slog.Info("PRUsecase.CreatePR: author not found",
				slog.String("pr_id", request.PrID),
				slog.String("author_id", request.AuthorId),
			)
			return nil, apperr.ErrNotFound
		}
		slog.Error("PRUsecase.CreatePR: failed to get author",
			slog.String("pr_id", request.PrID),
			slog.String("author_id", request.AuthorId),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("get author from provider: %w", err)
	}

	team, err := u.teamReader.GetTeam(ctx, author.TeamName)
	if err != nil {
		if err == apperr.ErrNotFound {
			slog.Info("PRUsecase.CreatePR: team not found",
				slog.String("pr_id", request.PrID),
				slog.String("team_name", author.TeamName),
			)
			return nil, apperr.ErrNotFound
		}
		slog.Error("PRUsecase.CreatePR: failed to get team",
			slog.String("pr_id", request.PrID),
			slog.String("team_name", author.TeamName),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("get team from provider: %w", err)
	}

	reviewers := prdomain.SelectReviewersForTeam(team, author.ID)

	pr := &prdomain.PullRequest{
		ID:                request.PrID,
		Name:              request.Name,
		AuthorId:          request.AuthorId,
		AssignedReviewers: reviewers,
	}

	created, err := u.prProvider.CreatePR(ctx, pr)
	if err != nil {
		slog.Error("PRUsecase.CreatePR: provider error",
			slog.String("pr_id", request.PrID),
			slog.String("author_id", request.AuthorId),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("create pull request in provider: %w", err)
	}

	slog.Info("PRUsecase.CreatePR: pull request created",
		slog.String("pr_id", created.ID),
		slog.String("author_id", created.AuthorId),
		slog.String("status", string(created.Status)),
	)

	return &dto.CreatePRResponse{
		PullRequest: *created,
	}, nil
}
