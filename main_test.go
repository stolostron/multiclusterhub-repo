package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLiveness(t *testing.T) {

	handler := http.HandlerFunc(livenessHandler)
	req := httptest.NewRequest("GET", "/liveness", nil)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestReadiness(t *testing.T) {
	handler := http.HandlerFunc(readinessHandler)
	req := httptest.NewRequest("GET", "/readiness", nil)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
