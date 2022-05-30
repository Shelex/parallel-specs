package controllers

import (
	"github.com/Shelex/split-specs-v2/api/middleware"
	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/Shelex/split-specs-v2/internal/errors"
	"github.com/Shelex/split-specs-v2/internal/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ApiKeyInput struct {
	Name     string `json:"name" validate:"required"`
	ExpireAt uint64 `json:"expireAt" validate:"required"`
}

// AddApiKey godoc
// @Tags api key
// @Summary add new api key
// @Accept  json
// @Param Authorization header string true "Set Bearer token"
// @Param  input body ApiKeyInput true "input" Example(ApiKeyInput)
// @Success 200 {object} tokenResponse "api token"
// @Router /api/keys [post]
func (c *Controller) AddApiKey(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)

	input := new(ApiKeyInput)

	if err := ctx.BodyParser(&input); err != nil {
		return errors.InternalError(ctx, err)
	}

	failed := errors.ValidateStruct(*input)
	if failed != nil {
		return errors.ValidationError(ctx, failed)
	}

	id := uuid.NewString()

	apiKey := entities.ApiKey{
		ID:       id,
		UserID:   user.ID,
		Name:     input.Name,
		ExpireAt: input.ExpireAt,
	}

	if err := c.app.Repository.AddApiKey(user.ID, apiKey); err != nil {
		return errors.BadRequest(ctx, err)
	}

	token, err := jwt.GenerateApiKey(user, apiKey)
	if err != nil {
		return errors.BadRequest(ctx, err)
	}

	return ctx.JSON(tokenResponse{
		Token: token,
	})
}

// GetApiKeys godoc
// @Tags api key
// @Summary get user api keys
// @Accept  json
// @Param Authorization header string true "Set Bearer token"
// @Success 200 {object} []entities.ApiKey "api keys"
// @Router /api/keys [get]
func (c *Controller) GetApiKeys(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)

	keys, err := c.app.Repository.GetApiKeys(user.ID)
	if err != nil {
		return errors.BadRequest(ctx, err)
	}

	return ctx.JSON(keys)
}

// DeleteApiKey godoc
// @Tags api key
// @Summary delete api key
// @Accept  json
// @Param Authorization header string true "Set Bearer token"
// @Param  id path string true "api key id" "uuid v4"
// @Success 200
// @Router /api/keys/{id} [delete]
func (c *Controller) DeleteApiKey(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)

	id := ctx.Params("id")

	if err := c.app.Repository.DeleteApiKey(user.ID, id); err != nil {
		return errors.BadRequest(ctx, err)
	}
	return nil
}
