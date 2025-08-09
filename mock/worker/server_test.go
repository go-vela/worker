// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFakeHandler(t *testing.T) {
	// Test that FakeHandler returns a valid http.Handler
	handler := FakeHandler()

	if handler == nil {
		t.Error("FakeHandler() returned nil")
	}

	// Test that the handler can serve HTTP requests
	req, err := http.NewRequest(http.MethodGet, "/api/v1/executors", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Should return a valid response (not 404)
	if w.Code == http.StatusNotFound {
		t.Error("FakeHandler() returned 404 for known route /api/v1/executors")
	}
}

func TestFakeHandler_Routes(t *testing.T) {
	handler := FakeHandler()

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "get executors",
			method: http.MethodGet,
			path:   "/api/v1/executors",
		},
		{
			name:   "get executor",
			method: http.MethodGet,
			path:   "/api/v1/executors/test-executor",
		},
		{
			name:   "get build",
			method: http.MethodGet,
			path:   "/api/v1/executors/test-executor/build",
		},
		{
			name:   "cancel build",
			method: http.MethodDelete,
			path:   "/api/v1/executors/test-executor/build/cancel",
		},
		{
			name:   "get pipeline",
			method: http.MethodGet,
			path:   "/api/v1/executors/test-executor/pipeline",
		},
		{
			name:   "get repo",
			method: http.MethodGet,
			path:   "/api/v1/executors/test-executor/repo",
		},
		{
			name:   "register",
			method: http.MethodPost,
			path:   "/register",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Route should exist (not return 404)
			if w.Code == http.StatusNotFound {
				t.Errorf("FakeHandler() returned 404 for route %s %s", tt.method, tt.path)
			}
		})
	}
}

func TestFakeHandler_UnknownRoute(t *testing.T) {
	handler := FakeHandler()

	req, err := http.NewRequest(http.MethodGet, "/unknown-route", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Unknown routes should return 404
	if w.Code != http.StatusNotFound {
		t.Errorf("FakeHandler() status = %v, want %v for unknown route", w.Code, http.StatusNotFound)
	}
}