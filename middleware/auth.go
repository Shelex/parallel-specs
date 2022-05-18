package middleware

import (
	"fmt"

	"github.com/Shelex/split-specs-v2/internal/errors"
	keys "github.com/Shelex/split-specs-v2/internal/jwt"
	"github.com/Shelex/split-specs-v2/internal/users"
	"github.com/Shelex/split-specs-v2/repository"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func Auth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningMethod: "RS256",
		SigningKey:    keys.SignKey.Public(),
		ContextKey:    "token",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return errors.Unauthorized(ctx, err)
		},
		SuccessHandler: func(ctx *fiber.Ctx) error {
			token, ok := ctx.Locals("token").(*jwt.Token)
			if !ok {
				return errors.Unauthorized(ctx, errors.AccessDenied)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return errors.Unauthorized(ctx, errors.AccessDenied)
			}
			user := users.User{
				Email: claims["email"].(string),
				ID:    claims["id"].(string),
			}

			entity := claims["entity"].(string)

			if entity != "user" {
				if err := checkApiKey(user, entity); err != nil {
					return errors.Unauthorized(ctx, err)
				}

			} else {
				if err := checkUser(user); err != nil {
					return errors.Unauthorized(ctx, errors.AccessDenied)
				}
			}

			ctx.Locals("user", user)
			return ctx.Next()
		},
	})
}

func GetUser(ctx *fiber.Ctx) users.User {
	return ctx.Locals("user").(users.User)
}

func checkUser(user users.User) error {
	if _, err := repository.DB.GetUserByEmail(user.Email); err != nil {
		return err
	}
	return nil
}

func checkApiKey(user users.User, entity string) error {
	isValid, err := repository.DB.IsApiKeyValid(user.ID, entity)

	if err != nil {
		return fmt.Errorf("failed to validate api key")
	}

	if !isValid {
		return fmt.Errorf("api key is invalid")
	}
	return nil
}
