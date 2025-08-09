// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealth(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful health check",
			expectedStatus: http.StatusOK,
			expectedBody:   "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create router and register the handler
			r := gin.New()
			r.GET("/health", Health)

			// Create request
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/health", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			r.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Health() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Check response body
			var response string
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("failed to unmarshal response: %v", err)
			}

			if response != tt.expectedBody {
				t.Errorf("Health() body = %v, want %v", response, tt.expectedBody)
			}
		})
	}
}

func TestHealth_ContentType(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", Health)

	// Create request
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("Health() content-type = %v, want application/json; charset=utf-8", contentType)
	}
}

func TestHealth_MethodNotAllowed(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", Health)

	// Test that POST is not allowed
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/health", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should return 404 Not Found since route is not defined for POST
	if w.Code != http.StatusNotFound {
		t.Errorf("Health() with POST status = %v, want %v", w.Code, http.StatusNotFound)
	}
}
