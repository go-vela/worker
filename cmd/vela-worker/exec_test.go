// SPDX-License-Identifier: Apache-2.0

package main

import (
	"strings"
	"sync"
	"testing"
	"time"
)

func TestGenerateCryptographicBuildID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "generates unique ID",
		},
		{
			name: "generates hex string",
		},
		{
			name: "generates consistent length",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Generate multiple IDs to test uniqueness
			id1 := generateCryptographicBuildID()
			id2 := generateCryptographicBuildID()

			// Test: IDs should not be empty
			if id1 == "" {
				t.Error("generateCryptographicBuildID returned empty string")
			}

			// Test: IDs should be unique
			if id1 == id2 {
				t.Error("generateCryptographicBuildID returned duplicate IDs")
			}

			// Test: IDs should be hex strings (32 chars for 16 bytes)
			if !strings.Contains(tc.name, "fallback") && len(id1) != 32 {
				t.Errorf("expected ID length of 32 for hex encoding, got %d", len(id1))
			}

			// Test: Validate hex encoding (should only contain hex characters)
			for _, r := range id1 {
				if (r < '0' || r > '9') && (r < 'a' || r > 'f') && r != '-' {
					t.Errorf("ID contains non-hex character: %c", r)
				}
			}
		})
	}
}

func TestWorker_GetBuildResources(t *testing.T) {
	tests := []struct {
		name           string
		cpuQuota       int
		memoryLimit    int
		pidsLimit      int
		expectedCPU    int64
		expectedMemory int64
		expectedPids   int64
	}{
		{
			name:           "standard resources",
			cpuQuota:       2000, // 2 cores in millicores
			memoryLimit:    4,    // 4 GB
			pidsLimit:      1024,
			expectedCPU:    2000,
			expectedMemory: 4294967296, // 4 GB in bytes
			expectedPids:   1024,
		},
		{
			name:           "minimal resources",
			cpuQuota:       500, // 0.5 cores
			memoryLimit:    1,   // 1 GB
			pidsLimit:      256,
			expectedCPU:    500,
			expectedMemory: 1073741824, // 1 GB in bytes
			expectedPids:   256,
		},
		{
			name:           "zero resources",
			cpuQuota:       0,
			memoryLimit:    0,
			pidsLimit:      0,
			expectedCPU:    0,
			expectedMemory: 0,
			expectedPids:   0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := &Worker{
				Config: &Config{
					Build: &Build{
						CPUQuota:    tc.cpuQuota,
						MemoryLimit: tc.memoryLimit,
						PidsLimit:   tc.pidsLimit,
					},
				},
			}

			resources := w.getBuildResources()

			if resources == nil {
				t.Fatal("getBuildResources returned nil")
			}

			if resources.CPUQuota != tc.expectedCPU {
				t.Errorf("expected CPU quota %d, got %d", tc.expectedCPU, resources.CPUQuota)
			}

			if resources.Memory != tc.expectedMemory {
				t.Errorf("expected memory %d, got %d", tc.expectedMemory, resources.Memory)
			}

			if resources.PidsLimit != tc.expectedPids {
				t.Errorf("expected pids limit %d, got %d", tc.expectedPids, resources.PidsLimit)
			}
		})
	}
}

func TestWorker_BuildContextTracking(t *testing.T) {
	t.Run("concurrent build context operations", func(t *testing.T) {
		w := &Worker{
			BuildContexts:      make(map[string]*BuildContext),
			BuildContextsMutex: sync.RWMutex{},
			Config: &Config{
				Build: &Build{
					CPUQuota:    1000,
					MemoryLimit: 2,
					PidsLimit:   512,
				},
			},
		}

		// Test concurrent writes
		var wg sync.WaitGroup

		numGoroutines := 10

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)

			go func(_ int) {
				defer wg.Done()

				buildID := generateCryptographicBuildID()
				buildContext := &BuildContext{
					BuildID:       buildID,
					WorkspacePath: "/tmp/test-" + buildID,
					StartTime:     time.Now(),
					Resources:     w.getBuildResources(),
					Environment:   make(map[string]string),
				}

				// Add context
				w.BuildContextsMutex.Lock()
				w.BuildContexts[buildID] = buildContext
				w.BuildContextsMutex.Unlock()

				// Simulate some work
				time.Sleep(10 * time.Millisecond)

				// Remove context
				w.BuildContextsMutex.Lock()
				delete(w.BuildContexts, buildID)
				w.BuildContextsMutex.Unlock()
			}(i)
		}

		wg.Wait()

		// Verify all contexts were cleaned up
		if len(w.BuildContexts) != 0 {
			t.Errorf("expected 0 build contexts after cleanup, got %d", len(w.BuildContexts))
		}
	})

	t.Run("build context initialization", func(t *testing.T) {
		w := &Worker{
			Config: &Config{
				Build: &Build{
					CPUQuota:    2000,
					MemoryLimit: 4,
					PidsLimit:   1024,
				},
			},
		}

		// Test that BuildContexts map is initialized properly
		if w.BuildContexts == nil {
			w.BuildContexts = make(map[string]*BuildContext)
		}

		buildID := generateCryptographicBuildID()
		buildContext := &BuildContext{
			BuildID:       buildID,
			WorkspacePath: "/tmp/vela-build-" + buildID,
			StartTime:     time.Now(),
			Resources:     w.getBuildResources(),
			Environment:   make(map[string]string),
		}

		w.BuildContextsMutex.Lock()
		w.BuildContexts[buildID] = buildContext
		w.BuildContextsMutex.Unlock()

		// Verify context was added
		w.BuildContextsMutex.RLock()
		ctx, exists := w.BuildContexts[buildID]
		w.BuildContextsMutex.RUnlock()

		if !exists {
			t.Error("build context was not added to map")
		}

		if ctx.BuildID != buildID {
			t.Errorf("expected build ID %s, got %s", buildID, ctx.BuildID)
		}

		if !strings.Contains(ctx.WorkspacePath, buildID) {
			t.Errorf("workspace path should contain build ID")
		}

		if ctx.Resources == nil {
			t.Error("build context resources should not be nil")
		}

		if ctx.Environment == nil {
			t.Error("build context environment should not be nil")
		}
	})
}
