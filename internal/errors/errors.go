package errors

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Code   int      `json:"code"`
	Errors []string `json:"errors"`
}

func BadRequest(ctx *fiber.Ctx, err error) error {
	return toJSON(ctx, fiber.StatusBadRequest, err.Error())
}

func Unauthorized(ctx *fiber.Ctx, err error) error {
	return toJSON(ctx, fiber.StatusUnauthorized, err.Error())
}

func InternalError(ctx *fiber.Ctx, err error) error {
	return toJSON(ctx, fiber.StatusInternalServerError, err.Error())
}

func ValidationError(ctx *fiber.Ctx, errors []*ValidationRule) error {
	messages := make([]string, len(errors))

	for index, err := range errors {
		if err.Value != "" {
			err.Value = " " + err.Value
		}
		messages[index] = fmt.Sprintf("validation failed for property '%s', should be %s%s", err.Field, err.Tag, err.Value)
	}

	return toJSON(ctx, fiber.StatusBadRequest, messages...)
}

func toJSON(ctx *fiber.Ctx, code int, messages ...string) error {
	return ctx.Status(code).JSON(ErrorResponse{
		Code:   code,
		Errors: messages,
	})
}

var AccessDenied = errors.New("access denied")                      //nolint
var InvalidEmailOrPassord = errors.New("invalid email or password") //nolint
var ProjectNotFound = errors.New("project not found")               //nolint
var SessionNotFound = errors.New("session not found")               //nolint
var SpecNotFound = errors.New("spec not found")                     //nolint
var ApiKeyNotFound = errors.New("api key not found")                //nolint
var SessionFinished = errors.New("session already finished")        //nolint
