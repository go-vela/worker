// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-vela/pkg-executor/executor"
	"github.com/go-vela/pkg-runtime/runtime"

	"github.com/sirupsen/logrus"
)

// exec is a helper function to poll the queue
// and execute Vela pipelines for the Worker.
// nolint:funlen
func (w *Worker) exec(index int) error {
	var err error

	// setup the runtime
	//
	// https://pkg.go.dev/github.com/go-vela/pkg-runtime/runtime?tab=doc#New
	w.Runtime, err = runtime.New(w.Config.Runtime)
	if err != nil {
		return err
	}

	// capture an item from the queue
	item, err := w.Queue.Pop()
	if err != nil {
		return err
	}

	// setup the executor
	//
	// https://godoc.org/github.com/go-vela/pkg-executor/executor#New
	_executor, err := executor.New(&executor.Setup{
		Driver:   w.Config.Executor.Driver,
		Client:   w.VelaClient,
		Hostname: w.Config.API.Address.Hostname(),
		Runtime:  w.Runtime,
		Build:    item.Build,
		Pipeline: item.Pipeline.Sanitize(w.Config.Runtime.Driver),
		Repo:     item.Repo,
		User:     item.User,
	})

	// add the executor to the worker
	w.Executors[index] = _executor

	// create logger with extra metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#WithFields
	logger := logrus.WithFields(logrus.Fields{
		"build": item.Build.GetNumber(),
		"host":  w.Config.API.Address.Hostname(),
		"repo":  item.Repo.GetFullName(),
	})

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

	// create channel for catching OS signals
	sigchan := make(chan os.Signal, 1)

	// add a cancelation signal to our current context
	ctx, sig := context.WithCancel(ctx)

	// set the OS signals the Worker will respond to
	signal.Notify(sigchan, syscall.SIGTERM)

	// defer canceling the context
	defer func() {
		signal.Stop(sigchan)
		sig()
	}()

	// spawn a goroutine to listen for the signals
	go func() {
		select {
		case <-sigchan:
			sig()
		case <-ctx.Done():
		}
	}()

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
