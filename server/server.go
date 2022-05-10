package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shelex/split-specs-v2/api"
	"github.com/Shelex/split-specs-v2/app"
	"github.com/Shelex/split-specs-v2/env"
)

func Start() {
	config := env.ReadEnv()

	ctx, cancel := context.WithCancel(context.Background())

	withGracefulShutdown(cancel)
	app, err := app.NewApp(ctx, config)
	if err != nil {
		log.Fatalf("Could not start HTTP server %s:\n", err)
	}
	defer app.Repository.ShutDown(ctx)

	api.RegisterControllers(app)
	api.RegisterSwagger(app)

	go func() {
		if err := app.Router.Listen(":" + config.HttpPort); err != nil {
			log.Fatalf("Could not start HTTP server %s:\n", err)
		}
	}()

	log.Println("Starting HTTP server")
	<-ctx.Done()
	log.Println("Stopping HTTP server")
}

func withGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Println("Got Interrupt signal")
		stop()
	}()
}
