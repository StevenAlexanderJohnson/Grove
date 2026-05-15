package grove_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/StevenAlexanderJohnson/grove"
)

func TestScopeWithRouteRegistersRoute(t *testing.T) {
	scope := grove.NewScope("test")

	scope.WithRoute("GET /me", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("me"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()

	scope.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "me" {
		t.Fatalf("body = %q; want %q", rec.Body.String(), "me")
	}
}

func TestScopeWithRouteEmptyPatternDoesNotRegister(t *testing.T) {
	logger := &testLogger{}
	scope := grove.NewScope("test", logger)

	scope.WithRoute("", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("should not run"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	scope.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusNotFound)
	}
	if len(logger.warnings) == 0 {
		t.Fatalf("expected warning for empty route pattern")
	}
}

func TestScopeWithRouteNilHandlerDoesNotRegister(t *testing.T) {
	logger := &testLogger{}
	scope := grove.NewScope("test", logger)

	scope.WithRoute("GET /nil", nil)

	req := httptest.NewRequest(http.MethodGet, "/nil", nil)
	rec := httptest.NewRecorder()

	scope.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusNotFound)
	}
	if len(logger.warnings) == 0 {
		t.Fatalf("expected warning for nil handler")
	}
}

func TestScopeWithMiddlewareRunsInRegisteredOrder(t *testing.T) {
	var calls []string

	first := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "first before")
			next.ServeHTTP(w, r)
			calls = append(calls, "first after")
		})
	}

	second := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "second before")
			next.ServeHTTP(w, r)
			calls = append(calls, "second after")
		})
	}

	scope := grove.NewScope("test").
		WithMiddleware(first).
		WithMiddleware(second)

	scope.WithRoute("GET /test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, "handler")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	scope.ServeHTTP(rec, req)

	want := []string{
		"first before",
		"second before",
		"handler",
		"second after",
		"first after",
	}

	if strings.Join(calls, ",") != strings.Join(want, ",") {
		t.Fatalf("calls = %v; want %v", calls, want)
	}
}

func TestScopeWithNilMiddlewareDoesNotRegister(t *testing.T) {
	var calls []string

	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "middleware")
			next.ServeHTTP(w, r)
		})
	}

	scope := grove.NewScope("test").
		WithMiddleware(nil).
		WithMiddleware(mw)

	scope.WithRoute("GET /test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, "handler")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	scope.ServeHTTP(rec, req)

	want := []string{"middleware", "handler"}

	if strings.Join(calls, ",") != strings.Join(want, ",") {
		t.Fatalf("calls = %v; want %v", calls, want)
	}
}

func TestScopeWithControllerRegistersRoutes(t *testing.T) {
	scope := grove.NewScope("test")

	scope.WithController(testController{
		pattern: "GET /controller",
		body:    "controller ok",
	})

	req := httptest.NewRequest(http.MethodGet, "/controller", nil)
	rec := httptest.NewRecorder()

	scope.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "controller ok" {
		t.Fatalf("body = %q; want %q", rec.Body.String(), "controller ok")
	}
}

func TestScopeWithNilControllerDoesNotRegister(t *testing.T) {
	logger := &testLogger{}
	scope := grove.NewScope("test", logger)

	scope.WithController(nil)

	req := httptest.NewRequest(http.MethodGet, "/controller", nil)
	rec := httptest.NewRecorder()

	scope.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusNotFound)
	}
	if len(logger.warnings) == 0 {
		t.Fatalf("expected warning for nil controller")
	}
}

func TestScopeWithNestedScopeMountsScope(t *testing.T) {
	api := grove.NewScope("api")
	users := grove.NewScope("users")

	users.WithRoute("GET /me", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("nested me"))
	}))

	api.WithScope("/users", users)

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	rec := httptest.NewRecorder()

	api.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "nested me" {
		t.Fatalf("body = %q; want %q", rec.Body.String(), "nested me")
	}
}

