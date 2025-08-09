// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"
)

func TestBuild_ResourceConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		limit       int
		cpuQuota    int
		memoryLimit int
		pidsLimit   int
		wantLimit   int
		wantCPU     int
		wantMemory  int
		wantPids    int
	}{
		{
			name:        "default configuration",
			limit:       1,
			cpuQuota:    1200,
			memoryLimit: 4,
			pidsLimit:   1024,
			wantLimit:   1,
			wantCPU:     1200,
			wantMemory:  4,
			wantPids:    1024,
		},
		{
			name:        "high resource configuration",
			limit:       4,
			cpuQuota:    2000,
			memoryLimit: 8,
			pidsLimit:   2048,
			wantLimit:   4,
			wantCPU:     2000,
			wantMemory:  8,
			wantPids:    2048,
		},
		{
			name:        "minimal configuration",
			limit:       1,
			cpuQuota:    500,
			memoryLimit: 1,
			pidsLimit:   256,
			wantLimit:   1,
			wantCPU:     500,
			wantMemory:  1,
			wantPids:    256,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			build := &Build{
				Limit:       tt.limit,
				CPUQuota:    tt.cpuQuota,
				MemoryLimit: tt.memoryLimit,
				PidsLimit:   tt.pidsLimit,
			}

			// Test that values are set correctly
			if build.Limit != tt.wantLimit {
				t.Errorf("Build.Limit = %v, want %v", build.Limit, tt.wantLimit)
			}

			if build.CPUQuota != tt.wantCPU {
				t.Errorf("Build.CPUQuota = %v, want %v", build.CPUQuota, tt.wantCPU)
			}

			if build.MemoryLimit != tt.wantMemory {
				t.Errorf("Build.MemoryLimit = %v, want %v", build.MemoryLimit, tt.wantMemory)
			}

			if build.PidsLimit != tt.wantPids {
				t.Errorf("Build.PidsLimit = %v, want %v", build.PidsLimit, tt.wantPids)
			}
		})
	}
}

func TestConfig_SecurityConfiguration(t *testing.T) {
	// Test that Config struct properly holds build configuration
	config := &Config{
		Build: &Build{
			Limit:       2,
			CPUQuota:    1500,
			MemoryLimit: 6,
			PidsLimit:   1536,
		},
	}

	if config.Build.Limit != 2 {
		t.Errorf("Config.Build.Limit = %v, want 2", config.Build.Limit)
	}

	if config.Build.CPUQuota != 1500 {
		t.Errorf("Config.Build.CPUQuota = %v, want 1500", config.Build.CPUQuota)
	}

	if config.Build.MemoryLimit != 6 {
		t.Errorf("Config.Build.MemoryLimit = %v, want 6", config.Build.MemoryLimit)
	}

	if config.Build.PidsLimit != 1536 {
		t.Errorf("Config.Build.PidsLimit = %v, want 1536", config.Build.PidsLimit)
	}
}

func TestWorker_ConfigurationIntegration(t *testing.T) {
	// Test that Worker properly integrates with Config and Build
	worker := &Worker{
		Config: &Config{
			Build: &Build{
				Limit:       3,
				CPUQuota:    1800,
				MemoryLimit: 8,
				PidsLimit:   2048,
			},
		},
	}

	// Test getBuildResources integration
	resources := worker.getBuildResources()

	expectedMemory := int64(8) * 1024 * 1024 * 1024
	if resources.Memory != expectedMemory {
		t.Errorf("getBuildResources().Memory = %v, want %v", resources.Memory, expectedMemory)
	}

	if resources.CPUQuota != 1800 {
		t.Errorf("getBuildResources().CPUQuota = %v, want 1800", resources.CPUQuota)
	}

	if resources.PidsLimit != 2048 {
		t.Errorf("getBuildResources().PidsLimit = %v, want 2048", resources.PidsLimit)
	}
}
