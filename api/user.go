package api

import (
	"github.com/Shelex/split-specs-v2/internal/appError"
	"github.com/Shelex/split-specs-v2/internal/jwt"
	"github.com/Shelex/split-specs-v2/internal/users"
	"github.com/Shelex/split-specs-v2/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type tokenResponse struct {
	Token string `json:"token"`
}

// Register godoc
// @Tags account
// @Summary register new account
// @Accept  json
// @Param  user body users.User true "user" Example(users.User)
// @Success 200 {object} tokenResponse "token response"
// @Router /api/register [post]
func (c *Controller) Register(ctx *fiber.Ctx) error {
	var user *users.User

	if err := ctx.BodyParser(&user); err != nil {
		return FailedToParseRequestBody(ctx, err.Error())
	}

	errors := ValidateStruct(*user)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	if user.Email == "" || user.Password == "" {
		return appError.InvalidEmailOrPassord
	}

	user.ID = uuid.NewString()

	if user.Exist() {
		return appError.InvalidEmailOrPassord
	}

	if err := user.Create(); err != nil {
		return err
	}

	token, err := jwt.GenerateToken(*user)
	if err != nil {
		return err
	}

	return ctx.JSON(tokenResponse{Token: token})
}

// Login godoc
// @Tags account
// @Summary get authorization token
// @Accept  json
// @Param  user body users.User true "user" Example(users.User)
// @Success 200 {object} tokenResponse "token response"
// @Router /api/auth [post]
func (c *Controller) Login(ctx *fiber.Ctx) error {
	var user *users.User

	if err := ctx.BodyParser(&user); err != nil {
		return FailedToParseRequestBody(ctx, err.Error())
	}

	dbUser, err := user.Authenticate()
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	user.ID = dbUser.ID

	token, err := jwt.GenerateToken(*user)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	return ctx.JSON(tokenResponse{Token: token})
}

type PasswordChange struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required"`
}

// ChangePassword godoc
// @Tags account
// @Summary change password for the account
// @Accept  json
// @Param  input body PasswordChange true "input" Example(PasswordChange)
// @Success 200
// @Router /api/new-password [post]
func (c *Controller) ChangePassword(ctx *fiber.Ctx) error {
	change := new(PasswordChange)

	if err := ctx.BodyParser(&change); err != nil {
		return FailedToParseRequestBody(ctx, err.Error())
	}

	errors := ValidateStruct(*change)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	user := middleware.GetUser(ctx)

	if err := user.ChangePassword(change.CurrentPassword, change.NewPassword); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)

	}

	return ctx.SendStatus(fiber.StatusOK)
}
