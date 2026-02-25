// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	_ "github.com/joho/godotenv/autoload"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/runtime"
)

// run executes the worker based
// off the configuration provided.
func run(ctx context.Context, c *cli.Command) error {
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
	// https://pkg.go.dev/github.com/sirupsen/logrus#WithFields
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

	outputsCtn := new(pipeline.Container)
	if len(c.String("executor.outputs-image")) > 0 {
		outputsCtn = &pipeline.Container{
			Detach:      true,
			Image:       c.String("executor.outputs-image"),
			Environment: make(map[string]string),
			Pull:        constants.PullNotPresent,
		}
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
				Limit:   c.Int32("build.limit"),
				Timeout: c.Duration("build.timeout"),
			},
			// build configuration
			CheckIn: c.Duration("checkIn"),
			// executor configuration
			Executor: &executor.Setup{
				Driver:              c.String("executor.driver"),
				MaxLogSize:          c.Uint("executor.max_log_size"),
				FileSizeLimit:       c.Int("storage.file-size-limit"),
				BuildFileSizeLimit:  c.Int("storage.build-file-size-limit"),
				LogStreamingTimeout: c.Duration("executor.log_streaming_timeout"),
				EnforceTrustedRepos: c.Bool("executor.enforce-trusted-repos"),
				OutputCtn:           outputsCtn,
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
				PodsTemplateFile: c.String("runtime.pods-template-file"),
				HostVolumes:      c.StringSlice("runtime.volumes"),
				PrivilegedImages: c.StringSlice("runtime.privileged-images"),
				DropCapabilities: c.StringSlice("runtime.drop-capabilities"),
			},
			// queue configuration
			Queue: &queue.Setup{
				Address: c.String("queue.addr"),
				Driver:  c.String("queue.driver"),
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

		RegisterToken: make(chan string, 1),

		RunningBuilds: make([]*api.Build, 0),
	}

	// set the worker address if no flag was provided
	if len(w.Config.API.Address.String()) == 0 {
		w.Config.API.Address, _ = url.Parse(fmt.Sprintf("http://%s", hostname))
	}

	// if server secret is provided, use as register token on start up
	if len(c.String("server.secret")) > 0 {
		logrus.Trace("registering worker with embedded server secret")

		w.RegisterToken <- c.String("server.secret")
	}

	// validate the worker
	err = w.Validate()
	if err != nil {
		return err
	}

	// start the worker
	return w.Start(ctx)
}
