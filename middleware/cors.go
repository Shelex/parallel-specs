package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Access-Control-Allow-Origin, Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, DELETE, OPTIONS",
	})
}
