package grove

import (
	"net/http"
)

// Scope is very similar to `App` except it doesn't have a dependency
// container, logger, or port.
// It should be used to segment your code application to apply middleware
// to specific routes. It can also be useful just for logical
// separation. Like `App`, you should not initialize this struct manually
// because it has private fields that should be built using the builder
// functions.
//
// While it was designed to work with App, this struct implements the net/http Handler interface.
// it can be used with the http standard library.
type Scope struct {
	mux        *http.ServeMux
	logger     ILogger
	middleware []Middleware
}

// Initializes a Scope. It sets default values that can be overwritten
// with the builder functions. The parameter `scopeName` is used to initialize the logger.
// The logger can be overwritten if you want to reuse an existing logger.
func NewScope(scopeName string) *Scope {
	return &Scope{
		mux:        http.NewServeMux(),
		logger:     NewDefaultLogger(scopeName),
		middleware: []Middleware{},
	}
}

// WithMiddleware registers a middleware function to the scope.
// Middleware is applied in the order that they were registered.
func (s *Scope) WithMiddleware(mw Middleware) *Scope {
	if mw == nil {
		return s
	}

	s.middleware = append(s.middleware, mw)
	return s
}

// WithLogger allows you to overwrite the default logger initialized from `NewScope`.
// This is useful if you want to use an existing logger.
func (s *Scope) WithLogger(logger ILogger) *Scope {
	if logger == nil {
		s.logger.Warning("Warning: Attempting to set a nil scope logger, using existing logger")
		return s
	}

	s.logger = logger
	return s
}

// WithRoute registers a handler for a specific path.
// It ensures the path starts and ends with a slash. This means that '/' is valid as well.
// If the handler is nil, it logs a warning and does not register the route.
// This method is used to add routes outside of the controllers.
// All middleware that have been registered for the scope are applied to the route in the
// order they were registered.
func (s *Scope) WithRoute(pattern string, handler http.Handler) *Scope {
	if pattern == "" || pattern[0] != '/' {
		pattern = "/" + pattern
	}
	if pattern[len(pattern)-1] != '/' {
		pattern += "/"
	}
	if handler == nil {
		s.logger.Warning("Warning: Attempting to add a nil route to scope.")
		return s
	}

	for _, mw := range s.middleware {
		handler = mw(handler)
	}

	s.mux.Handle(pattern, handler)
	return s
}

// WithController registers a controller with the application.
// It calls the RegisterRoutes method of the controller to set up its routes.
// If the controller is nil, it logs a warning and does not register it.
// Controller routes that are registered here will have the middleware registered
// to the scope applied.
func (s *Scope) WithController(controller IController) *Scope {
	if controller == nil {
		s.logger.Warning("Warning: Attempting to register a nil controller to scope.")
		return s
	}

	controller.RegisterRoutes(s.mux)
	return s
}

func (s *Scope) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = s.mux
	for _, mw := range s.middleware {
		handler = mw(handler)
	}

	handler.ServeHTTP(w, r)
}
