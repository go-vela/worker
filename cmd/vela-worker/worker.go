// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"net/url"
	"time"

	"github.com/go-vela/sdk-go/vela"
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
		Mock        bool // Mock should only be true for tests
		API         *API
		Build       *Build
		CheckIn     time.Duration
		Executor    *executor.Setup
		Logger      *Logger
		Queue       *queue.Setup
		Runtime     *runtime.Setup
		Server      *Server
		Certificate *Certificate
	}

	// Worker represents all configuration and
	// system processes for the worker.
	Worker struct {
		Config     *Config
		Executors  map[int]executor.Engine
		Queue      queue.Service
		Runtime    runtime.Engine
		VelaClient *vela.Client
	}
)
