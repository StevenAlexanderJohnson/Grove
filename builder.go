package grove

import "net/http"

func (app *App) WithController(controller IController) *App {
	controller.RegisterRoutes(app.mux)
	return app
}

func (app *App) WithControllerFactory(factory func(dependencies *Dependencies) IController) *App {
	controller := factory(app.deps)
	if controller == nil {
		app.logger.Error("Controller factory returned nil")
		return app
	}
	return app.WithController(controller)
}

func (app *App) WithPort(port string) *App {
	app.port = port
	return app
}

func (app *App) WithLogger(logger ILogger) *App {
	app.logger = logger
	return app
}

func (app *App) WithMiddleware(mw Middleware) *App {
	app.middleware = append(app.middleware, mw)
	return app
}

func (app *App) WithDependencies(deps *Dependencies) *App {
	app.deps = deps
	return app
}

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
