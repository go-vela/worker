// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"context"
	"time"

	"github.com/go-vela/server/queue"
	"github.com/go-vela/types"

	"github.com/sirupsen/logrus"

	"golang.org/x/sync/errgroup"
)

// Operate is a helper function to initiate all
// subprocesses for the operator to poll the
// queue and execute Vela pipelines.
func (w *Worker) Operate(ctx context.Context) error {
	var err error

	// setup the client
	w.VelaClient, err = setupClient(w.Config.Server)
	if err != nil {
		return err
	}

	// create the errgroup for managing operator subprocesses
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group
	executors, gctx := errgroup.WithContext(ctx)

	// Define the database representation of the worker
	// and register itself in the database
	registryWorker := ToLibrary(w)
	registryWorker.SetActive(true)
	registryWorker.SetLastCheckedIn(time.Now().UTC().Unix())

	// spawn goroutine for phoning home
	executors.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				logrus.Info("Completed looping on worker registration")
				return nil
			default:
				// set checking time to now and call the server
				registryWorker.SetLastCheckedIn(time.Now().UTC().Unix())

				// register or update the worker
				//nolint:contextcheck // ignore passing context
				err = w.checkIn(registryWorker)
				if err != nil {
					logrus.Error(err)
				}

				// if unable to update the worker, log the error but allow the worker to continue running
				if err != nil {
					logrus.Errorf("unable to update worker %s on the server: %v", registryWorker.GetHostname(), err)
				}

				// sleep for the configured time
				time.Sleep(w.Config.CheckIn)
			}
		}
	})

	// setup the queue
	//
	// https://pkg.go.dev/github.com/go-vela/server/queue?tab=doc#New
	//nolint:contextcheck // ignore passing context
	w.Queue, err = queue.New(w.Config.Queue)
	if err != nil {
		return err
	}

	// initialize build activity
	w.Activity = NewActivity()

	ch := make(chan *types.BuildPackage)
	w.PackageChannel = ch

	// iterate till the configured build limit
	for i := 0; i < w.Config.Build.Limit; i++ {
		// evaluate and capture i at each iteration
		//
		// https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		id := i

		// log a message indicating the start of an operator thread
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Info
		logrus.Infof("Thread ID %d listening for builds on channel...", id)

		// spawn errgroup routine for operator subprocess
		//
		// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group.Go

		// TODO: disable/enable this based on environment config
		executors.Go(func() error {
			// create an infinite loop to poll for builds
			for {
				switch w.Config.Executor.BuildMode {
				case "push":
					select {
					case <-gctx.Done():
						logrus.WithFields(logrus.Fields{
							"id": id,
						}).Info("Completed listening on worker executor package channel")
						return nil
					case pkg := <-w.PackageChannel:
						logrus.WithFields(logrus.Fields{
							"id": id,
						}).Info("Received execution package over channel.")
						// exec operator subprocess execute build from the queue
						// (do not pass the context to avoid errors in one
						// executor+build inadvertently canceling other builds)
						//nolint:contextcheck // ignore passing context
						err = w.Exec(id, pkg)
						if err != nil {
							// log the error received from the executor
							//
							// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Errorf
							logrus.Errorf("failing worker executor: %v", err)

							return err
						}
					}
					break
				case "pull":
					select {
					case <-gctx.Done():
						logrus.WithFields(logrus.Fields{
							"id": id,
						}).Info("Completed listening on worker executor package channel")
						return nil
					default:
						// capture an item from the queue
						item, err := w.Queue.Pop(context.Background())
						if err != nil {
							return err
						}

						if item == nil {
							return nil
						}
						// exec operator subprocess execute build from the queue
						// (do not pass the context to avoid errors in one
						// executor+build inadvertently canceling other builds)
						//nolint:contextcheck // ignore passing context
						err = w.Exec(id, &types.BuildPackage{
							Build: item.Build,
							// Secrets:  []*library.Secret{},
							Pipeline: item.Pipeline,
							Repo:     item.Repo,
							User:     item.User,
						})
						if err != nil {
							// log the error received from the executor
							//
							// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Errorf
							logrus.Errorf("failing worker executor: %v", err)

							return err
						}
					}
					break
				}
			}
		})
	}

	// wait for errors from operator subprocesses
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group.Wait
	return executors.Wait()
}
