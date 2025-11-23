package http

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"github.com/silentmol/avito-backend-trainee/internal/apperr"
	teamdto "github.com/silentmol/avito-backend-trainee/internal/team/dto"
)

func (h *Handle) AddTeam(c *fiber.Ctx) error {
	req := &teamdto.AddTeamRequest{}

	if err := c.BodyParser(req); err != nil {
		slog.Warn("AddTeam: invalid request body", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		slog.Warn("AddTeam: validation failed", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	seen := make(map[string]struct{})
	for _, m := range req.Members {
		if _, ok := seen[m.ID]; ok {
			slog.Warn("AddTeam: duplicate user_id in members",
				slog.String("team_name", req.Name),
				slog.String("user_id", m.ID),
			)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "DUPLICATE_USER_ID",
					"message": "user_id must be unique within team members",
				},
			})
		}
		seen[m.ID] = struct{}{}
	}

	resp, err := h.team.CreateTeam(c.Context(), req)
	if err != nil {
		if errors.Is(err, apperr.ErrTeamExists) {
			slog.Info("AddTeam: team already exists",
				slog.String("team_name", req.Name),
			)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "TEAM_EXISTS",
					"message": "team_name already exists",
				},
			})
		}

		slog.Error("AddTeam: failed to add team",
			slog.String("team_name", req.Name),
			slog.Any("error", err),
		)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to add team")
	}

	slog.Info("AddTeam: team created",
		slog.String("team_name", resp.Team.Name),
		slog.Int("members_count", len(resp.Team.Members)),
	)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"team": resp.Team,
	})
}

func (h *Handle) GetTeam(c *fiber.Ctx) error {
	req := &teamdto.GetTeamRequest{
		TeamName: c.Query("team_name"),
	}

	if req.TeamName == "" {
		slog.Warn("GetTeam: missing team_name")
		return fiber.NewError(fiber.StatusBadRequest, "team_name is required")
	}

	resp, err := h.team.GetTeam(c.Context(), req)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			slog.Info("GetTeam: team not found",
				slog.String("team_name", req.TeamName),
			)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "NOT_FOUND",
					"message": "team not found",
				},
			})
		}

		slog.Error("GetTeam: failed to get team",
			slog.String("team_name", req.TeamName),
			slog.Any("error", err),
		)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get team")
	}

	slog.Info("GetTeam: team fetched",
		slog.String("team_name", resp.Team.Name),
		slog.Int("members_count", len(resp.Team.Members)),
	)

	return c.Status(fiber.StatusOK).JSON(resp.Team)
}
