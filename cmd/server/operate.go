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
				ctx := context.Background()
				item, err := queue.Pop()
				if err != nil {
					return err
				}

				e.WithBuild(item.Build)
				e.WithPipeline(item.Pipeline)
				e.WithRepo(item.Repo)
				e.WithUser(item.User)

				logrus.Infof("creating %s build", item.Repo.GetFullName())
				err = e.CreateBuild(ctx)
				if err != nil {
					logrus.Errorf("unable to create build: %w", err)
					return err
				}

				logrus.Infof("executing %s build", item.Repo.GetFullName())
				err = e.ExecBuild(ctx)
				if err != nil {
					logrus.Errorf("unable to execute build: %w", err)
					return err
				}

				logrus.Infof("destroying %s build", item.Repo.GetFullName())
				err = e.DestroyBuild(ctx)
				if err != nil {
					logrus.Errorf("unable to destroy build: %w", err)
					return err
				}

				logrus.Infof("completed %s build", item.Repo.GetFullName())
			}
		})
	}

	err = threads.Wait()
	if err != nil {
		return err
	}

	return nil
}
