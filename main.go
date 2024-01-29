package main

import (
	"log"
	"runtime/debug"

	"github.com/Shelex/parallel-specs/server"
)

// @title Parallel-Specs API
// @version 2.0
// @description service for distributing test files among processes/machines/containers
// @schemes https
// @host parallel-specs.shelex.dev
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

	server.Start()
}
