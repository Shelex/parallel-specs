package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/Shelex/split-specs-v2/api"
	"github.com/Shelex/split-specs-v2/app"
	"github.com/Shelex/split-specs-v2/env"
)

// @title Split specs API
// @version 2.0
// @description service for distributing test files among processes/machines/containers
// @schemes http
// @host localhost:3000
// @BasePath /
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email ovr.shevtsov@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("%s: %s", r, string(debug.Stack()))
		}
	}()

	config := env.ReadEnv()

	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)

	app, err := app.NewApp(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer app.Repository.ShutDown(ctx)

	api.RegisterControllers(app)
	api.RegisterSwagger(app)

	go func() {
		if err := app.Router.Listen(":" + config.HttpPort); err != nil {
			log.Printf("Could not start HTTP server %s:\n", err)
		}
	}()

	log.Println("Starting HTTP server")
	<-ctx.Done()
	log.Println("Stopping HTTP server")
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Println("Got Interrupt signal")
		stop()
	}()
}
