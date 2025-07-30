package controllers

import (
	"net/http"
)

type HomeController struct{}

func NewHomeController() *HomeController {
	return &HomeController{}
}

func (c *HomeController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", c.Index)
}

func (c *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("<h1>Welcome to the Home Page</h1>"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
