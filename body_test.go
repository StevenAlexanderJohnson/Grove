package grove_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/StevenAlexanderJohnson/grove"
)

type testJSONBody struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func TestParseJsonBodyFromRequest(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/test",
		strings.NewReader(`{"email":"testing@example.com","name":"Testing"}`),
	)

	got, err := grove.ParseJsonBodyFromRequest[testJSONBody](req)
	if err != nil {
		t.Fatalf("ParseJsonBodyFromRequest() error = %v; want nil", err)
	}

	if got.Email != "testing@example.com" {
		t.Fatalf("Email = %q; want %q", got.Email, "testing@example.com")
	}

	if got.Name != "Testing" {
		t.Fatalf("Name = %q; want %q", got.Name, "Testing")
	}
}

func TestParseJsonBodyFromRequestWithInvalidJsonShouldFail(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/test",
		strings.NewReader(`{"email":`),
	)

	_, err := grove.ParseJsonBodyFromRequest[testJSONBody](req)
	if err == nil {
		t.Fatalf("ParseJsonBodyFromRequest() error = nil; want error")
	}
}

func TestWriteJsonBodyToResponse(t *testing.T) {
	rec := httptest.NewRecorder()

	err := grove.WriteJsonBodyToResponse(rec, testJSONBody{
		Email: "testing@example.com",
		Name:  "Testing",
	})
	if err != nil {
		t.Fatalf("WriteJsonBodyToResponse() error = %v; want nil", err)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q; want %q", got, "application/json")
	}

	var body testJSONBody
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("response body failed to decode: %v", err)
	}

	if body.Email != "testing@example.com" {
		t.Fatalf("Email = %q; want %q", body.Email, "testing@example.com")
	}

	if body.Name != "Testing" {
		t.Fatalf("Name = %q; want %q", body.Name, "Testing")
	}
}

func TestWriteJsonBodyToResponseWithUnencodableBodyShouldFail(t *testing.T) {
	rec := httptest.NewRecorder()

	body := map[string]any{
		"bad": func() {},
	}

	err := grove.WriteJsonBodyToResponse(rec, body)
	if err == nil {
		t.Fatalf("WriteJsonBodyToResponse() error = nil; want error")
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q; want %q", got, "application/json")
	}
}

func TestWriteErrorToResponse(t *testing.T) {
	rec := httptest.NewRecorder()

	grove.WriteErrorToResponse(rec, http.StatusUnauthorized, "not authorized")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusUnauthorized)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q; want %q", got, "application/json")
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("response body failed to decode: %v", err)
	}

	if body["error"] != "not authorized" {
		t.Fatalf("error = %q; want %q", body["error"], "not authorized")
	}
}
