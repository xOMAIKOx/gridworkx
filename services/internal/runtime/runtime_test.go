package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMuxReportsTheConfiguredBoundary(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	NewMux(Config{Role: "gridworks-api", Version: "0.1.0"}).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if got := recorder.Body.String(); got != "{\"role\":\"gridworks-api\",\"status\":\"ok\",\"version\":\"0.1.0\"}\n" {
		t.Fatalf("unexpected response: %s", got)
	}
}

func TestNewMuxRejectsNonGetRequests(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/healthz", nil)

	NewMux(Config{Role: "gridworks-worker", Version: "0.1.0"}).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", recorder.Code)
	}
}
