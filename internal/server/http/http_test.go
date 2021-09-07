package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gpiecyk/data-warehouse/internal/api"
)

// strings_test.go is a good example

func TestHandlers_Health(t *testing.T) {
	request := httptest.NewRequest("GET", "/health", nil)
	writer := httptest.NewRecorder()

	handlers := Handlers{api: &api.API{}}
	healthHandler := http.HandlerFunc(handlers.Health)

	healthHandler.ServeHTTP(writer, request)

	type health struct {
		Commit     string
		Env        string
		Status     string
		Version    string
		ReleasedOn string
		StartedAt  string
	}

	expectedHealth := health{
		Commit:  "<git commit hash>",
		Env:     "testing",
		Status:  "all systems up and running",
		Version: "v0.1.0",
	}

	actual := new(health)
	if err := json.NewDecoder(writer.Body).Decode(actual); err != nil {
		t.Error("cannot decode response body to the object")
	}

	if expectedHealth.Commit != actual.Commit {
		t.Error("field commit is wrong")
	}

	if expectedHealth.Env != actual.Env {
		t.Error("field env is wrong")
	}

	if expectedHealth.Status != actual.Status {
		t.Error("field status is wrong")
	}

	if expectedHealth.Version != actual.Version {
		t.Error("field version is wrong")
	}

	if status := writer.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
