// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/queue"
	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/runtime"
)

// flags is a helper function to return the all
// supported command line interface (CLI) flags
// for the Worker.
func flags() []cli.Flag {
	f := []cli.Flag{

		&cli.StringFlag{
			Name:    "worker.addr",
			Usage:   "Worker server address as a fully qualified url (<scheme>://<host>)",
			Sources: cli.EnvVars("WORKER_ADDR", "VELA_WORKER_ADDR", "VELA_WORKER"),
			Action: func(_ context.Context, _ *cli.Command, v string) error {
				// check if the worker address has a scheme
				if !strings.Contains(v, "://") {
					return fmt.Errorf("worker address must be fully qualified (<scheme>://<host>)")
				}

				// check if the worker address has a trailing slash
				if strings.HasSuffix(v, "/") {
					return fmt.Errorf("worker address must not have trailing slash")
				}

				return nil
			},
		},

		&cli.DurationFlag{
			Name:    "checkIn",
			Usage:   "time to wait in between checking in with the server",
			Sources: cli.EnvVars("WORKER_CHECK_IN", "VELA_CHECK_IN", "CHECK_IN"),
			Value:   15 * time.Minute,
		},

		// Build Flags

		&cli.IntFlag{
			Name:    "build.limit",
			Usage:   "maximum amount of builds that can run concurrently",
			Sources: cli.EnvVars("WORKER_BUILD_LIMIT", "VELA_BUILD_LIMIT", "BUILD_LIMIT"),
			Value:   1,
		},
		&cli.DurationFlag{
			Name:    "build.timeout",
			Usage:   "maximum amount of time a build can run for",
			Sources: cli.EnvVars("WORKER_BUILD_TIMEOUT", "VELA_BUILD_TIMEOUT", "BUILD_TIMEOUT"),
			Value:   30 * time.Minute,
		},
		&cli.IntFlag{
			Name:    "build.cpu-quota",
			Usage:   "CPU quota per build in millicores (1000 = 1 core)",
			Value:   1200, // 1.2 CPU cores per build
			Sources: cli.EnvVars("VELA_BUILD_CPU_QUOTA", "BUILD_CPU_QUOTA"),
		},
		&cli.IntFlag{
			Name:    "build.memory-limit",
			Usage:   "Memory limit per build in GB",
			Value:   4, // 4GB per build
			Sources: cli.EnvVars("VELA_BUILD_MEMORY_LIMIT", "BUILD_MEMORY_LIMIT"),
		},
		&cli.IntFlag{
			Name:    "build.pid-limit",
			Usage:   "Process limit per build container",
			Value:   1024, // Prevent fork bombs
			Sources: cli.EnvVars("VELA_BUILD_PID_LIMIT", "BUILD_PID_LIMIT"),
		},

		// Logger Flags

		&cli.StringFlag{
			Name:    "log.format",
			Usage:   "set log format for the worker",
			Sources: cli.EnvVars("WORKER_LOG_FORMAT", "VELA_LOG_FORMAT", "LOG_FORMAT"),
			Value:   "json",
		},
		&cli.StringFlag{
			Sources: cli.EnvVars("WORKER_LOG_LEVEL", "VELA_LOG_LEVEL", "LOG_LEVEL"),
			Name:    "log.level",
			Usage:   "set log level for the worker",
			Value:   "info",
		},

		// Server Flags

		&cli.StringFlag{
			Name:    "server.addr",
			Usage:   "Vela server address as a fully qualified url (<scheme>://<host>)",
			Sources: cli.EnvVars("WORKER_SERVER_ADDR", "VELA_SERVER_ADDR", "VELA_SERVER", "SERVER_ADDR"),
		},
		&cli.StringFlag{
			Name:    "server.secret",
			Usage:   "secret used for server <-> worker communication",
			Sources: cli.EnvVars("WORKER_SERVER_SECRET", "VELA_SERVER_SECRET", "SERVER_SECRET"),
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "server.cert",
			Usage:   "optional TLS certificate for https",
			Sources: cli.EnvVars("WORKER_SERVER_CERT", "VELA_SERVER_CERT", "SERVER_CERT"),
		},
		&cli.StringFlag{
			Name:    "server.cert-key",
			Usage:   "optional TLS certificate key",
			Sources: cli.EnvVars("WORKER_SERVER_CERT_KEY", "VELA_SERVER_CERT_KEY", "SERVER_CERT_KEY"),
		},
		&cli.StringFlag{
			Name:    "server.tls-min-version",
			Usage:   "optional TLS minimum version requirement",
			Sources: cli.EnvVars("WORKER_SERVER_TLS_MIN_VERSION", "VELA_SERVER_TLS_MIN_VERSION", "SERVER_TLS_MIN_VERSION"),
			Value:   "1.2",
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
