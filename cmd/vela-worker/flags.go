// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"time"

	"github.com/go-vela/pkg-executor/executor"
	"github.com/go-vela/pkg-queue/queue"
	"github.com/go-vela/pkg-runtime/runtime"

	"github.com/urfave/cli"
)

func flags() []cli.Flag {
	f := []cli.Flag{

		cli.StringFlag{
			EnvVar: "WORKER_LOG_LEVEL,VELA_LOG_LEVEL,LOG_LEVEL",
			Name:   "log.level",
			Usage:  "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:  "info",
		},

		// API Flags

		cli.StringFlag{
			EnvVar: "WORKER_API_PORT,VELA_API_PORT,API_PORT",
			Name:   "api.port",
			Usage:  "API port for the worker to listen on",
			Value:  ":8080",
		},

		// Build Flags

		cli.IntFlag{
			EnvVar: "WORKER_BUILD_LIMIT,VELA_BUILD_LIMIT,BUILD_LIMIT",
			Name:   "build.limit",
			Usage:  "maximum amount of builds that can run concurrently",
			Value:  1,
		},
		cli.DurationFlag{
			EnvVar: "WORKER_BUILD_TIMEOUT,VELA_BUILD_TIMEOUT,BUILD_TIMEOUT",
			Name:   "build.timeout",
			Usage:  "maximum amount of time a build can run for",
			Value:  30 * time.Minute,
		},

		// Server Flags

		cli.StringFlag{
			EnvVar: "WORKER_SERVER_ADDR,VELA_SERVER_ADDR,VELA_SERVER,SERVER_ADDR",
			Name:   "server.addr",
			Usage:  "Vela server address as a fully qualified url (<scheme>://<host>)",
		},
		cli.StringFlag{
			EnvVar: "WORKER_SERVER_SECRET,VELA_SERVER_SECRET,SERVER_SECRET",
			Name:   "server.secret",
			Usage:  "secret used for server <-> worker communication",
		},
	}

	// Executor Flags

	f = append(f, executor.Flags...)

	// Queue Flags

	f = append(f, queue.Flags...)

	// Runtime Flags

	f = append(f, runtime.Flags...)

	return f
}
