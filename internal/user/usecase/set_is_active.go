package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/silentmol/avito-backend-trainee/internal/user/dto"
)

func (u *UserUsecase) SetIsActive(ctx context.Context,
	setIsActiveRequest *dto.SetIsActiveRequest) (*dto.SetIsActiveResponse, error) {

	userID := setIsActiveRequest.UserID
	isActive := setIsActiveRequest.IsActive

	updatedUser, err := u.userProvider.SetIsActive(ctx, userID, isActive)

	if err != nil {
		slog.Error("UserUsecase.SetIsActive: provider error",
			slog.String("user_id", userID),
			slog.Bool("is_active", isActive),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("update isActive in postgres: %w", err)
	}

	slog.Info("UserUsecase.SetIsActive: user updated",
		slog.String("user_id", updatedUser.ID),
		slog.Bool("is_active", updatedUser.IsActive),
	)

	return &dto.SetIsActiveResponse{
		User: *updatedUser,
	}, nil
}
