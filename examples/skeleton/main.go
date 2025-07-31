package main

import (
	"net/http"

	"github.com/StevenAlexanderJohnson/grove"
)

func main() {
	scope := grove.NewScope().WithRoute("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Grove!"))
	}))

	app := grove.NewApp("Skeleton App").WithScope("/", scope)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
