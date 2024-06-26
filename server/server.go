package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shelex/parallel-specs/api/controllers"
	"github.com/Shelex/parallel-specs/app"
	"github.com/Shelex/parallel-specs/env"
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

	controllers.Register(app)
	controllers.Swagger(app)

	go func() {
		if err := app.Router.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port)); err != nil {
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
