// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"context"
	"time"

	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/runtime"
	"github.com/go-vela/worker/version"

	"github.com/sirupsen/logrus"
)

// exec is a helper function to poll the queue
// and execute Vela pipelines for the Worker.
//
// nolint:funlen // ignore function length due to comments and log messages
func (w *Worker) exec(index int) error {
	var err error

	// setup the version
	v := version.New()

	// capture an item from the queue
	item, err := w.Queue.Pop(context.Background())
	if err != nil {
		return err
	}

	if item == nil {
		return nil
	}

	// create logger with extra metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#WithFields
	logger := logrus.WithFields(logrus.Fields{
		"build":    item.Build.GetNumber(),
		"executor": w.Config.Executor.Driver,
		"host":     w.Config.API.Address.Hostname(),
		"repo":     item.Repo.GetFullName(),
		"runtime":  w.Config.Runtime.Driver,
		"user":     item.User.GetName(),
		"version":  v.Semantic(),
	})

	// setup the runtime
	//
	// https://pkg.go.dev/github.com/go-vela/worker/runtime?tab=doc#New
	w.Runtime, err = runtime.New(&runtime.Setup{
		Logger:           logger,
		Driver:           w.Config.Runtime.Driver,
		ConfigFile:       w.Config.Runtime.ConfigFile,
		HostVolumes:      w.Config.Runtime.HostVolumes,
		Namespace:        w.Config.Runtime.Namespace,
		PrivilegedImages: w.Config.Runtime.PrivilegedImages,
	})
	if err != nil {
		return err
	}

	// setup the executor
	//
	// https://godoc.org/github.com/go-vela/worker/executor#New
	_executor, err := executor.New(&executor.Setup{
		Logger:     logger,
		Driver:     w.Config.Executor.Driver,
		LogMethod:  w.Config.Executor.LogMethod,
		MaxLogSize: w.Config.Executor.MaxLogSize,
		Client:     w.VelaClient,
		Hostname:   w.Config.API.Address.Hostname(),
		Runtime:    w.Runtime,
		Build:      item.Build,
		Pipeline:   item.Pipeline.Sanitize(w.Config.Runtime.Driver),
		Repo:       item.Repo,
		User:       item.User,
		Version:    v.Semantic(),
	})

	// add the executor to the worker
	w.Executors[index] = _executor

	// capture the configured build timeout
	t := w.Config.Build.Timeout
	// check if the repository has a custom timeout
	if item.Repo.GetTimeout() > 0 {
		// update timeout variable to repository custom timeout
		t = time.Duration(item.Repo.GetTimeout()) * time.Minute
	}

	// create a background context
	ctx := context.Background()

	// add to the background context with a timeout
	// built in for ensuring a build doesn't run forever
	ctx, timeout := context.WithTimeout(ctx, t)
	defer timeout()

	defer func() {
		logger.Info("destroying build")
		// destroy the build with the executor
		err = _executor.DestroyBuild(context.Background())
		if err != nil {
			logger.Errorf("unable to destroy build: %v", err)
		}
	}()

	logger.Info("creating build")
	// create the build with the executor
	err = _executor.CreateBuild(ctx)
	if err != nil {
		logger.Errorf("unable to create build: %v", err)
		return nil
	}

	logger.Info("planning build")
	// plan the build with the executor
	err = _executor.PlanBuild(ctx)
	if err != nil {
		logger.Errorf("unable to plan build: %v", err)
		return nil
	}

	logger.Info("assembling build")
	// assemble the build with the executor
	err = _executor.AssembleBuild(ctx)
	if err != nil {
		logger.Errorf("unable to assemble build: %v", err)
		return nil
	}

	logger.Info("executing build")
	// execute the build with the executor
	err = _executor.ExecBuild(ctx)
	if err != nil {
		logger.Errorf("unable to execute build: %v", err)
		return nil
	}

	logger.Info("destroying build")
	// destroy the build with the executor
	err = _executor.DestroyBuild(context.Background())
	if err != nil {
		logger.Errorf("unable to destroy build: %v", err)
	}

	logger.Info("completed build")

	return nil
}
