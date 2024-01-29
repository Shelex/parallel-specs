package app

import (
	"context"

	"github.com/Shelex/parallel-specs/env"
	"github.com/Shelex/parallel-specs/internal/events"
	"github.com/Shelex/parallel-specs/repository"
	"github.com/Shelex/parallel-specs/repository/mock"
	"github.com/Shelex/parallel-specs/repository/postgres"
	"github.com/gofiber/fiber/v2"
)

type App struct {
	Router     *fiber.App
	Repository repository.Storage
	Events     *events.Orchestrator
}

func NewApp(ctx context.Context, config *env.Config) (*App, error) {
	router := ProvideRouter()

	var repo repository.Storage

	if config.Env == "dev" {
		mocked, err := mock.NewMockStorage(ctx)
		if err != nil {
			return nil, err
		}
		repo = mocked
	} else {
		postgres, err := postgres.NewPostgresStorage(ctx, config.DbConnectionUrl)
		if err != nil {
			return nil, err
		}
		repo = postgres
	}

	events := events.Start()

	app := &App{
		Router:     router,
		Repository: repo,
		Events:     events,
	}

	return app, nil
}
