package http

import (
	prusecase "github.com/silentmol/avito-backend-trainee/internal/pr/usecase"
	teamusecase "github.com/silentmol/avito-backend-trainee/internal/team/usecase"
	userusecase "github.com/silentmol/avito-backend-trainee/internal/user/usecase"
)

type Handle struct {
	user *userusecase.UserUsecase
	team *teamusecase.TeamUsecase
	pr   *prusecase.PRUsecase
}

func NewHandler(
	userUC *userusecase.UserUsecase,
	teamUC *teamusecase.TeamUsecase,
	prUC *prusecase.PRUsecase,
) *Handle {
	return &Handle{
		user: userUC,
		team: teamUC,
		pr:   prUC,
	}
}
