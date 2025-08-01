package main

import (
	"context"
	"net/http"

	"github.com/StevenAlexanderJohnson/grove"
)

func main() {
	logger := grove.NewDefaultLogger("scoped_routes")

	authScope := grove.
		NewScope().
		// Ad-hoc middleware dummy to simulate session token validation
		// In a real application, this would be replaced with actual session validation logic
		WithMiddleware(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sessionCookie, err := r.Cookie("session_token")
				if err != nil || sessionCookie == nil || sessionCookie.Value == "" {
					logger.Warning("Unauthorized request: no session token found")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				token := sessionCookie.Value
				// Validate the token using the authenticator
				tokenContext := context.WithValue(r.Context(), grove.AuthTokenKey, token)
				// If token is valid, proceed to the next handler
				next.ServeHTTP(w, r.WithContext(tokenContext))
			})
		}).
		WithController(&PrivateController{})

	publicScope := grove.
		NewScope().
		WithController(&PublicController{})

	app := grove.
		NewApp("scoped_routes").
		WithPort("8080").
		WithMiddleware(grove.DefaultRequestLoggerMiddleware(logger)).
		WithScope("/", publicScope).
		WithScope("/private/", authScope)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
