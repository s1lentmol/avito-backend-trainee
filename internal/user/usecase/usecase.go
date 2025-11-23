package usecase

import (
	"context"

	"github.com/silentmol/avito-backend-trainee/internal/user/domain"
)

type UserProvider interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
	SetIsActive(ctx context.Context, id string, isActive bool) (*domain.User, error)
}

type UserUsecase struct {
	userProvider UserProvider
}

func NewUserUsecase(repo UserProvider) *UserUsecase {
	return &UserUsecase{
		userProvider: repo,
	}
}
