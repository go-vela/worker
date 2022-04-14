// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"context"
	"time"

	"github.com/go-vela/server/queue"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"

	"golang.org/x/sync/errgroup"
)

// operate is a helper function to initiate all
// subprocesses for the operator to poll the
// queue and execute Vela pipelines.
func (w *Worker) operate(ctx context.Context) error {
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
	registryWorker := new(library.Worker)
	registryWorker.SetHostname(w.Config.API.Address.Hostname())
	registryWorker.SetAddress(w.Config.API.Address.String())
	registryWorker.SetRoutes(w.Config.Queue.Routes)
	registryWorker.SetActive(true)
	registryWorker.SetLastCheckedIn(time.Now().UTC().Unix())
	registryWorker.SetBuildLimit(int64(w.Config.Build.Limit))

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
	w.Queue, err = queue.New(w.Config.Queue)
	if err != nil {
		return err
	}

	// iterate till the configured build limit
	for i := 0; i < w.Config.Build.Limit; i++ {
		// evaluate and capture i at each iteration
		//
		// https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		id := i

		// log a message indicating the start of an operator thread
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Info
		logrus.Infof("Thread ID %d listening to queue...", id)

		// spawn errgroup routine for operator subprocess
		//
		// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group.Go
		executors.Go(func() error {
			// create an infinite loop to poll for builds
			for {
				select {
				case <-gctx.Done():
					logrus.WithFields(logrus.Fields{
						"id": id,
					}).Info("Completed looping on worker executor")
					return nil
				default:
					// exec operator subprocess to poll and execute builds
					err = w.exec(gctx, id)
					if err != nil {
						// log the error received from the executor
						//
						// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Errorf
						logrus.Errorf("failing worker executor: %v", err)

						return err
					}
				}
			}
		})
	}

	// wait for errors from operator subprocesses
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group.Wait
	return executors.Wait()
}
