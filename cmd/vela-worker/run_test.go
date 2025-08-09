// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/url"
	"testing"
	"time"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/runtime"
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

func TestWorkerCreation(t *testing.T) {
	// Test Worker struct initialization
	addr, err := url.Parse("http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to parse URL: %v", err)
	}

	outputsCtn := &pipeline.Container{
		Detach:      true,
		Image:       "alpine:latest",
		Environment: make(map[string]string),
		Pull:        constants.PullNotPresent,
	}

	worker := &Worker{
		Config: &Config{
			API: &API{
				Address: addr,
			},
			Build: &Build{
				Limit:       2,
				Timeout:     30 * time.Minute,
				CPUQuota:    1500,
				MemoryLimit: 4,
				PidsLimit:   1024,
			},
			CheckIn: 5 * time.Minute,
			Executor: &executor.Setup{
				Driver:              "linux",
				MaxLogSize:          2097152,
				LogStreamingTimeout: 30 * time.Second,
				EnforceTrustedRepos: false,
				OutputCtn:           outputsCtn,
			},
			Logger: &Logger{
				Format: "json",
				Level:  "info",
			},
			Runtime: &runtime.Setup{
				Driver:           "docker",
				ConfigFile:       "",
				Namespace:        "vela",
				PodsTemplateName: "",
				PodsTemplateFile: "",
				HostVolumes:      []string{},
				PrivilegedImages: []string{"alpine"},
				DropCapabilities: []string{"NET_RAW"},
			},
			Queue: &queue.Setup{
				Address: "redis://localhost:6379",
				Driver:  "redis",
				Cluster: false,
				Routes:  []string{"vela"},
				Timeout: 30 * time.Second,
			},
			Server: &Server{
				Address: "http://localhost:8080",
				Secret:  "test-secret",
			},
			Certificate: &Certificate{
				Cert: "",
				Key:  "",
			},
			TLSMinVersion: "1.2",
		},
		Executors:     make(map[int]executor.Engine),
		RegisterToken: make(chan string, 1),
		RunningBuilds: make([]*api.Build, 0),
	}

	// Test that the worker structure is properly initialized
	if worker.Config.API.Address.String() != "http://localhost:8080" {
		t.Errorf("Worker API Address = %v, want http://localhost:8080", worker.Config.API.Address.String())
	}

	if worker.Config.Build.Limit != 2 {
		t.Errorf("Worker Build Limit = %v, want 2", worker.Config.Build.Limit)
	}

	if worker.Config.Build.CPUQuota != 1500 {
		t.Errorf("Worker Build CPUQuota = %v, want 1500", worker.Config.Build.CPUQuota)
	}

	if worker.Config.Executor.Driver != "linux" {
		t.Errorf("Worker Executor Driver = %v, want linux", worker.Config.Executor.Driver)
	}

	if worker.Config.Runtime.Driver != "docker" {
		t.Errorf("Worker Runtime Driver = %v, want docker", worker.Config.Runtime.Driver)
	}

	if worker.Config.Queue.Driver != "redis" {
		t.Errorf("Worker Queue Driver = %v, want redis", worker.Config.Queue.Driver)
	}

	if len(worker.Executors) != 0 {
		t.Errorf("Worker Executors length = %v, want 0", len(worker.Executors))
	}

	if len(worker.RunningBuilds) != 0 {
		t.Errorf("Worker RunningBuilds length = %v, want 0", len(worker.RunningBuilds))
	}

	if cap(worker.RegisterToken) != 1 {
		t.Errorf("Worker RegisterToken capacity = %v, want 1", cap(worker.RegisterToken))
	}
}

func TestAPI_Configuration(t *testing.T) {
	api := &API{
		Address: mustParseURL("https://vela.example.com:8443"),
	}

	if api.Address.Scheme != "https" {
		t.Errorf("API Address Scheme = %v, want https", api.Address.Scheme)
	}

	if api.Address.Host != "vela.example.com:8443" {
		t.Errorf("API Address Host = %v, want vela.example.com:8443", api.Address.Host)
	}
}

func TestServer_Configuration(t *testing.T) {
	server := &Server{
		Address: "https://api.vela.example.com",
		Secret:  "super-secret-key",
	}

	if server.Address != "https://api.vela.example.com" {
		t.Errorf("Server Address = %v, want https://api.vela.example.com", server.Address)
	}

	if server.Secret != "super-secret-key" {
		t.Errorf("Server Secret = %v, want super-secret-key", server.Secret)
	}
}

func TestCertificate_Configuration(t *testing.T) {
	cert := &Certificate{
		Cert: "/path/to/cert.pem",
		Key:  "/path/to/key.pem",
	}

	if cert.Cert != "/path/to/cert.pem" {
		t.Errorf("Certificate Cert = %v, want /path/to/cert.pem", cert.Cert)
	}

	if cert.Key != "/path/to/key.pem" {
		t.Errorf("Certificate Key = %v, want /path/to/key.pem", cert.Key)
	}
}

func TestLogger_Configuration(t *testing.T) {
	logger := &Logger{
		Format: "json",
		Level:  "debug",
	}

	if logger.Format != "json" {
		t.Errorf("Logger Format = %v, want json", logger.Format)
	}

	if logger.Level != "debug" {
		t.Errorf("Logger Level = %v, want debug", logger.Level)
	}
}

// Helper function for tests
func mustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}
