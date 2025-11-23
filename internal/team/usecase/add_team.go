package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/silentmol/avito-backend-trainee/internal/team/domain"
	"github.com/silentmol/avito-backend-trainee/internal/team/dto"
)

func (t *TeamUsecase) CreateTeam(ctx context.Context, addTeamRequest *dto.AddTeamRequest) (*dto.AddTeamResponse, error) {
	team := &domain.Team{
		Name:    addTeamRequest.Name,
		Members: addTeamRequest.Members,
	}

	createdTeam, err := t.teamProvider.CreateTeam(ctx, team)

	if err != nil {
		slog.Error("TeamUsecase.CreateTeam: provider error",
			slog.String("team_name", team.Name),
			slog.Int("members_count", len(team.Members)),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("create team in postgres: %w", err)
	}

	slog.Info("TeamUsecase.CreateTeam: team created",
		slog.String("team_name", createdTeam.Name),
		slog.Int("members_count", len(createdTeam.Members)),
	)

	return &dto.AddTeamResponse{
		Team: *createdTeam,
	}, nil
}
