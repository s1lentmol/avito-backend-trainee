package http

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"github.com/silentmol/avito-backend-trainee/internal/apperr"
	prdto "github.com/silentmol/avito-backend-trainee/internal/pr/dto"
	userdto "github.com/silentmol/avito-backend-trainee/internal/user/dto"
)

func (h *Handle) SetIsActive(c *fiber.Ctx) error {
	req := &userdto.SetIsActiveRequest{}

	if err := c.BodyParser(req); err != nil {
		slog.Warn("SetIsActive: invalid request body", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		slog.Warn("SetIsActive: validation failed", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	resp, err := h.user.SetIsActive(c.Context(), req)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			slog.Info("SetIsActive: user not found",
				slog.String("user_id", req.UserID),
			)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "NOT_FOUND",
					"message": "user not found",
				},
			})
		}

		slog.Error("SetIsActive: failed to update is_active",
			slog.String("user_id", req.UserID),
			slog.Any("error", err),
		)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update is_active")
	}

	slog.Info("SetIsActive: user updated",
		slog.String("user_id", resp.User.ID),
		slog.Bool("is_active", resp.User.IsActive),
	)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": resp.User,
	})
}

func (h *Handle) GetReview(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		slog.Warn("GetReview: missing user_id")
		return fiber.NewError(fiber.StatusBadRequest, "user_id is required")
	}

	req := &prdto.GetReviewRequest{
		UserId: userID,
	}

	resp, err := h.pr.GetReview(c.Context(), req)
	if err != nil {
		slog.Error("GetReview: failed to get user reviews",
			slog.String("user_id", userID),
			slog.Any("error", err),
		)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get user reviews")
	}

	slog.Info("GetReview: reviews fetched",
		slog.String("user_id", resp.UserId),
		slog.Int("pull_requests_count", len(resp.PullRequests)),
	)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":       resp.UserId,
		"pull_requests": resp.PullRequests,
	})
}
