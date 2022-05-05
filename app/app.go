package app

import (
	"context"

	"github.com/Shelex/split-specs-v2/env"
	"github.com/Shelex/split-specs-v2/internal/events"
	"github.com/Shelex/split-specs-v2/repository"
	"github.com/Shelex/split-specs-v2/repository/mock"
	"github.com/Shelex/split-specs-v2/repository/postgres"
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

	return &App{
		Router:     router,
		Repository: repo,
		Events:     events,
	}, nil
}
