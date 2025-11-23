package usecase

import (
	"context"

	"github.com/silentmol/avito-backend-trainee/internal/pr/domain"
	teamdomain "github.com/silentmol/avito-backend-trainee/internal/team/domain"
	userdomain "github.com/silentmol/avito-backend-trainee/internal/user/domain"
)

type UserReader interface {
	GetUser(ctx context.Context, id string) (*userdomain.User, error)
}

type TeamReader interface {
	GetTeam(ctx context.Context, teamName string) (*teamdomain.Team, error)
}

type PRProvider interface {
	GetPR(ctx context.Context, id string) (*domain.PullRequest, error)
	UpdatePR(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error)
	CreatePR(ctx context.Context, pullRequest *domain.PullRequest) (*domain.PullRequest, error)
	MergePR(ctx context.Context, id string) (*domain.PullRequest, error)
	GetReview(ctx context.Context, userId string) (*[]domain.PullRequest, error)
}

type PRUsecase struct {
	prProvider PRProvider
	userReader UserReader
	teamReader TeamReader
}

func NewPRUsecase(repo PRProvider, userReader UserReader, teamReader TeamReader) *PRUsecase {
	return &PRUsecase{
		prProvider: repo,
		userReader: userReader,
		teamReader: teamReader,
	}
}
