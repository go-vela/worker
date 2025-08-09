// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"testing"

	"github.com/docker/go-units"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
)

func TestDocker_hostConfig(t *testing.T) {
	// setup logger
	logger := logrus.NewEntry(logrus.StandardLogger())

	tests := []struct {
		name           string
		id             string
		ulimits        pipeline.UlimitSlice
		volumes        []string
		dropCaps       []string
		resourceLimits *ResourceLimits
		wantMemory     int64
		wantCPUQuota   int64
		wantPidsLimit  int64
		wantCapDrop    []string
		wantCapAdd     []string
	}{
		{
			name:          "with resource limits",
			id:            "test-build-1",
			ulimits:       pipeline.UlimitSlice{},
			volumes:       []string{},
			dropCaps:      []string{},
			resourceLimits: &ResourceLimits{
				Memory:    int64(2) * 1024 * 1024 * 1024,
				CPUQuota:  int64(1500),
				CPUPeriod: 100000,
				PidsLimit: 512,
			},
			wantMemory:    int64(2) * 1024 * 1024 * 1024,
			wantCPUQuota:  1500,
			wantPidsLimit: 512,
			wantCapDrop:   []string{"ALL"},
			wantCapAdd:    []string{"CHOWN", "SETUID", "SETGID"},
		},
		{
			name:           "without resource limits (defaults)",
			id:             "test-build-2",
			ulimits:        pipeline.UlimitSlice{},
			volumes:        []string{},
			dropCaps:       []string{},
			resourceLimits: nil,
			wantMemory:     int64(4) * 1024 * 1024 * 1024,
			wantCPUQuota:   int64(1.2 * 100000),
			wantPidsLimit:  1024,
			wantCapDrop:    []string{"ALL"},
			wantCapAdd:     []string{"CHOWN", "SETUID", "SETGID"},
		},
		{
			name: "with custom ulimits",
			id:   "test-build-3",
			ulimits: pipeline.UlimitSlice{
				{
					Name: "nofile",
					Hard: 2048,
					Soft: 2048,
				},
			},
			volumes:        []string{},
			dropCaps:       []string{"NET_ADMIN", "SYS_ADMIN"},
			resourceLimits: nil,
			wantMemory:     int64(4) * 1024 * 1024 * 1024,
			wantCPUQuota:   int64(1.2 * 100000),
			wantPidsLimit:  1024,
			wantCapDrop:    []string{"NET_ADMIN", "SYS_ADMIN"},
			wantCapAdd:     []string{"CHOWN", "SETUID", "SETGID"},
		},
		{
			name:     "with volumes",
			id:       "test-build-4",
			ulimits:  pipeline.UlimitSlice{},
			volumes:  []string{"/host/path:/container/path:ro"},
			dropCaps: []string{},
			resourceLimits: &ResourceLimits{
				Memory:    int64(8) * 1024 * 1024 * 1024,
				CPUQuota:  int64(2000),
				CPUPeriod: 100000,
				PidsLimit: 2048,
			},
			wantMemory:    int64(8) * 1024 * 1024 * 1024,
			wantCPUQuota:  2000,
			wantPidsLimit: 2048,
			wantCapDrop:   []string{"ALL"},
			wantCapAdd:    []string{"CHOWN", "SETUID", "SETGID"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := hostConfig(logger, tt.id, tt.ulimits, tt.volumes, tt.dropCaps, tt.resourceLimits)

			// Check resource limits
			if config.Resources.Memory != tt.wantMemory {
				t.Errorf("hostConfig() Memory = %v, want %v", config.Resources.Memory, tt.wantMemory)
			}
			if config.Resources.CPUQuota != tt.wantCPUQuota {
				t.Errorf("hostConfig() CPUQuota = %v, want %v", config.Resources.CPUQuota, tt.wantCPUQuota)
			}
			if config.Resources.PidsLimit != nil && *config.Resources.PidsLimit != tt.wantPidsLimit {
				t.Errorf("hostConfig() PidsLimit = %v, want %v", *config.Resources.PidsLimit, tt.wantPidsLimit)
			}

			// Check capabilities
			if len(config.CapDrop) != len(tt.wantCapDrop) {
				t.Errorf("hostConfig() CapDrop length = %v, want %v", len(config.CapDrop), len(tt.wantCapDrop))
			}
			if len(config.CapAdd) != len(tt.wantCapAdd) {
				t.Errorf("hostConfig() CapAdd length = %v, want %v", len(config.CapAdd), len(tt.wantCapAdd))
			}

			// Check security options are set
			if len(config.SecurityOpt) != 2 {
				t.Errorf("hostConfig() SecurityOpt length = %v, want 2", len(config.SecurityOpt))
			}

			// Check default mount is created
			if len(config.Mounts) < 1 {
				t.Errorf("hostConfig() should have at least one mount")
			}
			if config.Mounts[0].Target != constants.WorkspaceMount {
				t.Errorf("hostConfig() first mount target = %v, want %v", config.Mounts[0].Target, constants.WorkspaceMount)
			}

			// Check ulimits are applied
			if len(tt.ulimits) > 0 {
				if len(config.Resources.Ulimits) != len(tt.ulimits) {
					t.Errorf("hostConfig() Ulimits length = %v, want %v", len(config.Resources.Ulimits), len(tt.ulimits))
				}
			} else if tt.resourceLimits == nil {
				// Should have default ulimits
				if len(config.Resources.Ulimits) != 2 {
					t.Errorf("hostConfig() should have default ulimits when none provided")
				}
			}
		})
	}
}

func TestResourceLimitsDefaults(t *testing.T) {
	logger := logrus.NewEntry(logrus.StandardLogger())
	
	// Test that nil resource limits apply secure defaults
	config := hostConfig(logger, "test-id", nil, nil, nil, nil)
	
	// Check secure defaults are applied
	if config.Resources.Memory != int64(4)*1024*1024*1024 {
		t.Errorf("Default Memory = %v, want %v", config.Resources.Memory, int64(4)*1024*1024*1024)
	}
	if config.Resources.CPUQuota != int64(1.2*100000) {
		t.Errorf("Default CPUQuota = %v, want %v", config.Resources.CPUQuota, int64(1.2*100000))
	}
	if config.Resources.PidsLimit == nil || *config.Resources.PidsLimit != 1024 {
		t.Errorf("Default PidsLimit not set correctly")
	}
	
	// Check default security ulimits
	foundNofile := false
	foundNproc := false
	for _, ulimit := range config.Resources.Ulimits {
		if ulimit.Name == "nofile" && ulimit.Hard == 1024 && ulimit.Soft == 1024 {
			foundNofile = true
		}
		if ulimit.Name == "nproc" && ulimit.Hard == 512 && ulimit.Soft == 512 {
			foundNproc = true
		}
	}
	if !foundNofile {
		t.Error("Default nofile ulimit not found")
	}
	if !foundNproc {
		t.Error("Default nproc ulimit not found")
	}
	
	// Check security hardening is applied
	if !contains(config.CapDrop, "ALL") {
		t.Error("Should drop ALL capabilities by default")
	}
	if !contains(config.SecurityOpt, "no-new-privileges:true") {
		t.Error("Should have no-new-privileges security option")
	}
	if !contains(config.SecurityOpt, "seccomp=docker/default") {
		t.Error("Should have seccomp security option")
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsUlimit(ulimits []*units.Ulimit, name string, hard, soft int64) bool {
	for _, u := range ulimits {
		if u.Name == name && u.Hard == hard && u.Soft == soft {
			return true
		}
	}
	return false
}