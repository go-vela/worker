// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetrics(t *testing.T) {
	// Test that Metrics returns a valid http.Handler
	handler := Metrics()

	if handler == nil {
		t.Error("Metrics() returned nil")
	}

	// Test that the handler can serve HTTP requests
	req, err := http.NewRequest(http.MethodGet, "/metrics", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Prometheus metrics endpoint should return 200 OK
	if w.Code != http.StatusOK {
		t.Errorf("Metrics() status = %v, want %v", w.Code, http.StatusOK)
	}

	// Check that response contains some metrics content
	body := w.Body.String()
	if len(body) == 0 {
		t.Error("Metrics() returned empty response body")
	}

	// Basic check for Prometheus metrics format (should contain # TYPE or # HELP)
	if !containsMetricsContent(body) {
		t.Error("Metrics() response does not appear to be Prometheus metrics format")
	}
}

func TestMetrics_ContentType(t *testing.T) {
	handler := Metrics()

	req, err := http.NewRequest(http.MethodGet, "/metrics", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Prometheus handler should set appropriate content type
	contentType := w.Header().Get("Content-Type")
	if contentType == "" {
		t.Error("Metrics() did not set Content-Type header")
	}
}

// containsMetricsContent checks if the response body contains Prometheus metrics content
func containsMetricsContent(body string) bool {
	// At minimum, should have some content and look like metrics
	return len(body) > 0 && (body[0] == '#' || len(body) > 100)
}