package grove

import "net/http"

// WithController registers a controller with the application.
// It calls the RegisterRoutes method of the controller to set up its routes.
// If the controller is nil, it logs an error and does not register it.
// Controller routes that are registered directly to the App will use all the middleware registered to the app.
func (app *App) WithController(controller IController) *App {
	if controller == nil {
		app.logger.Error("Controller is nil, cannot register")
		return app
	}
	controller.RegisterRoutes(app.mux)
	return app
}

// WithControllerFactory registers a controller using a factory function.
// The factory is called with the application's dependencies to create the controller.
// If the factory returns nil, it logs an error and does not register the controller.
// This method allows for dynamic controller creation based on the application's dependencies.
// This method should be used if you are not manually bootstrapping the controller.
func (app *App) WithControllerFactory(factory ControllerFactory) *App {
	controller := factory(app.deps)
	if controller == nil {
		app.logger.Error("Controller factory returned nil")
		return app
	}
	return app.WithController(controller)
}

// WithPort sets the port for the application.
// If the port is empty, it logs a warning and defaults to "8080".
// This method allows the application to listen on a specific port for incoming HTTP requests.
func (app *App) WithPort(port string) *App {
	if port == "" {
		app.logger.Warning("Warning: Attempting to set an empty port, defaulting to '8080'")
		port = "8080"
	}
	app.port = port
	return app
}

// WithLogger sets the logger for the application.
// If the logger is nil, it logs a warning and uses the default logger.
// This method allows the application to use a custom logger for logging messages.
func (app *App) WithLogger(logger ILogger) *App {
	if logger == nil {
		app.logger.Warning("Warning: Attempting to set a nil logger, using default logger")
		logger = NewDefaultLogger()
	}
	app.logger = logger
	return app
}

// WithMiddleware registers a middleware function to the application.
// If the middleware is nil, it logs a warning and does not register it.
// This method is used to add middleware that can modify the request/response cycle.
// The middleware will be applied to all routes handled by the application.
func (app *App) WithMiddleware(mw Middleware) *App {
	if mw == nil {
		app.logger.Warning("Warning: Attempting to register a nil middleware")
		return app
	}
	app.middleware = append(app.middleware, mw)
	return app
}

// WithRoute registers a handler for a specific path.
// It ensures the path starts and ends with a slash.
// If the handler is nil, it logs a warning and does not register the route.
// This method is used to add routes outside of controllers.
// The it will use all the middleware registered to the app.
func (app *App) WithRoute(path string, handler http.Handler) *App {
	if path == "" || path[0] != '/' {
		path = "/" + path
	}
	if path[len(path)-1] != '/' {
		path += "/"
	}
	if handler == nil {
		app.logger.Warning("Warning: Attempting to register a nil handler at", path)
		return app
	}
	app.mux.Handle(path, http.StripPrefix(path[:len(path)-1], handler))
	return app
}

func (app *App) WithDependencies(deps *Dependencies) *App {
	if deps == nil {
		app.logger.Warning("Warning: Attempting to set nil dependencies, using existing dependencies")
		return app
	}
	app.deps = deps
	return app
}

// WithScope registers a scope at the specified path.
// It ensures the path starts and ends with a slash.
// If the scope is nil, it logs a warning and does not register the scope.
// This method is used to add scopes that can be used for grouping routes or resources.
func (app *App) WithScope(path string, scope *Scope) *App {
	if path == "" || path[0] != '/' {
		path = "/" + path
	}
	if path[len(path)-1] != '/' {
		path += "/"
	}
	if scope == nil {
		app.logger.Warning("Warning: Attempting to register a nil scope at", path)
		return app
	}
	app.mux.Handle(path, http.StripPrefix(path[:len(path)-1], scope))
	return app
}
