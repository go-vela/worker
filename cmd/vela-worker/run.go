// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/gin-gonic/gin"

	"github.com/go-vela/pkg-executor/executor"
	"github.com/go-vela/pkg-queue/queue"
	"github.com/go-vela/pkg-runtime/runtime"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"

	_ "github.com/joho/godotenv/autoload"
)

// run executes the worker based off the configuration provided.
func run(c *cli.Context) error {
	// set log level for the worker
	switch c.String("log.level") {
	case "t", "trace", "Trace", "TRACE":
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.DebugLevel)
	case "i", "info", "Info", "INFO":
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.InfoLevel)
	case "w", "warn", "Warn", "WARN":
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.PanicLevel)
	}

	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/worker/",
		"docs":     "https://go-vela.github.io/docs/concepts/infrastructure/worker/",
		"registry": "https://hub.docker.com/r/target/vela-worker/",
	}).Info("Vela Worker")

	// create the worker
	w := &Worker{
		// worker configuration
		Config: &Config{
			// api configuration
			API: &API{
				Port: c.String("api.port"),
			},
			// build configuration
			Build: &Build{
				Limit:   c.Int("build.limit"),
				Timeout: c.Duration("build.timeout"),
			},
			// executor configuration
			Executor: &executor.Setup{
				Driver: c.String("executor.driver"),
			},
			// runtime configuration
			Runtime: &runtime.Setup{
				Driver:    c.String("runtime.driver"),
				Config:    c.String("runtime.config"),
				Namespace: c.String("runtime.namespace"),
			},
			// queue configuration
			Queue: &queue.Setup{
				Driver:  c.String("queue.driver"),
				Config:  c.String("queue.config"),
				Cluster: c.Bool("queue.cluster"),
				Routes:  c.StringSlice("queue.routes"),
			},
			// server configuration
			Server: &Server{
				Address: c.String("server.addr"),
				Secret:  c.String("server.secret"),
			},
		},
		Executors: make(map[int]executor.Engine),
	}

	// validate the worker
	err := w.Validate()
	if err != nil {
		return err
	}

	// start the worker
	return w.Start()
}
