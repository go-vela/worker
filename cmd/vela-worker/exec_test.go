// SPDX-License-Identifier: Apache-2.0

package main

import (
	"sync"
	"testing"
	"time"
)

func TestGenerateCryptographicBuildID(t *testing.T) {
	// Test that generateCryptographicBuildID returns a valid hex string
	id1 := generateCryptographicBuildID()
	if len(id1) != 32 { // 16 bytes = 32 hex characters
		t.Errorf("generateCryptographicBuildID() returned ID with length %d, want 32", len(id1))
	}

	// Test that it generates unique IDs
	id2 := generateCryptographicBuildID()
	if id1 == id2 {
		t.Errorf("generateCryptographicBuildID() returned duplicate IDs: %s", id1)
	}

	// Verify it's valid hex
	for _, c := range id1 {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			t.Errorf("generateCryptographicBuildID() returned non-hex character: %c", c)
		}
	}
}

func TestWorker_getBuildResources(t *testing.T) {
	tests := []struct {
		name        string
		cpuQuota    int
		memoryLimit int
		pidsLimit   int
		wantCPU     int64
		wantMemory  int64
		wantPids    int64
	}{
		{
			name:        "default values",
			cpuQuota:    1200,
			memoryLimit: 4,
			pidsLimit:   1024,
			wantCPU:     1200,
			wantMemory:  4 * 1024 * 1024 * 1024,
			wantPids:    1024,
		},
		{
			name:        "custom values",
			cpuQuota:    2000,
			memoryLimit: 8,
			pidsLimit:   2048,
			wantCPU:     2000,
			wantMemory:  8 * 1024 * 1024 * 1024,
			wantPids:    2048,
		},
		{
			name:        "minimum values",
			cpuQuota:    100,
			memoryLimit: 1,
			pidsLimit:   256,
			wantCPU:     100,
			wantMemory:  1 * 1024 * 1024 * 1024,
			wantPids:    256,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				Config: &Config{
					Build: &Build{
						CPUQuota:    tt.cpuQuota,
						MemoryLimit: tt.memoryLimit,
						PidsLimit:   tt.pidsLimit,
					},
				},
			}

			resources := w.getBuildResources()

			if resources.CPUQuota != tt.wantCPU {
				t.Errorf("getBuildResources() CPUQuota = %v, want %v", resources.CPUQuota, tt.wantCPU)
			}
			if resources.Memory != tt.wantMemory {
				t.Errorf("getBuildResources() Memory = %v, want %v", resources.Memory, tt.wantMemory)
			}

			if resources.PidsLimit != tt.wantPids {
				t.Errorf("getBuildResources() PidsLimit = %v, want %v", resources.PidsLimit, tt.wantPids)
			}
		})
	}
}

func TestWorker_BuildContextManagement(t *testing.T) {
	// Test build context initialization and cleanup
	w := &Worker{
		BuildContexts:      nil,
		BuildContextsMutex: sync.RWMutex{},
		Config: &Config{
			Build: &Build{
				CPUQuota:    1200,
				MemoryLimit: 4,
				PidsLimit:   1024,
			},
		},
	}

	// Test BuildContexts initialization
	if w.BuildContexts == nil {
		w.BuildContexts = make(map[string]*BuildContext)
	}

	buildID := "test-build-123"
	buildContext := &BuildContext{
		BuildID:       buildID,
		WorkspacePath: "/tmp/vela-build-" + buildID,
		StartTime:     time.Now(),
		Resources:     w.getBuildResources(),
		Environment:   make(map[string]string),
	}

	// Test context storage
	w.BuildContextsMutex.Lock()
	w.BuildContexts[buildID] = buildContext
	w.BuildContextsMutex.Unlock()

	// Verify context is stored
	w.BuildContextsMutex.RLock()
	stored, exists := w.BuildContexts[buildID]
	w.BuildContextsMutex.RUnlock()

	if !exists {
		t.Error("Build context was not stored")
	}

	if stored.BuildID != buildID {
		t.Errorf("Stored build ID = %v, want %v", stored.BuildID, buildID)
	}

	// Test context cleanup
	w.BuildContextsMutex.Lock()
	delete(w.BuildContexts, buildID)
	w.BuildContextsMutex.Unlock()

	// Verify context is cleaned up
	w.BuildContextsMutex.RLock()
	_, exists = w.BuildContexts[buildID]
	w.BuildContextsMutex.RUnlock()

	if exists {
		t.Error("Build context was not cleaned up")
	}
}

func TestBuildContext(t *testing.T) {
	buildID := "test-build-456"
	workspace := "/tmp/vela-build-" + buildID
	startTime := time.Now()

	resources := &BuildResources{
		CPUQuota:  1200,
		Memory:    4 * 1024 * 1024 * 1024,
		PidsLimit: 1024,
	}

	env := make(map[string]string)
	env["TEST_VAR"] = "test_value"

	context := &BuildContext{
		BuildID:       buildID,
		WorkspacePath: workspace,
		StartTime:     startTime,
		Resources:     resources,
		Environment:   env,
	}

	// Test all fields are set correctly
	if context.BuildID != buildID {
		t.Errorf("BuildContext.BuildID = %v, want %v", context.BuildID, buildID)
	}
	if context.WorkspacePath != workspace {
		t.Errorf("BuildContext.WorkspacePath = %v, want %v", context.WorkspacePath, workspace)
	}
	if context.Resources.CPUQuota != 1200 {
		t.Errorf("BuildContext.Resources.CPUQuota = %v, want 1200", context.Resources.CPUQuota)
	}
	if context.Environment["TEST_VAR"] != "test_value" {
		t.Errorf("BuildContext.Environment[TEST_VAR] = %v, want test_value", context.Environment["TEST_VAR"])
	}
}

func TestBuildResources(t *testing.T) {
	resources := &BuildResources{
		CPUQuota:  2000,
		Memory:    8 * 1024 * 1024 * 1024,
		PidsLimit: 2048,
	}

	if resources.CPUQuota != 2000 {
		t.Errorf("BuildResources.CPUQuota = %v, want 2000", resources.CPUQuota)
	}
	if resources.Memory != 8*1024*1024*1024 {
		t.Errorf("BuildResources.Memory = %v, want %v", resources.Memory, 8*1024*1024*1024)
	}
	if resources.PidsLimit != 2048 {
		t.Errorf("BuildResources.PidsLimit = %v, want 2048", resources.PidsLimit)
	}
}
