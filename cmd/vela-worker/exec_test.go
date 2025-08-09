// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"
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
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
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