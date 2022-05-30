package app

import (
	"github.com/Shelex/split-specs-v2/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProvideRouter() *fiber.App {
	router := fiber.New()
	router.Use(middleware.Logger(), middleware.Cors())
	return router
}
