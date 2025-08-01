package grove

import "net/http"

type Scope struct {
	mux        *http.ServeMux
	middleware []Middleware
}

func NewScope() *Scope {
	return &Scope{
		mux:        http.NewServeMux(),
		middleware: []Middleware{},
	}
}

func (s *Scope) WithMiddleware(mw Middleware) *Scope {
	s.middleware = append(s.middleware, mw)
	return s
}

func (s *Scope) WithRoute(pattern string, handler http.Handler) *Scope {
	s.mux.Handle(pattern, handler)
	return s
}

func (s *Scope) WithController(controller IController) *Scope {
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
