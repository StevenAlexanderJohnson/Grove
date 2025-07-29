package main

import (
	"net/http"
	"time"

	"github.com/StevenAlexanderJohnson/grove/v1"
	"github.com/golang-jwt/jwt/v5"
)

type HomeController struct {
	authenticator *grove.Authenticator[*CustomClaims]
	logger        grove.ILogger
}

func NewHomeController(authenticator *grove.Authenticator[*CustomClaims], logger grove.ILogger) *HomeController {
	return &HomeController{
		authenticator: authenticator,
		logger:        logger,
	}
}

func (c *HomeController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", c.handleHome)
	mux.HandleFunc("/login", c.sampleLogin)
}

func (c *HomeController) handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Home Page!"))
}

func (c *HomeController) sampleLogin(w http.ResponseWriter, r *http.Request) {
	token, err := c.authenticator.GenerateToken(&CustomClaims{
		UserID: "steven",
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    "example.com",
			Audience:  []string{"example.com"},
			Subject:   "steven",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        "unique-session-id",
		},
	})
	if err != nil {
		c.logger.Errorf("Failed to generate token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: token,
		Path:  "/",
	})
	// This is a placeholder for a login handler.
	// In a real application, you would handle user authentication here.
	w.Write([]byte("Login successful! A token would be generated and returned."))
}
