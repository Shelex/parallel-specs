package controllers

import (
	"github.com/Shelex/split-specs-v2/app"
	_ "github.com/Shelex/split-specs-v2/docs"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Swagger(app *app.App) {
	app.Router.Get("/swagger/*", fiberSwagger.WrapHandler)
}
