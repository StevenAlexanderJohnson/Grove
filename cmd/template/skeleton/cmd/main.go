package main

import (
	"template/internal/controllers"

	"github.com/StevenAlexanderJohnson/grove"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	*jwt.RegisteredClaims
}

func main() {
	authConfig, err := grove.LoadAuthenticatorConfigFromEnv()
	if err != nil {
		panic(err)
	}

	logger := grove.NewDefaultLogger("grove-test")
	authenticator := grove.NewAuthenticator[*Claims](authConfig)

	deps := grove.NewDependencies()
	deps.Set("logger", logger)

	api := grove.
		NewScope().
		UseMiddleware(grove.DefaultRequestLoggerMiddleware(logger)).
		UseMiddleware(grove.DefaultAuthMiddleware(
			authenticator,
			logger,
			func() *Claims { return &Claims{} },
		))

	app := grove.
		NewApp("skeleton-project").
		WithDependencies(deps).
		WithScope("/api", api).
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
