package main

import (
	"net/http"
	"time"
)

type PrivateController struct{}

func (c *PrivateController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("private content"))
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Path:    "/",
			Expires: time.Now().Add(-24 * time.Hour), // Expire the cookie
		})
		w.Write([]byte("logout successful"))
	})
}
