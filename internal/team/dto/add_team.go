package dto

import "github.com/silentmol/avito-backend-trainee/internal/team/domain"

type AddTeamRequest struct {
	domain.Team
}

type AddTeamResponse struct {
	Team domain.Team `json:"team"`
}
