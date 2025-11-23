package usecase

import (
	"context"

	"github.com/silentmol/avito-backend-trainee/internal/team/domain"
)

type TeamProvider interface {
	CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error)
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
}

type TeamUsecase struct {
	teamProvider TeamProvider
}

func NewTeamUsecase(repo TeamProvider) *TeamUsecase {
	return &TeamUsecase{
		teamProvider: repo,
	}
}
