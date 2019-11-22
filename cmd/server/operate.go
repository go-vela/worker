// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"context"

	"github.com/go-vela/worker/executor"

	"github.com/go-vela/worker/queue"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func operate(queue queue.Service, executors map[int]executor.Engine) (err error) {
	threads := new(errgroup.Group)

	for id, e := range executors {
		logrus.Infof("Thread ID %d listening to queue...", id)
		threads.Go(func() error {
			for {
				// create the context for the executor
				ctx := context.Background()

				// pop an item from the queue
				item, err := queue.Pop()
				if err != nil {
					return err
				}

				// create logger with extra metadata
				logger := logrus.WithFields(logrus.Fields{
					"build": item.Build.GetNumber(),
					"repo":  item.Repo.GetFullName(),
				})

				// add build metadata to the executor
				e.WithBuild(item.Build)
				e.WithPipeline(item.Pipeline)
				e.WithRepo(item.Repo)
				e.WithUser(item.User)

				// create the build on the executor
				logger.Infof("creating build")
				err = e.CreateBuild(ctx)
				if err != nil {
					logger.Errorf("unable to create build: %w", err)
					return err
				}

				// execute the build on the executor
				logger.Infof("executing build")
				err = e.ExecBuild(ctx)
				if err != nil {
					logger.Errorf("unable to execute build: %w", err)
					return err
				}

				// destroy the build on the executor
				logger.Info("destroying build")
				err = e.DestroyBuild(ctx)
				if err != nil {
					logger.Errorf("unable to destroy build: %w", err)
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
