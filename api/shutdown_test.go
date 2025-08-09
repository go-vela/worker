// SPDX-License-Identifier: Apache-2.0

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestShutdown(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "shutdown endpoint returns not implemented",
			expectedStatus: http.StatusNotImplemented,
			expectedBody:   "This endpoint is not yet implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create router and register the handler
			r := gin.New()
			r.POST("/api/v1/shutdown", Shutdown)

			// Create request
			req, err := http.NewRequest(http.MethodPost, "/api/v1/shutdown", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			r.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Shutdown() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Check response body
			var response string
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("failed to unmarshal response: %v", err)
			}

			if response != tt.expectedBody {
				t.Errorf("Shutdown() body = %v, want %v", response, tt.expectedBody)
			}
		})
	}
}

func TestShutdown_ContentType(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/shutdown", Shutdown)

	// Create request
	req, err := http.NewRequest(http.MethodPost, "/api/v1/shutdown", nil)
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
		t.Errorf("Shutdown() content-type = %v, want application/json; charset=utf-8", contentType)
	}
}