// SPDX-License-Identifier: Apache-2.0

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/version"
)

func TestVersion(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "successful version request",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create router and register the handler
			r := gin.New()
			r.GET("/version", Version)

			// Create request
			req, err := http.NewRequest(http.MethodGet, "/version", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			r.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Version() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Check response body contains version information
			var v version.Version
			if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
				t.Errorf("failed to unmarshal response: %v", err)
			}

			// Verify basic version structure
			if v.Canonical == "" {
				t.Error("version.Canonical should not be empty")
			}

			if v.Metadata.Architecture == "" {
				t.Error("version.Metadata.Architecture should not be empty")
			}
		})
	}
}

func TestVersion_ContentType(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/version", Version)

	// Create request
	req, err := http.NewRequest(http.MethodGet, "/version", nil)
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
		t.Errorf("Version() content-type = %v, want application/json; charset=utf-8", contentType)
	}
}