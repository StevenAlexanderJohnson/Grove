package main

import (
	"fmt"
	"net/http"

	"github.com/StevenAlexanderJohnson/grove"
)

type HomeParameters struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *HomeParameters) Validate() error {
	if h.Title == "" {
		return fmt.Errorf("Title is required")
	}
	if h.Description == "" {
		return fmt.Errorf("Description is required")
	}
	return nil
}

type HomeController struct{}

func (h *HomeController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /", h.Index)
}

func (h *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	parameters, err := grove.ParseJsonBodyFromRequest[HomeParameters](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := parameters.Validate(); err != nil {
		grove.WriteErrorToResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := grove.WriteJsonBodyToResponse(w, parameters); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
