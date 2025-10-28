package grove

import "net/http"

// IController defines the interface for controllers in the Grove application.
// Controllers implement this interface to register their routes with the application's HTTP multiplexer.
// The RegisterRoutes method is called with the application's ServeMux and Dependencies to set up the routes
// and any necessary dependencies for the controller.
type IController interface {
	// RegisterRoutes registers the controller's routes with the provided ServeMux.
	// It is called by the Grove application to set up the HTTP routes for the controller.
	RegisterRoutes(mux *http.ServeMux)
}

// Type alias for used with `App.WithControllerFactory`.
// This isn't a recommended method for initializing controller because it can hide errors until runtime,
// not compile time. It is provided for convenience if you would like to use it.
type ControllerFactory func(*Dependencies) IController
