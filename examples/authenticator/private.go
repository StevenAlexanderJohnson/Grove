package main

import (
	"net/http"
	"time"

	"github.com/StevenAlexanderJohnson/grove"
)

type PrivateController struct{}

func (c *PrivateController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", c.handlePrivate)
	mux.HandleFunc("/logout", c.logout)
}
func (c *PrivateController) handlePrivate(w http.ResponseWriter, r *http.Request) {
	authToken, ok := r.Context().Value(grove.AuthTokenKey).(*CustomClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Write([]byte("Welcome to the Private Area! Your token is valid: " + authToken.Subject))
}

func (c *PrivateController) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-24 * time.Hour), // Expire the cookie
	})
	w.Write([]byte("You have been logged out."))
}
