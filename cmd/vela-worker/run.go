// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/queue"
	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/runtime"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

// run executes the worker based
// off the configuration provided.
func run(c *cli.Context) error {
	// set log format for the worker
	switch c.String("log.format") {
	case "t", "text", "Text", "TEXT":
		logrus.SetFormatter(&logrus.TextFormatter{})
	case "j", "json", "Json", "JSON":
		fallthrough
	default:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	// set log level for the worker
	switch c.String("log.level") {
	case "t", "trace", "Trace", "TRACE":
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.DebugLevel)
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
	case "i", "info", "Info", "INFO":
		fallthrough
	default:
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.InfoLevel)
	}

	// create a log entry with extra metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#WithFields
	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/worker/",
		"docs":     "https://go-vela.github.io/docs/concepts/infrastructure/worker/",
		"registry": "https://hub.docker.com/r/target/vela-worker/",
	}).Info("Vela Worker")

	// parse the workers address, returning any errors.
	addr, err := url.Parse(c.String("worker.addr"))
	if err != nil {
		return fmt.Errorf("unable to parse address: %w", err)
	}

	// create the worker
	w := &Worker{
		// worker configuration
		Config: &Config{
			// api configuration
			API: &API{
				Address: addr,
			},
			// build configuration
			Build: &Build{
				Limit:   c.Int("build.limit"),
				Timeout: c.Duration("build.timeout"),
			},
			// build configuration
			CheckIn: c.Duration("checkIn"),
			// executor configuration
			Executor: &executor.Setup{
				Driver:              c.String("executor.driver"),
				LogMethod:           c.String("executor.log_method"),
				MaxLogSize:          c.Uint("executor.max_log_size"),
				LogStreamingTimeout: c.Duration("executor.log_streaming_timeout"),
				EnforceTrustedRepos: c.Bool("executor.enforce-trusted-repos"),
			},
			// logger configuration
			Logger: &Logger{
				Format: c.String("log.format"),
				Level:  c.String("log.level"),
			},
			// runtime configuration
			Runtime: &runtime.Setup{
				Driver:           c.String("runtime.driver"),
				ConfigFile:       c.String("runtime.config"),
				Namespace:        c.String("runtime.namespace"),
				PodsTemplateName: c.String("runtime.pods-template-name"),
				PodsTemplateFile: c.Path("runtime.pods-template-file"),
				HostVolumes:      c.StringSlice("runtime.volumes"),
				PrivilegedImages: c.StringSlice("runtime.privileged-images"),
			},
			// queue configuration
			Queue: &queue.Setup{
				Driver:  c.String("queue.driver"),
				Address: c.String("queue.addr"),
				Cluster: c.Bool("queue.cluster"),
				Routes:  c.StringSlice("queue.routes"),
				Timeout: c.Duration("queue.pop.timeout"),
			},
			// server configuration
			Server: &Server{
				Address: c.String("server.addr"),
				Secret:  c.String("server.secret"),
			},
			// Certificate configuration
			Certificate: &Certificate{
				Cert: c.String("server.cert"),
				Key:  c.String("server.cert-key"),
			},
			// TLS minimum version enforced
			TLSMinVersion: c.String("server.tls-min-version"),
		},
		Executors: make(map[int]executor.Engine),

		AuthToken: make(chan string, 1),
	}

	// set the worker address if no flag was provided
	if len(w.Config.API.Address.String()) == 0 {
		w.Config.API.Address, _ = url.Parse(fmt.Sprintf("http://%s", hostname))
	}

	if len(c.String("server.secret")) > 0 {
		w.AuthToken <- c.String("server.secret")
	}

	// validate the worker
	err = w.Validate()
	if err != nil {
		return err
	}

	// start the worker
	return w.Start()
}
