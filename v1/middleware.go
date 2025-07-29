package grove

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Middleware defines a function type that takes an http.Handler and returns an http.Handler.
// This allows for chaining middleware functions in the Grove application.
type Middleware func(next http.Handler) http.Handler

type requestIDKeyType struct{}

var RequestIDKey = requestIDKeyType{}

type authTokenKeyType struct{}

var AuthTokenKey = authTokenKeyType{}

// DefaultAuthMiddleware is a middleware that provides default authentication logic.
// It checks for a valid token in the request header and denies access if the token is missing
func DefaultAuthMiddleware[T jwt.Claims](authenticator *Authenticator[T], logger ILogger, claimsFactory func() T) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Implement default authentication logic here
			// For example, check for a valid token in the request header
			token := r.Header.Get("Authorization")
			if token == "" {
				sessionCookie, err := r.Cookie("session_token")
				if err != nil || sessionCookie == nil || sessionCookie.Value == "" {
					// If no token is found, return unauthorized
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				token = sessionCookie.Value
			}
			claims := claimsFactory()
			// Validate the token using the authenticator
			parsedClaims, err := authenticator.VerifyToken(token, claims)
			if err != nil {
				logger.Errorf("Invalid token: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			authContext := context.WithValue(r.Context(), AuthTokenKey, parsedClaims)
			// If token is valid, proceed to the next handler
			next.ServeHTTP(w, r.WithContext(authContext))
		})
	}
}

// DefaultRequestLoggerMiddleware logs the request details including a unique request ID.
// It uses the Logger instance to log the start and completion of each request.
// The request ID is generated using the uuid package and is stored in the request context.
// This middleware can be used to trace requests through the application.
// It logs the request method, URL, and duration of the request.
// This is useful for debugging and monitoring purposes.
func DefaultRequestLoggerMiddleware(logger ILogger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := uuid.New().String()
			ctx := context.WithValue(r.Context(), RequestIDKey, requestId)
			r = r.WithContext(ctx)

			start := time.Now()
			logger.Tracef("Start Request ID: %s, Method: %s, URL: %s", requestId, r.Method, r.URL.Path)
			defer func() {
				duration := time.Since(start)
				logger.Tracef("End Request ID: %s, Method: %s, URL: %s, Duration: %s", requestId, r.Method, r.URL.Path, duration)
			}()
			next.ServeHTTP(w, r)
		})
	}
}
