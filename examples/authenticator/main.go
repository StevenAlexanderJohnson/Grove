package main

import (
	"os"

	"github.com/StevenAlexanderJohnson/grove/v1"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	*jwt.RegisteredClaims
}

func main() {
	// This is just for demonstration purposes.
	// In a real application, you would set these environment variables
	// through your deployment configuration or secrets management.
	os.Setenv("JWT_PRIVATE_KEY_PATH", "jwt_private_key.pem")
	os.Setenv("JWT_ISSUER", "example.com")
	os.Setenv("JWT_AUDIENCE", "example.com")
	os.Setenv("JWT_SECRET", "super_secret_key")

	authConfig, err := grove.LoadAuthenticatorConfigFromEnv()
	if err != nil {
		panic(err)
	}
	if err := authConfig.Validate(); err != nil {
		panic("Invalid authenticator configuration")
	}
	authenticator := grove.NewAuthenticator[*CustomClaims](authConfig)

	logger := grove.NewDefaultLogger("authenticator")

	authScope := grove.
		NewScope().
		UseMiddleware(grove.DefaultAuthMiddleware(
			authenticator,
			logger,
			func() *CustomClaims { return &CustomClaims{} },
		)).
		AddController(&PrivateController{})

	if err := grove.
		NewApp("authenticator").
		WithScope("/private/", authScope).
		WithController(NewHomeController(authenticator, logger)).
		WithPort("8080").
		Run(); err != nil {
		panic(err)
	}
}
