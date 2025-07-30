package grove

import "net/http"

type App struct {
	port       string
	mux        *http.ServeMux
	middleware []Middleware
	logger     ILogger
	deps       *Dependencies
}

func NewApp(appName string) *App {
	return &App{
		port:       "8080",
		mux:        http.NewServeMux(),
		middleware: []Middleware{},
		logger:     NewDefaultLogger(appName),
		deps:       NewDependencies(),
	}
}

// ServeHTTP implements the http.Handler interface.
// It allows the App to be used as a handler in an HTTP server.
// It delegates the request handling to the ServeMux.
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.mux.ServeHTTP(w, r)
}

// Run starts to listen the HTTP server on the specified port.
// It applies all registered middleware to the handler.
// It returns an error if the server fails to start.
func (app *App) Run() error {
	app.logger.Info("Starting server on port", app.port)
	addr := ":" + app.port
	var handler http.Handler = app.mux
	for _, mw := range app.middleware {
		handler = mw(handler)
	}
	return http.ListenAndServe(addr, handler)
}
