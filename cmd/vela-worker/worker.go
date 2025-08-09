// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/url"
	"sync"
	"time"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/runtime"
)

type (
	// API represents the worker configuration for API information.
	API struct {
		Address *url.URL
	}

	// Build represents the worker configuration for build information.
	Build struct {
		Limit       int
		Timeout     time.Duration
		CPUQuota    int // CPU quota per build in millicores
		MemoryLimit int // Memory limit per build in GB
		PidsLimit   int // Process limit per build container
	}

	// Logger represents the worker configuration for logger information.
	Logger struct {
		Format string
		Level  string
	}

	// Server represents the worker configuration for server information.
	Server struct {
		Address string
		Secret  string
	}

	// Certificate represents the optional cert and key to enable TLS.
	Certificate struct {
		Cert string
		Key  string
	}

	// Config represents the worker configuration.
	Config struct {
		Mock          bool // Mock should only be true for tests
		API           *API
		Build         *Build
		CheckIn       time.Duration
		Executor      *executor.Setup
		Logger        *Logger
		Queue         *queue.Setup
		Runtime       *runtime.Setup
		Server        *Server
		Certificate   *Certificate
		TLSMinVersion string
	}

	// BuildContext represents isolated build execution context.
	BuildContext struct {
		BuildID       string            // Cryptographic ID for build isolation
		WorkspacePath string            // Isolated workspace path
		StartTime     time.Time         // Build start time
		Resources     *BuildResources   // Resource allocation
		Environment   map[string]string // Environment variables
	}

	// BuildResources represents resource limits for a build.
	BuildResources struct {
		CPUQuota  int64 // CPU limit in millicores (1000 = 1 core)
		Memory    int64 // Memory in bytes
		PidsLimit int64 // Process limit
	}

	// Worker represents all configuration and
	// system processes for the worker.
	Worker struct {
		Config             *Config
		Executors          map[int]executor.Engine
		Queue              queue.Service
		Runtime            runtime.Engine
		VelaClient         *vela.Client
		RegisterToken      chan string
		CheckedIn          bool
		RunningBuilds      []*api.Build
		QueueCheckedIn     bool
		RunningBuildsMutex sync.Mutex
		// Security-focused build tracking (works with single builds, scales to concurrent)
		BuildContexts      map[string]*BuildContext // Thread-safe build context tracking
		BuildContextsMutex sync.RWMutex             // Thread-safe access to build contexts
	}
)
