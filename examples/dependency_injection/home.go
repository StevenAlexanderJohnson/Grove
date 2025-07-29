package main

import (
	"net/http"

	"github.com/StevenAlexanderJohnson/grove"
)

type HomeController struct {
	names map[string]string
}

// RegisterRoutes implements grove.IController.
func (c *HomeController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", c.Index)
}

func (c *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, " + c.names["first"] + " " + c.names["last"] + "!"))
}

func NewHomeController(dependencies *grove.Dependencies) *HomeController {
	names := grove.DependencyMustGet[map[string]string](dependencies, "names")
	return &HomeController{
		names: names,
	}
}
