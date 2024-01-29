package controllers

import (
	"github.com/Shelex/parallel-specs/app"
	_ "github.com/Shelex/parallel-specs/docs"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Swagger(app *app.App) {
	app.Router.Get("/swagger/*", fiberSwagger.WrapHandler)
}
