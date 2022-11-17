// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"net/url"
	"sync"
	"time"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime"
)

type (
	// API represents the worker configuration for API information.
	API struct {
		Address *url.URL
	}

	// Build represents the worker configuration for build information.
	Build struct {
		Limit   int
		Timeout time.Duration
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

	// Worker represents all configuration and
	// system processes for the worker.
	Worker struct {
		Config         *Config
		Executors      map[int]executor.Engine
		Queue          queue.Service
		Runtime        runtime.Engine
		VelaClient     *vela.Client
		PackageChannel chan *types.BuildPackage
		*message.Activity
	}
)

func ToLibrary(w *Worker) *library.Worker {
	_w := new(library.Worker)
	_w.SetHostname(w.Config.API.Address.Hostname())
	_w.SetAddress(w.Config.API.Address.String())
	_w.SetRoutes(w.Config.Queue.Routes)
	_w.SetBuildLimit(int64(w.Config.Build.Limit))
	return _w
}

func NewActivity() *message.Activity {
	return &message.Activity{
		ActiveBuilds: []*library.Build{},
		Mutex:        new(sync.Mutex),
		Channel:      make(chan message.BuildActivity),
	}
}

func ToWorkerStatus(w *Worker, a *message.Activity) string {
	if len(a.ActiveBuilds) < w.Config.Build.Limit {
		return "ready"
	} else if len(a.ActiveBuilds) > 0 {
		return "building"
	} else {
		return "error"
	}
}
