package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/silentmol/avito-backend-trainee/internal/apperr"
	prdomain "github.com/silentmol/avito-backend-trainee/internal/pr/domain"
	"github.com/silentmol/avito-backend-trainee/internal/pr/dto"
)

func (u *PRUsecase) ReassignPR(ctx context.Context,
	request *dto.ReassignPRRequest) (*dto.ReassignPRResponse, error) {

	// загружаем PR и сразу проверяем, можно ли его переназначать
	pr, err := u.prProvider.GetPR(ctx, request.PrID)
	if err != nil {
		if err == apperr.ErrNotFound {
			slog.Info("PRUsecase.ReassignPR: PR not found",
				slog.String("pr_id", request.PrID),
			)
			return nil, apperr.ErrNotFound
		}
		slog.Error("PRUsecase.ReassignPR: failed to get PR",
			slog.String("pr_id", request.PrID),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("get pull request from provider: %w", err)
	}

	// ранняя проверка, чтобы не ходить за пользователем/командой при уже слитом PR
	if err := pr.CanReassign(); err != nil {
		slog.Info("PRUsecase.ReassignPR: cannot reassign",
			slog.String("pr_id", request.PrID),
			slog.Any("error", err),
		)
		return nil, err
	}

	oldReviewer, err := u.userReader.GetUser(ctx, request.OldReviewerId)
	if err != nil {
		if err == apperr.ErrNotFound {
			slog.Info("PRUsecase.ReassignPR: old reviewer not found",
				slog.String("pr_id", request.PrID),
				slog.String("old_reviewer_id", request.OldReviewerId),
			)
			return nil, apperr.ErrNotFound
		}
		slog.Error("PRUsecase.ReassignPR: failed to get old reviewer",
			slog.String("pr_id", request.PrID),
			slog.String("old_reviewer_id", request.OldReviewerId),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("get old reviewer from provider: %w", err)
	}

	team, err := u.teamReader.GetTeam(ctx, oldReviewer.TeamName)
	if err != nil {
		if err == apperr.ErrNotFound {
			slog.Info("PRUsecase.ReassignPR: team not found",
				slog.String("pr_id", request.PrID),
				slog.String("team_name", oldReviewer.TeamName),
			)
			return nil, apperr.ErrNotFound
		}
		slog.Error("PRUsecase.ReassignPR: failed to get team",
			slog.String("pr_id", request.PrID),
			slog.String("team_name", oldReviewer.TeamName),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("get team from provider: %w", err)
	}

	newReviewerID, err := prdomain.ReassignReviewer(pr, team, request.OldReviewerId)
	if err != nil {
		slog.Info("PRUsecase.ReassignPR: cannot find replacement",
			slog.String("pr_id", request.PrID),
			slog.String("old_reviewer_id", request.OldReviewerId),
			slog.Any("error", err),
		)
		return nil, err
	}

	updated, err := u.prProvider.UpdatePR(ctx, pr)
	if err != nil {
		slog.Error("PRUsecase.ReassignPR: failed to update PR",
			slog.String("pr_id", request.PrID),
			slog.String("new_reviewer_id", newReviewerID),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("update pull request in provider: %w", err)
	}

	slog.Info("PRUsecase.ReassignPR: reviewer reassigned",
		slog.String("pr_id", updated.ID),
		slog.String("old_reviewer_id", request.OldReviewerId),
		slog.String("new_reviewer_id", newReviewerID),
	)

	return &dto.ReassignPRResponse{
		PullRequest: *updated,
		ReplacedBy:  newReviewerID,
	}, nil
}
