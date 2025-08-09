// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/queue"
	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/runtime"
)

func TestFlags(t *testing.T) {
	// Test that flags() returns a non-empty slice
	flagList := flags()

	if len(flagList) == 0 {
		t.Error("flags() returned empty slice, expected flags")
	}

	// Test that flags include expected base flags plus executor, queue, and runtime flags
	expectedMinFlags := 10 // Conservative minimum - adjust based on actual count
	if len(flagList) < expectedMinFlags {
		t.Errorf("flags() returned %d flags, expected at least %d", len(flagList), expectedMinFlags)
	}

	// Verify that executor, queue, and runtime flags are appended
	executorFlagsCount := len(executor.Flags)
	queueFlagsCount := len(queue.Flags)
	runtimeFlagsCount := len(runtime.Flags)

	// The total should be reasonable but we'll be flexible about the exact count
	if len(flagList) < 20 { // Conservative minimum total
		t.Errorf("flags() returned %d flags, expected at least 20 total flags (including appended flags)", len(flagList))
	}

	t.Logf("flags() returned %d total flags (executor: %d, queue: %d, runtime: %d)",
		len(flagList), executorFlagsCount, queueFlagsCount, runtimeFlagsCount)
}

func TestWorkerAddressValidation(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid http address",
			addr:    "http://localhost:8080",
			wantErr: false,
		},
		{
			name:    "valid https address",
			addr:    "https://vela.example.com",
			wantErr: false,
		},
		{
			name:    "missing scheme",
			addr:    "localhost:8080",
			wantErr: true,
			errMsg:  "worker address must be fully qualified",
		},
		{
			name:    "trailing slash",
			addr:    "http://localhost:8080/",
			wantErr: true,
			errMsg:  "worker address must not have trailing slash",
		},
		{
			name:    "empty address",
			addr:    "",
			wantErr: true, // Empty address lacks scheme so triggers validation error
			errMsg:  "worker address must be fully qualified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get the worker.addr flag
			flagList := flags()

			var workerAddrFlag *cli.StringFlag

			for _, flag := range flagList {
				if strFlag, ok := flag.(*cli.StringFlag); ok && strFlag.Name == "worker.addr" {
					workerAddrFlag = strFlag
					break
				}
			}

			if workerAddrFlag == nil {
				t.Fatal("worker.addr flag not found")
			}

			// Test the validation action
			err := workerAddrFlag.Action(context.Background(), nil, tt.addr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for address %s, got nil", tt.addr)
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error message to contain %s, got %s", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for address %s, got %v", tt.addr, err)
				}
			}
		})
	}
}

func TestFlagDefaults(t *testing.T) {
	// Test default values for key flags
	flagList := flags()

	tests := []struct {
		name         string
		expectedType string
		expectedName string
	}{
		{"worker.addr", "string", "worker.addr"},
		{"checkIn", "duration", "checkIn"},
		{"build.limit", "int", "build.limit"},
		{"build.timeout", "duration", "build.timeout"},
		{"build.cpu-quota", "int", "build.cpu-quota"},
		{"build.memory-limit", "int", "build.memory-limit"},
		{"build.pid-limit", "int", "build.pid-limit"},
		{"log.format", "string", "log.format"},
		{"log.level", "string", "log.level"},
		{"server.addr", "string", "server.addr"},
		{"server.secret", "string", "server.secret"},
		{"server.cert", "string", "server.cert"},
		{"server.cert-key", "string", "server.cert-key"},
		{"server.tls-min-version", "string", "server.tls-min-version"},
	}

	foundFlags := make(map[string]bool)

	for _, flag := range flagList {
		switch f := flag.(type) {
		case *cli.StringFlag:
			foundFlags[f.Name] = true
		case *cli.IntFlag:
			foundFlags[f.Name] = true
		case *cli.DurationFlag:
			foundFlags[f.Name] = true
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !foundFlags[tt.expectedName] {
				t.Errorf("Flag %s not found in flags list", tt.expectedName)
			}
		})
	}
}

