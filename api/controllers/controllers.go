package controllers

import (
	"github.com/Shelex/split-specs-v2/api/middleware"
	"github.com/Shelex/split-specs-v2/app"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Controller struct {
	app *app.App
}

func Register(app *app.App) {
	controller := Controller{app}
	app.Router.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// app.Router.Get("/dashboard", middleware.Monitor())

	app.Router.Use(middleware.Limiter())

	api := app.Router.Group("api")
	api.Post("/register", controller.Register)
	api.Post("/auth", controller.Login)
	api.Post("/new-password", controller.ChangePassword)

	api.Get("/listen", websocket.New(controller.Listener))

	api.Use(middleware.Auth())

	keys := api.Group("keys")

	keys.Get("/", controller.GetApiKeys)
	keys.Post("/", controller.AddApiKey)
	keys.Delete("/:id", controller.DeleteApiKey)

	projects := api.Group("projects")

	projects.Get("/", controller.GetProjects)
	projects.Get("/:id/sessions", controller.GetProjectSessions)
	projects.Post("/:id/share/:email", controller.ShareProject)
	projects.Delete("/:id", controller.DeleteProject)

	session := api.Group("session")

	session.Get("/:id", controller.GetSession)
	session.Get("/:id/next", controller.GetNextSpec)
	session.Post("/", controller.AddSession)
	session.Delete("/:id", controller.DeleteSession)

	spec := api.Group("spec")
	spec.Get("/:id", controller.GetSpecExecutions)
}
