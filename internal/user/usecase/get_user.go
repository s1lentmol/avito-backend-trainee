package usecase

import (
	"context"
	"fmt"

	"github.com/silentmol/avito-backend-trainee/internal/user/dto"
)

func (u *UserUsecase) GetUser(ctx context.Context,
	getUserRequest *dto.GetUserRequest) (*dto.GetUserResponse, error) {
	userID := getUserRequest.UserID

	user, err := u.userProvider.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get profile from postgres: %w", err)
	}

	return &dto.GetUserResponse{
		User: *user,
	}, nil
}
