package dto

import (
	"github.com/silentmol/avito-backend-trainee/internal/user/domain"
)

type GetUserRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

type GetUserResponse struct {
	domain.User
}
