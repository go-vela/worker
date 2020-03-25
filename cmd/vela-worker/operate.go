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
	"github.com/go-vela/pkg-queue/queue"
	"github.com/go-vela/pkg-runtime/runtime"

	"github.com/sirupsen/logrus"

	"golang.org/x/sync/errgroup"
)

// operate is a helper function to ...
func (w *Worker) operate() error {
	// setup the client
	_client, err := setupClient(w.Server)
	if err != nil {
		logrus.Fatal(err)
	}

	// setup the queue
	_queue, err := queue.New(w.Queue)
	if err != nil {
		logrus.Fatal(err)
	}

	executors := new(errgroup.Group)

	for i := 0; i < w.Build.Limit; i++ {
		logrus.Infof("Thread ID %d listening to queue...", i)
		executors.Go(func() error {
			for {
				// setup the runtime
				_runtime, err := runtime.New(w.Runtime)
				if err != nil {
					logrus.Fatal(err)
				}

				// capture an item from the queue
				item, err := _queue.Pop()
				if err != nil {
					return err
				}

				_executor, err := executor.New(&executor.Setup{
					Driver:   w.Executor.Driver,
					Client:   _client,
					Runtime:  _runtime,
					Build:    item.Build,
					Pipeline: item.Pipeline.Sanitize(w.Runtime.Driver),
					Repo:     item.Repo,
					User:     item.User,
				})

				// create logger with extra metadata
				logger := logrus.WithFields(logrus.Fields{
					"build": item.Build.GetNumber(),
					"repo":  item.Repo.GetFullName(),
				})

				t := w.Build.Timeout
				// check if the repository has a custom timeout
				if item.Repo.GetTimeout() > 0 {
					// update timeout variable to repository custom timeout
					t = time.Duration(item.Repo.GetTimeout()) * time.Minute
				}

				ctx := context.Background()

				// add to the background context with a timeout
				// built in for ensuring a build doesn't run forever
				ctx, timeout := context.WithTimeout(ctx, t)
				defer timeout()

				// add signals to the parent context so
				// users can cancel builds
				sigchan := make(chan os.Signal, 1)
				ctx, sig := context.WithCancel(ctx)
				signal.Notify(sigchan, syscall.SIGTERM)
				defer func() {
					signal.Stop(sigchan)
					sig()
				}()
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
					return err
				}

				logger.Info("planning build")
				// plan the build with the executor
				err = _executor.PlanBuild(ctx)
				if err != nil {
					logger.Errorf("unable to plan build: %v", err)
					return err
				}

				logger.Info("executing build")
				// execute the build with the executor
				err = _executor.ExecBuild(ctx)
				if err != nil {
					logger.Errorf("unable to execute build: %v", err)
					return err
				}

				logger.Info("destroying build")
				// destroy the build with the executor
				err = _executor.DestroyBuild(context.Background())
				if err != nil {
					logger.Errorf("unable to destroy build: %v", err)
				}

				logger.Info("completed build")
			}
		})
	}

	return executors.Wait()
}
