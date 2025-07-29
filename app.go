package grove

import "net/http"

type App struct {
	port       string
	mux        *http.ServeMux
	middleware []Middleware
	logger     ILogger
	deps       *Dependencies
}

func NewApp() *App {
	return &App{
		port:       "8080",
		mux:        http.NewServeMux(),
		middleware: []Middleware{},
		logger:     NewDefaultLogger(),
		deps:       NewDependencies(),
	}
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.mux.ServeHTTP(w, r)
}

func (app *App) Run() error {
	app.logger.Info("Starting server on port", app.port)
	addr := ":" + app.port
	var handler http.Handler = app.mux
	for _, mw := range app.middleware {
		handler = mw(handler)
	}
	return http.ListenAndServe(addr, handler)
}
