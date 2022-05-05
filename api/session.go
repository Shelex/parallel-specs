package api

import (
	"github.com/Shelex/split-specs-v2/internal/appError"
	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/Shelex/split-specs-v2/internal/events"
	"github.com/Shelex/split-specs-v2/internal/execution"
	"github.com/Shelex/split-specs-v2/internal/projects"
	"github.com/Shelex/split-specs-v2/middleware"
	"github.com/Shelex/split-specs-v2/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SessionInput struct {
	ProjectName string   `json:"projectName" validate:"required"`
	SpecFiles   []string `json:"specFiles" validate:"required"`
}

type AddSessionResponse struct {
	ProjectName string `json:"projectName"`
	ProjectID   string `json:"projectId"`
	SessionID   string `json:"sessionId"`
}

// AddSession godoc
// @Tags         session
// @Summary add new session
// @Accept  json
// @Param  input body SessionInput true "input" Example(SessionInput)
// @Success 200 {object} AddSessionResponse "session created"
// @Router /api/session [post]
func (c *Controller) AddSession(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)

	input := new(SessionInput)

	if err := ctx.BodyParser(&input); err != nil {
		return FailedToParseRequestBody(ctx, err.Error())
	}

	errors := ValidateStruct(*input)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	project, isNew, err := projects.GetByNameOrCreateNew(user.ID, input.ProjectName)
	if err != nil {
		return SendError(ctx, fiber.StatusBadRequest, err)
	}

	if isNew {
		events.Handler.Publish(events.Project, events.ProjectEvent{
			Kind:   events.Created,
			ID:     project.ID,
			Name:   project.Name,
			UserID: user.ID,
		})
	}

	specs, err := repository.DB.AddSpecsMaybe(project.ID, input.SpecFiles)
	if err != nil {
		return SendError(ctx, fiber.StatusBadRequest, err)
	}

	executions, err := execution.SpecsToExecutions(specs)
	if err != nil {
		return SendError(ctx, fiber.StatusBadRequest, err)
	}

	sessionExecution := entities.Session{
		ID:        uuid.NewString(),
		ProjectID: project.ID,
		CreatedAt: repository.GetTimestamp(),
	}

	if err := repository.DB.AddSession(sessionExecution); err != nil {
		return SendError(ctx, fiber.StatusBadRequest, err)
	}

	events.Handler.Publish(events.Session, events.SessionEvent{
		Kind:      events.Created,
		ID:        sessionExecution.ID,
		ProjectID: sessionExecution.ProjectID,
	})

	if err := repository.DB.AddExecutions(sessionExecution.ID, executions); err != nil {
		return SendError(ctx, fiber.StatusBadRequest, err)
	}

	return ctx.JSON(AddSessionResponse{
		ProjectName: input.ProjectName,
		ProjectID:   sessionExecution.ProjectID,
		SessionID:   sessionExecution.ID,
	})
}

// GetSession godoc
// @Tags  session
// @Summary get session with executions by id
// @Accept  json
// @Param  id path string true "session id" "uuid v4"
// @Success 200 {object} entities.Session "session"
// @Router /api/session/{id} [get]
func (c *Controller) GetSession(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)
	ID := ctx.Params("id")
	session, err := repository.DB.GetSessionWithExecution(ID)
	if err != nil {
		return SendError(ctx, fiber.StatusBadRequest, err)
	}

	hasAccess, err := repository.DB.IsProjectAccessible(user.ID, session.ProjectID)
	if err != nil {
		return SendError(ctx, fiber.StatusBadRequest, err)
	}

	if !hasAccess {

		return SendError(ctx, fiber.StatusBadRequest, appError.SessionNotFound)
	}

	return ctx.JSON(session)
}

// DeleteProject godoc
// @Tags        session
// @Summary delete session by id
// @Accept  json
// @Param  id path string true "session id" "uuid v4"
// @Success 200
// @Router /api/session/{id} [delete]
func (c *Controller) DeleteSession(ctx *fiber.Ctx) error {
	ID := ctx.Params("id")

	if err := repository.DB.DeleteSession(ID); err != nil {
		return SendError(ctx, fiber.StatusBadRequest, err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

type NextSpecResponse struct {
	Next string `json:"next"`
}

// GetNextSpec godoc
// @Tags  session
// @Summary get next spec file to be executed
// @Accept  json
// @Param  id path string true "session id" "uuid v4"
// @Param machineId query string false "specify machine id" "default"
// @Param previousStatus query string false "specify status of previous spec execution" "unknown"
// @Success 200 {object} NextSpecResponse "next"
// @Router /api/session/{id}/next [get]
func (c *Controller) GetNextSpec(ctx *fiber.Ctx) error {
	type NextOptions struct {
		MachineID      string `query:"machineId"`
		PreviousStatus string `query:"previousStatus"`
	}

	opts := new(NextOptions)

	sessionID := ctx.Params("id")

	if err := ctx.QueryParser(opts); err != nil {
		return FailedToParseRequestBody(ctx, err.Error())
	}

	machineID := "default"
	if opts != nil && opts.MachineID != "" {
		machineID = opts.MachineID
	}

	previousStatus := "unknown"
	if opts != nil && opts.PreviousStatus != "" {
		previousStatus = opts.PreviousStatus
	}

	next, err := execution.Next(sessionID, machineID, previousStatus)
	if err != nil {
		return SendError(ctx, fiber.StatusBadRequest, err)
	}

	return ctx.JSON(NextSpecResponse{
		Next: next,
	})
}
