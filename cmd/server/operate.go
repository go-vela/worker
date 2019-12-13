// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-vela/worker/executor"

	"github.com/go-vela/worker/queue"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func operate(q queue.Service, e map[int]executor.Engine, t time.Duration) (err error) {
	threads := new(errgroup.Group)

	for id, executor := range e {
		logrus.Infof("Thread ID %d listening to queue...", id)
		threads.Go(func() error {
			for {
				// pop an item from the queue
				item, err := q.Pop()
				if err != nil {
					return err
				}

				// create logger with extra metadata
				logger := logrus.WithFields(logrus.Fields{
					"build": item.Build.GetNumber(),
					"repo":  item.Repo.GetFullName(),
				})

				// add build metadata to the executor
				executor.WithBuild(item.Build)
				executor.WithPipeline(item.Pipeline)
				executor.WithRepo(item.Repo)
				executor.WithUser(item.User)

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

				logger.Info("pulling secrets")
				// pull secrets for the build on the executor
				err = executor.PullSecret(ctx)
				if err != nil {
					logger.Errorf("unable to pull secrets: %v", err)
					return err
				}

				defer func() {
					// destroying the build on the executor
					logger.Info("destroying build")
					err = executor.DestroyBuild(context.Background())
					if err != nil {
						logger.Errorf("unable to destroy build: %v", err)
					}
				}()

				// create the build on the executor
				logger.Info("creating build")
				err = executor.CreateBuild(ctx)
				if err != nil {
					logger.Errorf("unable to create build: %v", err)
					return err
				}

				// execute the build on the executor
				logger.Info("executing build")
				err = executor.ExecBuild(ctx)
				if err != nil {
					logger.Errorf("unable to execute build: %v", err)
					return err
				}

				logger.Info("completed build")
			}
		})
	}

	err = threads.Wait()
	if err != nil {
		return err
	}

	return nil
}
