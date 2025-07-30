package main

import (
	"template/internal/controllers"

	"github.com/StevenAlexanderJohnson/grove/v1"
)

func main() {
	logger := grove.NewDefaultLogger("grove-test")

	deps := grove.NewDependencies()
	deps.Set("logger", logger)

	app := grove.
		NewApp().
		WithDependencies(deps).
		WithControllerFactory(
			func(dependencies *grove.Dependencies) grove.IController {
				return controllers.NewHomeController(
					grove.DependencyMustGet[grove.ILogger](dependencies, "logger"),
				)
			},
		).
		WithMiddleware(grove.DefaultRequestLoggerMiddleware(logger))
	if err := app.Run(); err != nil {
		panic(err)
	}
}
