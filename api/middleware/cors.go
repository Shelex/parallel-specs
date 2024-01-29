package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, https://parallel-specs.shelex.dev",
		AllowHeaders: "Access-Control-Allow-Origin, Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, DELETE, OPTIONS",
	})
}
