// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"os"
	"time"

	"github.com/go-vela/worker/version"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "vela-worker"
	app.Action = server
	app.Version = version.Version.String()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server-port",
			Usage:  "API port to listen on",
			EnvVar: "VELA_PORT",
			Value:  ":8080",
		},
		cli.StringFlag{
			EnvVar: "VELA_LOG_LEVEL,LOG_LEVEL",
			Name:   "log-level",
			Usage:  "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:  "info",
		},
		cli.StringFlag{
			EnvVar: "VELA_ADDR,VELA_HOST",
			Name:   "server-addr",
			Usage:  "server address as a fully qualified url (<scheme>://<host>)",
		},
		cli.StringFlag{
			EnvVar: "VELA_SECRET",
			Name:   "vela-secret",
			Usage:  "secret used for server <-> worker communication",
		},

		// Executor Flags
		cli.StringFlag{
			EnvVar: "VELA_EXECUTOR_DRIVER,EXECUTOR_DRIVER",
			Name:   "executor-driver",
			Usage:  "executor driver",
			Value:  "linux",
		},
		cli.IntFlag{
			EnvVar: "VELA_EXECUTOR_THREADS,EXECUTOR_THREADS",
			Name:   "executor-threads",
			Usage:  "number of executor threads to create",
			Value:  1,
		},
		cli.DurationFlag{
			EnvVar: "VELA_EXECUTOR_TIMEOUT,EXECUTOR_TIMEOUT",
			Name:   "executor-timeout",
			Usage:  "max time an executor will run a build",
			Value:  60 * time.Minute,
		},

		// Queue Flags
		cli.StringFlag{
			EnvVar: "VELA_QUEUE_DRIVER,QUEUE_DRIVER",
			Name:   "queue-driver",
			Usage:  "queue driver",
		},
		cli.StringFlag{
			EnvVar: "VELA_QUEUE_CONFIG,QUEUE_CONFIG",
			Name:   "queue-config",
			Usage:  "queue driver configuration string",
		},
		cli.BoolFlag{
			EnvVar: "VELA_QUEUE_CLUSTER,QUEUE_CLUSTER",
			Name:   "queue-cluster",
			Usage:  "queue client is setup for clusters",
		},
		// By default all builds are pushed to the "vela" route
		cli.StringSliceFlag{
			EnvVar: "VELA_QUEUE_WORKER_ROUTES,QUEUE_WORKER_ROUTES",
			Name:   "queue-worker-routes",
			Usage:  "queue worker routes is configuration for routing builds",
		},

		// Runtime Flags
		cli.StringFlag{
			EnvVar: "VELA_RUNTIME_DRIVER,RUNTIME_DRIVER",
			Name:   "runtime-driver",
			Usage:  "runtime driver",
		},
	}

	// set logrus to log in JSON format
	log.SetFormatter(&log.JSONFormatter{})

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