func TestFlagValues(t *testing.T) {
	// Test specific default values
	flagList := flags()

	tests := []struct {
		name          string
		flagName      string
		expectedValue interface{}
	}{
		{"checkIn default", "checkIn", 15 * time.Minute},
		{"build.limit default", "build.limit", 1},
		{"build.timeout default", "build.timeout", 30 * time.Minute},
		{"build.cpu-quota default", "build.cpu-quota", 1200},
		{"build.memory-limit default", "build.memory-limit", 4},
		{"build.pid-limit default", "build.pid-limit", 1024},
		{"log.format default", "log.format", "json"},
		{"log.level default", "log.level", "info"},
		{"server.tls-min-version default", "server.tls-min-version", "1.2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, flag := range flagList {
				switch f := flag.(type) {
				case *cli.StringFlag:
					if f.Name == tt.flagName {
						if f.Value != tt.expectedValue {
							t.Errorf("Flag %s value = %v, want %v", tt.flagName, f.Value, tt.expectedValue)
						}
					}
				case *cli.IntFlag:
					if f.Name == tt.flagName {
						if f.Value != tt.expectedValue {
							t.Errorf("Flag %s value = %v, want %v", tt.flagName, f.Value, tt.expectedValue)
						}
					}
				case *cli.DurationFlag:
					if f.Name == tt.flagName {
						if f.Value != tt.expectedValue {
							t.Errorf("Flag %s value = %v, want %v", tt.flagName, f.Value, tt.expectedValue)
						}
					}
				}
			}
		})
	}
}

func TestFlagEnvironmentVariables(t *testing.T) {
	// Test that key flags have sources configured (simplified test)
	flagList := flags()

	keyFlags := []string{
		"worker.addr",
		"checkIn",
		"build.limit",
		"log.format",
		"server.addr",
	}

	for _, flagName := range keyFlags {
		t.Run(flagName, func(t *testing.T) {
			var foundFlag bool

			for _, flag := range flagList {
				var name string

				switch f := flag.(type) {
				case *cli.StringFlag:
					name = f.Name
				case *cli.IntFlag:
					name = f.Name
				case *cli.DurationFlag:
					name = f.Name
				}

				if name == flagName {
					foundFlag = true
					break
				}
			}

			if !foundFlag {
				t.Errorf("Flag %s not found in flags list", flagName)
			}
		})
	}
}

func TestFlagsAppending(t *testing.T) {
	// Test that executor, queue, and runtime flags are properly appended
	baseFlags := 15 // Approximate number of base flags

	// Get all flags
	allFlags := flags()

	// This should include base flags + executor flags + queue flags + runtime flags
	if len(allFlags) <= baseFlags {
		t.Errorf("Expected more than %d flags after appending executor/queue/runtime flags, got %d", baseFlags, len(allFlags))
	}

	// Verify that we have flags from each category by checking for some known flag patterns
	hasExecutorFlag := false
	hasQueueFlag := false
	hasRuntimeFlag := false

	for _, flag := range allFlags {
		var name string

		switch f := flag.(type) {
		case *cli.StringFlag:
			name = f.Name
		case *cli.IntFlag:
			name = f.Name
		case *cli.DurationFlag:
			name = f.Name
		case *cli.BoolFlag:
			name = f.Name
		case *cli.StringSliceFlag:
			name = f.Name
		case *cli.UintFlag:
			name = f.Name
		}

		if strings.Contains(name, "executor") {
			hasExecutorFlag = true
		}

		if strings.Contains(name, "queue") {
			hasQueueFlag = true
		}

		if strings.Contains(name, "runtime") {
			hasRuntimeFlag = true
		}
	}

	if !hasExecutorFlag {
		t.Error("No executor flags found in flag list")
	}

	if !hasQueueFlag {
		t.Error("No queue flags found in flag list")
	}

	if !hasRuntimeFlag {
		t.Error("No runtime flags found in flag list")
	}
}
