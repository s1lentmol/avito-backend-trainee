package usecase

import (
	"context"
	"fmt"

	"github.com/silentmol/avito-backend-trainee/internal/team/dto"
)

func (t *TeamUsecase) GetTeam(ctx context.Context, getTeamRequest *dto.GetTeamRequest) (*dto.GetTeamResponse, error) {
	teamName := getTeamRequest.TeamName

	team, err := t.teamProvider.GetTeam(ctx, teamName)

	if err != nil {
		return nil, fmt.Errorf("get team from postgres: %w", err)
	}

	return &dto.GetTeamResponse{
		Team: *team,
	}, nil
}
