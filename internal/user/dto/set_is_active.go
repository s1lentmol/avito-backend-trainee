package dto

import (
	"github.com/silentmol/avito-backend-trainee/internal/user/domain"
)

type SetIsActiveRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveResponse struct {
	domain.User
}
