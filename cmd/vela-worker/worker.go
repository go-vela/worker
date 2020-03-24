// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"time"

	"github.com/go-vela/pkg-executor/executor"
	"github.com/go-vela/pkg-queue/queue"
	"github.com/go-vela/pkg-runtime/runtime"
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

	// Server represents the worker configuration for server information.
	Server struct {
		Address string
		Secret  string
	}

	// Worker represents the worker configuration.
	Worker struct {
		API      *API
		Build    *Build
		Executor *executor.Setup
		Queue    *queue.Setup
		Runtime  *runtime.Setup
		Server   *Server
	}
)
