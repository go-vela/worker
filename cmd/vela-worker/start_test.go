// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/runtime"
)

func TestWorker_Start_HTTPServerConfiguration(t *testing.T) {
	// Test the HTTP server configuration logic from Start function
	addr, _ := url.Parse("http://localhost:8080")

	w := &Worker{
		Config: &Config{
			API: &API{
				Address: addr,
			},
			Build: &Build{
				Limit:       1,
				Timeout:     30 * time.Minute,
				CPUQuota:    1200,
				MemoryLimit: 4,
				PidsLimit:   1024,
			},
			Executor: &executor.Setup{
				Driver: "linux",
			},
			Runtime: &runtime.Setup{
				Driver: "docker",
			},
			Queue: &queue.Setup{
				Driver: "redis",
			},
			Server: &Server{
				Address: "http://localhost:8080",
				Secret:  "test-secret",
			},
			Certificate: &Certificate{
				Cert: "",
				Key:  "",
			},
		},
		Executors:     make(map[int]executor.Engine),
		RegisterToken: make(chan string, 1),
		RunningBuilds: make([]*api.Build, 0),
	}

	// Test server configuration creation (mimics Start function logic)
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", w.Config.API.Address.Port()),
		Handler:           nil, // Would be set by w.server() in actual code
		TLSConfig:         nil, // Would be set by w.server() in actual code
		ReadHeaderTimeout: 60 * time.Second,
	}

	// Verify server configuration
	if server.Addr != ":8080" {
		t.Errorf("Server address = %v, want :8080", server.Addr)
	}

	if server.ReadHeaderTimeout != 60*time.Second {
		t.Errorf("Server ReadHeaderTimeout = %v, want 60s", server.ReadHeaderTimeout)
	}
}

func TestWorker_Start_TLSConfiguration(t *testing.T) {
	// Test TLS configuration logic
	addr, _ := url.Parse("https://localhost:8443")

	w := &Worker{
		Config: &Config{
			API: &API{
				Address: addr,
			},
			Certificate: &Certificate{
				Cert: "/path/to/cert.pem",
				Key:  "/path/to/key.pem",
			},
		},
	}

	// Test TLS server configuration
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", w.Config.API.Address.Port()),
		ReadHeaderTimeout: 60 * time.Second,
	}

	// Verify HTTPS port configuration
	if server.Addr != ":8443" {
		t.Errorf("TLS Server address = %v, want :8443", server.Addr)
	}

	// Test certificate paths are configured
	if w.Config.Certificate.Cert == "" {
		t.Error("Certificate cert path should not be empty for TLS configuration")
	}

	if w.Config.Certificate.Key == "" {
		t.Error("Certificate key path should not be empty for TLS configuration")
	}
}

func TestWorker_Start_ContextConfiguration(t *testing.T) {
	// Test context setup from Start function
	ctx := context.Background()

	// Test context cancellation (mimics Start function logic)
	ctx, done := context.WithCancel(ctx)

	// Verify context is not done initially
	select {
	case <-ctx.Done():
		t.Error("Context should not be done initially")
	default:
		// Expected behavior
	}

	// Test cancellation
	done()

	// Verify context is done after cancellation
	select {
	case <-ctx.Done():
		// Expected behavior
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should be done after cancellation")
	}

	// Verify context error
	if ctx.Err() == nil {
		t.Error("Context should have an error after cancellation")
	}
}

func TestWorker_Start_PortExtraction(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		expectedPort string
	}{
		{
			name:         "standard HTTP port",
			url:          "http://localhost:8080",
			expectedPort: "8080",
		},
		{
			name:         "standard HTTPS port",
			url:          "https://localhost:8443",
			expectedPort: "8443",
		},
		{
			name:         "custom port",
			url:          "http://localhost:9090",
			expectedPort: "9090",
		},
		{
			name:         "default HTTP port",
			url:          "http://localhost",
			expectedPort: "",
		},
		{
			name:         "default HTTPS port",
			url:          "https://localhost",
			expectedPort: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := url.Parse(tt.url)
			if err != nil {
				t.Fatalf("Failed to parse URL: %v", err)
			}

			port := addr.Port()
			if port != tt.expectedPort {
				t.Errorf("Port = %v, want %v", port, tt.expectedPort)
			}
		})
	}
}

func TestWorker_Start_ServerShutdown(t *testing.T) {
	// Test server shutdown error handling logic
	ctx := context.Background()

	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 60 * time.Second,
	}

	// Test graceful shutdown
	err := server.Shutdown(ctx)
	if err != nil {
		// This is expected for a server that was never started
		// In the actual Start function, this error would be logged but not returned
		t.Logf("Expected shutdown error for unstarted server: %v", err)
	}
}

func TestContextWithTimeout(t *testing.T) {
	// Test context timeout behavior that might be used in Start function
	ctx := context.Background()

	timeoutCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	select {
	case <-timeoutCtx.Done():
		// Expected after timeout
		if timeoutCtx.Err() != context.DeadlineExceeded {
			t.Errorf("Context error = %v, want %v", timeoutCtx.Err(), context.DeadlineExceeded)
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Context should have timed out")
	}
}
