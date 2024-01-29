package controllers

import (
	"github.com/Shelex/parallel-specs/api/middleware"
	"github.com/Shelex/parallel-specs/internal/entities"
	"github.com/Shelex/parallel-specs/internal/errors"
	"github.com/gofiber/fiber/v2"
)

type SpecResponse struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Executions []entities.Execution `json:"executions"`
	Project    entities.Project     `json:"project"`
}

// GetSpecExecutions godoc
// @Tags  spec
// @Summary get spec executions by id
// @Accept  json
// @Param Authorization header string true "Set Bearer token"
// @Param  id path string true "spec id" "uuid v4"
// @Param limit query integer false "pagination" 15
// @Success 200 {object} entities.Session "session"
// @Router /api/session/{id} [get]
func (c *Controller) GetSpecExecutions(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)

	specID := ctx.Params("id")

	spec, err := c.app.Repository.GetSpec(specID)
	if err != nil {
		return errors.BadRequest(ctx, err)
	}

	hasAccess, err := c.app.Repository.IsProjectAccessible(user.ID, spec.ProjectID)
	if err != nil {
		return errors.BadRequest(ctx, err)
	}

	if !hasAccess {
		return errors.BadRequest(ctx, errors.SessionNotFound)
	}

	project, err := c.app.Repository.GetProjectByID(spec.ProjectID)
	if err != nil {
		return errors.BadRequest(ctx, err)
	}

	history, err := c.app.Repository.GetExecutionHistory(specID, 15)
	if err != nil {
		return errors.BadRequest(ctx, err)
	}

	return ctx.JSON(SpecResponse{
		ID:         spec.ID,
		Name:       spec.FilePath,
		Project:    *project,
		Executions: history,
	})
}
