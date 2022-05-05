package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func SendError(args ...interface{}) error {
	ctx := args[0].(*fiber.Ctx)
	code := args[1].(int)
	err := args[2].(error)
	return ctx.Status(code).SendString(err.Error())
}

func FailedToParseRequestBody(ctx *fiber.Ctx, cause string) error {
	err := fmt.Errorf("failed to parse request body: %s", cause)
	return SendError(ctx, fiber.StatusInternalServerError, err)
}
