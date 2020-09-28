// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"time"

	"github.com/go-vela/pkg-executor/executor"
	"github.com/go-vela/pkg-queue/queue"
	"github.com/go-vela/pkg-runtime/runtime"
	"github.com/go-vela/sdk-go/vela"
)

type (
	// API represents the worker configuration for API information.
	API struct {
		Port string
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

	// Config represents the worker configuration.
	Config struct {
		API      *API
		Build    *Build
		Executor *executor.Setup
		Hostname string
		Logger   *Logger
		Queue    *queue.Setup
		Runtime  *runtime.Setup
		Server   *Server
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
