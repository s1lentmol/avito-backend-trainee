package http

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"github.com/silentmol/avito-backend-trainee/internal/apperr"
	prdto "github.com/silentmol/avito-backend-trainee/internal/pr/dto"
)

func (h *Handle) CreatePR(c *fiber.Ctx) error {
	req := &prdto.CreatePRRequest{}

	if err := c.BodyParser(req); err != nil {
		slog.Warn("CreatePR: invalid request body", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		slog.Warn("CreatePR: validation failed", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	resp, err := h.pr.CreatePR(c.Context(), req)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			slog.Info("CreatePR: author or team not found",
				slog.String("pr_id", req.PrID),
				slog.String("author_id", req.AuthorId),
			)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "NOT_FOUND",
					"message": "author or team not found",
				},
			})
		}

		if errors.Is(err, apperr.ErrPRExists) {
			slog.Info("CreatePR: PR already exists",
				slog.String("pr_id", req.PrID),
			)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "PR_EXISTS",
					"message": "PR id already exists",
				},
			})
		}

		slog.Error("CreatePR: failed to create pull request",
			slog.String("pr_id", req.PrID),
			slog.String("author_id", req.AuthorId),
			slog.Any("error", err),
		)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create pull request")
	}

	slog.Info("CreatePR: pull request created",
		slog.String("pr_id", resp.PullRequest.ID),
		slog.String("author_id", resp.PullRequest.AuthorId),
		slog.String("status", string(resp.PullRequest.Status)),
	)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"pr": resp.PullRequest,
	})
}

func (h *Handle) MergePR(c *fiber.Ctx) error {
	req := &prdto.MergePRRequest{}

	if err := c.BodyParser(req); err != nil {
		slog.Warn("MergePR: invalid request body", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	resp, err := h.pr.MergePR(c.Context(), req)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			slog.Info("MergePR: pull request not found",
				slog.String("pr_id", req.PrID),
			)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "NOT_FOUND",
					"message": "pull request not found",
				},
			})
		}

		slog.Error("MergePR: failed to merge pull request",
			slog.String("pr_id", req.PrID),
			slog.Any("error", err),
		)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to merge pull request")
	}

	slog.Info("MergePR: pull request merged",
		slog.String("pr_id", resp.PullRequest.ID),
		slog.String("status", string(resp.PullRequest.Status)),
	)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"pr": resp.PullRequest,
	})
}

func (h *Handle) ReassignPR(c *fiber.Ctx) error {
	req := &prdto.ReassignPRRequest{}

	if err := c.BodyParser(req); err != nil {
		slog.Warn("ReassignPR: invalid request body", slog.Any("error", err))
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	resp, err := h.pr.ReassignPR(c.Context(), req)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			slog.Info("ReassignPR: pull request or user not found",
				slog.String("pr_id", req.PrID),
				slog.String("old_reviewer_id", req.OldReviewerId),
			)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "NOT_FOUND",
					"message": "pull request or user not found",
				},
			})
		}

		if errors.Is(err, apperr.ErrPRMerged) {
			slog.Info("ReassignPR: PR already merged",
				slog.String("pr_id", req.PrID),
			)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "PR_MERGED",
					"message": "cannot reassign on merged PR",
				},
			})
		}

		if errors.Is(err, apperr.ErrNotAssigned) {
			slog.Info("ReassignPR: reviewer not assigned to PR",
				slog.String("pr_id", req.PrID),
				slog.String("old_reviewer_id", req.OldReviewerId),
			)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "NOT_ASSIGNED",
					"message": "reviewer is not assigned to this PR",
				},
			})
		}

		if errors.Is(err, apperr.ErrNoCandidate) {
			slog.Info("ReassignPR: no candidate found in team",
				slog.String("pr_id", req.PrID),
				slog.String("old_reviewer_id", req.OldReviewerId),
			)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "NO_CANDIDATE",
					"message": "no active replacement candidate in team",
				},
			})
		}

		slog.Error("ReassignPR: cannot reassign reviewer",
			slog.String("pr_id", req.PrID),
			slog.String("old_reviewer_id", req.OldReviewerId),
			slog.Any("error", err),
		)
		return fiber.NewError(fiber.StatusConflict, "cannot reassign reviewer")
	}

	slog.Info("ReassignPR: reviewer reassigned",
		slog.String("pr_id", resp.PullRequest.ID),
		slog.String("replaced_by", resp.ReplacedBy),
	)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"pr":          resp.PullRequest,
		"replaced_by": resp.ReplacedBy,
	})
}
