// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"net/http"
	"time"

	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/router"
	"github.com/go-vela/worker/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"
	tomb "gopkg.in/tomb.v2"
)

func server(c *cli.Context) error {
	// validate all input
	err := validate(c)
	if err != nil {
		return err
	}

	// set log level for logrus
	switch c.String("log-level") {
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

	// create a vela client
	vela, err := setupClient(c)
	if err != nil {
		return err
	}

	// create a runtime client
	runtime, err := setupRuntime(c)
	if err != nil {
		return err
	}

	// create a queue client
	queue, err := setupQueue(c)
	if err != nil {
		return err
	}

	// create the executor clients
	executors := make(map[int]executor.Engine)
	for i := 0; i < c.Int("executor-threads"); i++ {
		executor, err := setupExecutor(c, vela, runtime)
		if err != nil {
			return err
		}

		executors[i] = executor
	}

	router := router.Load(
		middleware.RequestVersion,
		// TODO: middleware.Executor(executors),
		middleware.Logger(logrus.StandardLogger(), time.RFC3339, true),
	)

	tomb := new(tomb.Tomb)
	tomb.Go(func() error {
		// Start server
		srv := &http.Server{Addr: c.String("server-port"), Handler: router}

		go func() {
			logrus.Info("Starting HTTP server...")
			err := srv.ListenAndServe()
			if err != nil {
				tomb.Kill(err)
			}
		}()

		go func() {
			logrus.Info("Starting operator...")
			// TODO: refactor due to one thread killing entire worker
			err := operate(queue, executors, c.Duration("executor-timeout"))
			if err != nil {
				tomb.Kill(err)
			}
		}()

		for {
			select {
			case <-tomb.Dying():
				logrus.Info("Stopping HTTP server...")
				return srv.Shutdown(nil)
			}
		}
	})

	// Wait for stuff and watch for errors
	tomb.Wait()
	return tomb.Err()
}
