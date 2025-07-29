package main

import (
	"net/http"
	"time"
)

type PublicController struct{}

func (c *PublicController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte("public content"))
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "example_token_value",
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour), // Set an appropriate expiration time
		})
		// Handle login logic here
		w.Write([]byte("login successful"))
	})
}