func TestScopeWithNestedScopeNormalizesMountPath(t *testing.T) {
	tests := []string{
		"users",
		"/users",
		"/users/",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			api := grove.NewScope("api")
			users := grove.NewScope("users")

			users.WithRoute("GET /me", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("nested me"))
			}))

			api.WithScope(path, users)

			req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
			rec := httptest.NewRecorder()

			api.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
			}
			if rec.Body.String() != "nested me" {
				t.Fatalf("body = %q; want %q", rec.Body.String(), "nested me")
			}
		})
	}
}

func TestScopeWithRootNestedScope(t *testing.T) {
	api := grove.NewScope("api")
	protected := grove.NewScope("protected")

	api.WithRoute("POST /login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("login"))
	}))

	protected.WithRoute("GET /me", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("protected me"))
	}))

	api.WithScope("/", protected)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()

	api.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "protected me" {
		t.Fatalf("body = %q; want %q", rec.Body.String(), "protected me")
	}
}

func TestScopeWithRootNestedScopeDoesNotOverrideMoreSpecificRoute(t *testing.T) {
	api := grove.NewScope("api")
	protected := grove.NewScope("protected")

	api.WithRoute("POST /login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("login"))
	}))

	protected.WithRoute("GET /me", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("protected me"))
	}))

	api.WithScope("/", protected)

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	rec := httptest.NewRecorder()

	api.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "login" {
		t.Fatalf("body = %q; want %q", rec.Body.String(), "login")
	}
}

func TestScopeWithNilNestedScopeDoesNotRegister(t *testing.T) {
	logger := &testLogger{}
	api := grove.NewScope("api", logger)

	api.WithScope("/users", nil)

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	rec := httptest.NewRecorder()

	api.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusNotFound)
	}
	if len(logger.warnings) == 0 {
		t.Fatalf("expected warning for nil nested scope")
	}
}

func TestScopeMiddlewareAppliesToNestedScope(t *testing.T) {
	var calls []string

	apiMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "api middleware")
			next.ServeHTTP(w, r)
		})
	}

	usersMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "users middleware")
			next.ServeHTTP(w, r)
		})
	}

	api := grove.NewScope("api").WithMiddleware(apiMiddleware)
	users := grove.NewScope("users").WithMiddleware(usersMiddleware)

	users.WithRoute("GET /me", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, "handler")
	}))

	api.WithScope("/users", users)

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	rec := httptest.NewRecorder()

	api.ServeHTTP(rec, req)

	want := []string{
		"api middleware",
		"users middleware",
		"handler",
	}

	if strings.Join(calls, ",") != strings.Join(want, ",") {
		t.Fatalf("calls = %v; want %v", calls, want)
	}
}

func TestScopePublicAndProtectedRoutesUnderSameBaseScope(t *testing.T) {
	var calls []string

	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "auth")
			next.ServeHTTP(w, r)
		})
	}

	api := grove.NewScope("api")
	protected := grove.NewScope("protected").WithMiddleware(authMiddleware)

	api.WithRoute("POST /login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, "login")
		_, _ = w.Write([]byte("login"))
	}))

	protected.WithRoute("GET /me", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, "me")
		_, _ = w.Write([]byte("me"))
	}))

	api.WithScope("/", protected)

	loginReq := httptest.NewRequest(http.MethodPost, "/login", nil)
	loginRec := httptest.NewRecorder()

	api.ServeHTTP(loginRec, loginReq)

	if loginRec.Code != http.StatusOK {
		t.Fatalf("login status = %d; want %d", loginRec.Code, http.StatusOK)
	}
	if strings.Join(calls, ",") != "login" {
		t.Fatalf("login calls = %v; want [login]", calls)
	}

	calls = nil

	meReq := httptest.NewRequest(http.MethodGet, "/me", nil)
	meRec := httptest.NewRecorder()

	api.ServeHTTP(meRec, meReq)

	if meRec.Code != http.StatusOK {
		t.Fatalf("me status = %d; want %d", meRec.Code, http.StatusOK)
	}
	if strings.Join(calls, ",") != "auth,me" {
		t.Fatalf("me calls = %v; want [auth me]", calls)
	}
}
