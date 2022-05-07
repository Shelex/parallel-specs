package api

import (
	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/Shelex/split-specs-v2/internal/errors"
	"github.com/Shelex/split-specs-v2/internal/events"
	"github.com/Shelex/split-specs-v2/internal/projects"
	"github.com/Shelex/split-specs-v2/middleware"
	"github.com/gofiber/fiber/v2"
)

type ProjectsResponse struct {
	Projects []*entities.Project `json:"projects"`
}

// GetProjects godoc
// @Tags        project
// @Summary get projects for user
// @Accept  json
// @Success 200 {object} ProjectsResponse "projects"
// @Router /api/projects [get]
func (c *Controller) GetProjects(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)

	projects, err := projects.GetUserProjects(user.ID)
	if err != nil {
		return errors.BadRequest(ctx, err)
	}
	return ctx.JSON(ProjectsResponse{
		Projects: projects,
	})
}

type ProjectSessions struct {
	Sessions []entities.Session `json:"sessions"`
	Total    int                `json:"total"`
}

// GetProjectSessions godoc
// @Tags        project
// @Summary get project recorded sessions
// @Accept  json
// @Param  id path string true "project id" "uuid v4"
// @Param limit query integer false "pagination" 20
// @Param offset query integer false "pagination" 0
// @Success 200 {object} ProjectSessions "sessions"
// @Router /api/projects/{id}/sessions [get]
func (c *Controller) GetProjectSessions(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)
	projectID := ctx.Params("id")

	pagination := new(entities.Pagination)

	if err := ctx.QueryParser(pagination); err != nil {
		return errors.InternalError(ctx, err)
	}

	hasAccess, err := c.app.Repository.IsProjectAccessible(user.ID, projectID)
	if err != nil {
		return errors.BadRequest(ctx, err)
	}

	if !hasAccess {
		return errors.BadRequest(ctx, errors.ProjectNotFound)
	}

	sessions, total, err := c.app.Repository.GetProjectSessions(projectID, pagination)
	if err != nil {
		return errors.BadRequest(ctx, err)
	}
	return ctx.JSON(ProjectSessions{
		Sessions: sessions,
		Total:    total,
	})
}

// DeleteProject godoc
// @Tags        project
// @Summary delete project by id
// @Accept  json
// @Param  id path string true "project id" "uuid v4"
// @Success 200
// @Router /api/projects/{id} [delete]
func (c *Controller) DeleteProject(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)

	projectID := ctx.Params("id")

	if err := c.app.Repository.DeleteProject(user.ID, projectID); err != nil {
		return errors.BadRequest(ctx, err)
	}

	events.Handler.Publish(events.Project, events.ProjectEvent{
		Kind:   events.Deleted,
		ID:     projectID,
		Name:   "",
		UserID: user.ID,
	})

	return ctx.SendStatus(fiber.StatusOK)
}

// GetProjects godoc
// @Tags        project
// @Summary share project with another user
// @Accept  json
// @Param  id path string true "project id" "uuid v4"
// @Param  email path string true "other account email" "address@example.com"
// @Success 200
// @Router /api/projects/{id}/share/{email} [post]
func (c *Controller) ShareProject(ctx *fiber.Ctx) error {
	projectID := ctx.Params("id")
	email := ctx.Params("email")

	user, err := c.app.Repository.GetUserByEmail(email)
	if err != nil {
		return errors.BadRequest(ctx, err)
	}

	if err := c.app.Repository.AddUserProject(user.ID, projectID); err != nil {
		return errors.BadRequest(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}
