package app

import (
	"github.com/Shelex/parallel-specs/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProvideRouter() *fiber.App {
	router := fiber.New()
	router.Use(middleware.Logger(), middleware.Cors())
	return router
}
