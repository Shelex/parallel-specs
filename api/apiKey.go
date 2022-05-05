package api

import (
	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/Shelex/split-specs-v2/internal/jwt"
	"github.com/Shelex/split-specs-v2/middleware"
	"github.com/Shelex/split-specs-v2/repository"
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
// @Param  input body ApiKeyInput true "input" Example(ApiKeyInput)
// @Success 200 {object} tokenResponse "api token"
// @Router /api/keys [post]
func (c *Controller) AddApiKey(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)

	input := new(ApiKeyInput)

	if err := ctx.BodyParser(&input); err != nil {
		return FailedToParseRequestBody(ctx, err.Error())
	}

	errors := ValidateStruct(*input)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	id := uuid.NewString()

	apiKey := entities.ApiKey{
		ID:       id,
		UserID:   user.ID,
		Name:     input.Name,
		ExpireAt: input.ExpireAt,
	}

	if err := repository.DB.AddApiKey(user.ID, apiKey); err != nil {
		return SendError(ctx, fiber.StatusBadRequest, err)
	}

	token, err := jwt.GenerateApiKey(user, apiKey)
	if err != nil {
		return SendError(ctx, fiber.StatusBadRequest, err)
	}

	return ctx.JSON(fiber.Map{
		"token": token,
	})
}

// GetApiKeys godoc
// @Tags api key
// @Summary get user api keys
// @Accept  json
// @Success 200 {object} []entities.ApiKey "api keys"
// @Router /api/keys [get]
func (c *Controller) GetApiKeys(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)

	keys, err := repository.DB.GetApiKeys(user.ID)
	if err != nil {
		return SendError(ctx, fiber.StatusBadRequest, err)
	}

	return ctx.JSON(keys)
}

// DeleteApiKey godoc
// @Tags api key
// @Summary delete api key
// @Accept  json
// @Param  id path string true "api key id" "uuid v4"
// @Success 200
// @Router /api/keys/{id} [delete]
func (c *Controller) DeleteApiKey(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)

	id := ctx.Params("id")

	if err := repository.DB.DeleteApiKey(user.ID, id); err != nil {
		return SendError(ctx, fiber.StatusBadRequest, err)
	}
	return nil
}
