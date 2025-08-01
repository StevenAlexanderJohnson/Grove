package main

import "github.com/StevenAlexanderJohnson/grove"

func main() {
	app := grove.NewApp("Parsing Body Example").
		WithController(&HomeController{})

	if err := app.Run(); err != nil {
		panic(err)
	}
}
