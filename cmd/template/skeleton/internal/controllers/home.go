package controllers

import (
	"net/http"

	"github.com/StevenAlexanderJohnson/grove"
)

type HomeController struct {
	logger grove.ILogger
}

func NewHomeController(logger grove.ILogger) *HomeController {
	return &HomeController{logger: logger}
}

func (c *HomeController) RegisterRoutes(mux *http.ServeMux) {
	c.logger.Info("Registering HomeController routes")
	mux.HandleFunc("/", c.Index)
}

func (c *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Handling request for home page")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("<h1>Welcome to the Home Page</h1>"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
