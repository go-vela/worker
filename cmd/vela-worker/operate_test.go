// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/queue"
)

func TestWorkerRegistryConfiguration(t *testing.T) {
	// Test worker registry configuration logic from operate function
	w := &Worker{
		Config: &Config{
			Build: &Build{
				Limit: 5,
			},
			Queue: &queue.Setup{
				Routes: []string{"vela", "test"},
			},
		},
	}

	// Test build limit setting logic
	registryWorker := new(api.Worker)
	registryWorker.SetHostname("test-worker")
	registryWorker.SetActive(true)

	// Test normal build limit
	if w.Config.Build.Limit <= int(^uint32(0)>>1) {
		registryWorker.SetBuildLimit(int32(w.Config.Build.Limit))
		if registryWorker.GetBuildLimit() != 5 {
			t.Errorf("Build limit = %v, want 5", registryWorker.GetBuildLimit())
		}
	}

	// Test routes setting
	if len(w.Config.Queue.Routes) > 0 && w.Config.Queue.Routes[0] != "NONE" && w.Config.Queue.Routes[0] != "" {
		registryWorker.SetRoutes(w.Config.Queue.Routes)
		routes := registryWorker.GetRoutes()
		if len(routes) != 2 {
			t.Errorf("Routes length = %v, want 2", len(routes))
		}
		if routes[0] != "vela" || routes[1] != "test" {
			t.Errorf("Routes = %v, want [vela test]", routes)
		}
	}
}

func TestWorkerRegistryLimitBoundary(t *testing.T) {
	// Test the upper bound logic for build limits
	w := &Worker{
		Config: &Config{
			Build: &Build{
				Limit: int(^uint32(0)>>1) + 1, // Exceed max int32
			},
		},
	}

	registryWorker := new(api.Worker)

	// This should clamp to max int32
	if w.Config.Build.Limit > int(^uint32(0)>>1) {
		registryWorker.SetBuildLimit(int32(^uint32(0) >> 1))
		expectedMax := int32(^uint32(0) >> 1)
		if registryWorker.GetBuildLimit() != expectedMax {
			t.Errorf("Build limit = %v, want %v", registryWorker.GetBuildLimit(), expectedMax)
		}
	}
}

func TestWorkerRegistryNoRoutes(t *testing.T) {
	// Test routes handling with empty or NONE values
	tests := []struct {
		name   string
		routes []string
		expect bool
	}{
		{
			name:   "empty routes",
			routes: []string{},
			expect: false,
		},
		{
			name:   "NONE routes",
			routes: []string{"NONE"},
			expect: false,
		},
		{
			name:   "empty string routes",
			routes: []string{""},
			expect: false,
		},
		{
			name:   "valid routes",
			routes: []string{"vela"},
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				Config: &Config{
					Queue: &queue.Setup{
						Routes: tt.routes,
					},
				},
			}

			registryWorker := new(api.Worker)

			shouldSetRoutes := len(w.Config.Queue.Routes) > 0 && w.Config.Queue.Routes[0] != "NONE" && w.Config.Queue.Routes[0] != ""

			if shouldSetRoutes != tt.expect {
				t.Errorf("Should set routes = %v, want %v", shouldSetRoutes, tt.expect)
			}

			if shouldSetRoutes {
				registryWorker.SetRoutes(w.Config.Queue.Routes)
				if len(registryWorker.GetRoutes()) == 0 {
					t.Errorf("Routes were not set when they should have been")
				}
			}
		})
	}
}

func TestWorkerStatusUpdate(t *testing.T) {
	// Test worker status update patterns from operate function
	registryWorker := new(api.Worker)
	registryWorker.SetHostname("test-worker")

	// Test error status setting
	registryWorker.SetStatus(constants.WorkerStatusError)
	if registryWorker.GetStatus() != constants.WorkerStatusError {
		t.Errorf("Worker status = %v, want %v", registryWorker.GetStatus(), constants.WorkerStatusError)
	}

	// Test active setting
	registryWorker.SetActive(true)
	if !registryWorker.GetActive() {
		t.Errorf("Worker active = %v, want true", registryWorker.GetActive())
	}

	// Test hostname setting
	if registryWorker.GetHostname() != "test-worker" {
		t.Errorf("Worker hostname = %v, want test-worker", registryWorker.GetHostname())
	}
}