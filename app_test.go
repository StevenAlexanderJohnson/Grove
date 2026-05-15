package grove_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/StevenAlexanderJohnson/grove"
)

type testLogger struct {
	infos    []string
	warnings []string
	errors   []string
	debugs   []string
	traces   []string
	fatals   []string
	logs     []string
}

func (l *testLogger) Log(v ...any) {
	l.logs = append(l.logs, joinArgs(v...))
}

func (l *testLogger) Logf(format string, v ...any) {
	l.logs = append(l.logs, sprintf(format, v...))
}

func (l *testLogger) Info(v ...any) {
	l.infos = append(l.infos, joinArgs(v...))
}

func (l *testLogger) Infof(format string, v ...any) {
	l.infos = append(l.infos, sprintf(format, v...))
}

func (l *testLogger) Error(v ...any) {
	l.errors = append(l.errors, joinArgs(v...))
}

func (l *testLogger) Errorf(format string, v ...any) {
	l.errors = append(l.errors, sprintf(format, v...))
}

func (l *testLogger) Debug(v ...any) {
	l.debugs = append(l.debugs, joinArgs(v...))
}

func (l *testLogger) Debugf(format string, v ...any) {
	l.debugs = append(l.debugs, sprintf(format, v...))
}

func (l *testLogger) Warning(v ...any) {
	l.warnings = append(l.warnings, joinArgs(v...))
}

func (l *testLogger) Warningf(format string, v ...any) {
	l.warnings = append(l.warnings, sprintf(format, v...))
}

func (l *testLogger) Trace(v ...any) {
	l.traces = append(l.traces, joinArgs(v...))
}

func (l *testLogger) Tracef(format string, v ...any) {
	l.traces = append(l.traces, sprintf(format, v...))
}

func (l *testLogger) Fatal(v ...any) {
	l.fatals = append(l.fatals, joinArgs(v...))
}

func (l *testLogger) Fatalf(format string, v ...any) {
	l.fatals = append(l.fatals, sprintf(format, v...))
}

func sprintf(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func joinArgs(args ...any) string {
	var b strings.Builder

	for i, arg := range args {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(anyToString(arg))
	}

	return b.String()
}

func anyToString(v any) string {
	switch value := v.(type) {
	case string:
		return value
	case int:
		return strconv.Itoa(value)
	default:
		return ""
	}
}

type testController struct {
	pattern string
	body    string
}

func (c testController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(c.pattern, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(c.body))
	})
}

func TestAppWithRouteRegistersRoute(t *testing.T) {
	app := grove.NewApp("test")

	app.WithRoute("GET /health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("body = %q; want %q", rec.Body.String(), "ok")
	}
}

func TestAppWithRouteEmptyPatternDoesNotPanic(t *testing.T) {
	app := grove.NewApp("test")

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("WithRoute empty pattern panicked: %v", r)
		}
	}()

	app.WithRoute("", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
}

func TestAppWithRouteNilHandlerDoesNotRegister(t *testing.T) {
	logger := &testLogger{}
	app := grove.NewApp("test").WithLogger(logger)

	app.WithRoute("GET /missing", nil)

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusNotFound)
	}
	if len(logger.warnings) == 0 {
		t.Fatalf("expected warning for nil handler")
	}
}

