package dto

import "github.com/silentmol/avito-backend-trainee/internal/team/domain"

type GetTeamRequest struct {
	TeamName string `query:"team_name" validate:"required"`
}

type GetTeamResponse struct {
	domain.Team
}
