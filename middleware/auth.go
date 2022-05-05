package middleware

import (
	"fmt"

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
		SuccessHandler: func(ctx *fiber.Ctx) error {
			token := ctx.Locals("token").(*jwt.Token)

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

				user := users.User{
					Email: claims["email"].(string),
					ID:    claims["id"].(string),
				}

				entity := claims["entity"].(string)

				if entity != "user" {
					isValid, err := repository.DB.IsApiKeyValid(user.ID, entity)

					if err != nil {
						return fmt.Errorf("failed to validate api key")
					}

					if !isValid {
						return fmt.Errorf("api key is invalid")
					}
				}

				ctx.Locals("user", user)
				return ctx.Next()

			}
			return fmt.Errorf("could not parse claims from jwt token")
		},
	})
}

func GetUser(ctx *fiber.Ctx) users.User {
	return ctx.Locals("user").(users.User)
}