func TestAppWithMiddlewareRunsInRegisteredOrder(t *testing.T) {
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

	app := grove.NewApp("test").
		WithMiddleware(first).
		WithMiddleware(second)

	app.WithRoute("GET /test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, "handler")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

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

func TestAppWithNilMiddlewareDoesNotRegister(t *testing.T) {
	logger := &testLogger{}

	app := grove.NewApp("test").WithLogger(logger)
	app.WithMiddleware(nil)

	if len(logger.warnings) == 0 {
		t.Fatalf("expected warning for nil middleware")
	}
}

func TestAppWithControllerRegistersRoutes(t *testing.T) {
	app := grove.NewApp("test")

	app.WithController(testController{
		pattern: "GET /controller",
		body:    "controller ok",
	})

	req := httptest.NewRequest(http.MethodGet, "/controller", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "controller ok" {
		t.Fatalf("body = %q; want %q", rec.Body.String(), "controller ok")
	}
}

func TestAppWithNilControllerLogsError(t *testing.T) {
	logger := &testLogger{}

	app := grove.NewApp("test").WithLogger(logger)
	app.WithController(nil)

	if len(logger.errors) == 0 {
		t.Fatalf("expected error for nil controller")
	}
}

func TestAppWithControllerFactoryRegistersController(t *testing.T) {
	app := grove.NewApp("test")

	app.WithControllerFactory(func(deps *grove.Dependencies) grove.IController {
		if deps == nil {
			t.Fatalf("deps = nil; want dependency container")
		}

		return testController{
			pattern: "GET /factory",
			body:    "factory ok",
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/factory", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "factory ok" {
		t.Fatalf("body = %q; want %q", rec.Body.String(), "factory ok")
	}
}

func TestAppWithControllerFactoryNilControllerLogsError(t *testing.T) {
	logger := &testLogger{}

	app := grove.NewApp("test").WithLogger(logger)

	app.WithControllerFactory(func(deps *grove.Dependencies) grove.IController {
		return nil
	})

	if len(logger.errors) == 0 {
		t.Fatalf("expected error for nil controller from factory")
	}
}

func TestAppWithMuxUsesProvidedMux(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /custom", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("custom mux"))
	})

	app := grove.NewApp("test").WithMux(mux)

	req := httptest.NewRequest(http.MethodGet, "/custom", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "custom mux" {
		t.Fatalf("body = %q; want %q", rec.Body.String(), "custom mux")
	}
}

func TestAppWithNilMuxKeepsExistingMux(t *testing.T) {
	logger := &testLogger{}

	app := grove.NewApp("test").WithLogger(logger)
	app.WithMux(nil)

	app.WithRoute("GET /existing", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("existing"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/existing", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
	}
	if len(logger.warnings) == 0 {
		t.Fatalf("expected warning for nil mux")
	}
}

func TestAppWithLoggerNilKeepsExistingLogger(t *testing.T) {
	app := grove.NewApp("test")

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("WithLogger(nil) panicked: %v", r)
		}
	}()

	app.WithLogger(nil)
	app.WithPort("")
}

func TestAppWithDependenciesNilDoesNotPanic(t *testing.T) {
	app := grove.NewApp("test")

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("WithDependencies(nil) panicked: %v", r)
		}
	}()

	app.WithDependencies(nil)
}

func TestAppWithDependenciesUsesProvidedContainerInFactory(t *testing.T) {
	app := grove.NewApp("test")
	deps := grove.NewDependencies()

	app.WithDependencies(deps)

	app.WithControllerFactory(func(got *grove.Dependencies) grove.IController {
		if got != deps {
			t.Fatalf("deps = %p; want %p", got, deps)
		}

		return testController{
			pattern: "GET /deps",
			body:    "deps ok",
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/deps", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
	}
}

func TestAppWithScopeMountsScope(t *testing.T) {
	app := grove.NewApp("test")
	scope := grove.NewScope("api")

	scope.WithRoute("GET /me", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("me"))
	}))

	app.WithScope("/api", scope)

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "me" {
		t.Fatalf("body = %q; want %q", rec.Body.String(), "me")
	}
}

func TestAppWithScopeNormalizesMountPath(t *testing.T) {
	tests := []string{
		"api",
		"/api",
		"/api/",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			app := grove.NewApp("test")
			scope := grove.NewScope("api")

			scope.WithRoute("GET /me", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("me"))
			}))

			app.WithScope(path, scope)

			req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
			rec := httptest.NewRecorder()

			app.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
			}
			if rec.Body.String() != "me" {
				t.Fatalf("body = %q; want %q", rec.Body.String(), "me")
			}
		})
	}
}

func TestAppWithRootScope(t *testing.T) {
	app := grove.NewApp("test")
	scope := grove.NewScope("root")

	scope.WithRoute("GET /me", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("root me"))
	}))

	app.WithScope("/", scope)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "root me" {
		t.Fatalf("body = %q; want %q", rec.Body.String(), "root me")
	}
}

func TestAppWithNilScopeDoesNotRegister(t *testing.T) {
	logger := &testLogger{}

	app := grove.NewApp("test").WithLogger(logger)
	app.WithScope("/api", nil)

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusNotFound)
	}
	if len(logger.warnings) == 0 {
		t.Fatalf("expected warning for nil scope")
	}
}
