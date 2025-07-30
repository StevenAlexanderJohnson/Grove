package main

import (
	"template/internal/controllers"

	"github.com/StevenAlexanderJohnson/grove"
)

func main() {
	app := grove.
		NewApp().
		WithController(controllers.NewHomeController())
	if err := app.Run(); err != nil {
		panic(err)
	}
}
