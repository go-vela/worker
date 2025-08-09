// SPDX-License-Identifier: Apache-2.0

package version

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		setupTag string
		wantTag  string
	}{
		{
			name:     "with valid semantic version",
			setupTag: "v1.2.3",
			wantTag:  "v1.2.3",
		},
		{
			name:     "with empty tag (default fallback)",
			setupTag: "",
			wantTag:  "v0.0.0",
		},
		{
			name:     "with prerelease version",
			setupTag: "v1.0.0-alpha.1",
			wantTag:  "v1.0.0-alpha.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			originalTag := Tag
			Tag = tt.setupTag
			defer func() { Tag = originalTag }()

			// Test
			v := New()

			// Verify
			if v == nil {
				t.Error("New() returned nil")
				return
			}

			if v.Canonical != tt.wantTag {
				t.Errorf("New().Canonical = %v, want %v", v.Canonical, tt.wantTag)
			}

			// Verify metadata is populated
			if v.Metadata.Architecture == "" {
				t.Error("Metadata.Architecture should not be empty")
			}
			if v.Metadata.Compiler == "" {
				t.Error("Metadata.Compiler should not be empty")
			}
			if v.Metadata.GoVersion == "" {
				t.Error("Metadata.GoVersion should not be empty")
			}
			if v.Metadata.OperatingSystem == "" {
				t.Error("Metadata.OperatingSystem should not be empty")
			}
		})
	}
}

func TestNew_WithCommitAndDate(t *testing.T) {
	// Setup
	originalTag := Tag
	originalCommit := Commit
	originalDate := Date
	
	Tag = "v1.0.0"
	Commit = "abc123"
	Date = "2023-01-01T00:00:00Z"
	
	defer func() {
		Tag = originalTag
		Commit = originalCommit
		Date = originalDate
	}()

	// Test
	v := New()

	// Verify
	if v.Metadata.GitCommit != "abc123" {
		t.Errorf("Metadata.GitCommit = %v, want abc123", v.Metadata.GitCommit)
	}
	if v.Metadata.BuildDate != "2023-01-01T00:00:00Z" {
		t.Errorf("Metadata.BuildDate = %v, want 2023-01-01T00:00:00Z", v.Metadata.BuildDate)
	}
}

func TestPackageVariables(t *testing.T) {
	// Test that package variables are set correctly
	if Arch == "" {
		t.Error("Arch should not be empty")
	}
	if Compiler == "" {
		t.Error("Compiler should not be empty")
	}
	if Go == "" {
		t.Error("Go should not be empty")
	}
	if OS == "" {
		t.Error("OS should not be empty")
	}
}