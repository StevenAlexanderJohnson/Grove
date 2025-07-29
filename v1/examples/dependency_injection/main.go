package main

import (
	"github.com/StevenAlexanderJohnson/grove/v1"
)

func main() {
	namesKey := grove.DependencyKey("names")
	names := map[string]string{
		"first": "Steven",
		"last":  "Johnson",
	}

	deps := grove.NewDependencies()
	deps.Set(namesKey, names)

	app := grove.
		NewApp().
		WithPort("8080").
		WithDependencies(deps).
		WithControllerFactory(func(dependencies *grove.Dependencies) grove.IController {
			return NewHomeController(dependencies)
		})

	if err := app.Run(); err != nil {
		panic(err)
	}
}
