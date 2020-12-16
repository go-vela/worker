// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"time"

	"github.com/go-vela/pkg-queue/queue"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"

	"golang.org/x/sync/errgroup"
)

// operate is a helper function to initiate all
// subprocesses for the operator to poll the
// queue and execute Vela pipelines.
func (w *Worker) operate() error {
	var err error

	// setup the client
	w.VelaClient, err = setupClient(w.Config.Server)
	if err != nil {
		return err
	}
	// Define the database representation of the worker
	// and register itself in the database
	registryWorker := new(library.Worker)
	registryWorker.SetHostname(w.Config.Hostname)
	registryWorker.SetAddress(w.Config.Server.Address)
	registryWorker.SetRoutes(w.Config.Queue.Routes)
	registryWorker.SetActive(true)
	registryWorker.SetLastCheckedIn(time.Now().UTC().Unix())
	registryWorker.SetBuildLimit(int64(w.Config.Build.Limit))
	err = w.register(registryWorker)
	if err != nil {
		logrus.Error(err)
	}

	// spawn goroutine for phoning home
	go func() {
		for {
			// sleep for the configured time
			time.Sleep(w.Config.CheckIn)

			// set checking time to now and call the server
			registryWorker.SetLastCheckedIn(time.Now().UTC().Unix())
			_, _, err := w.VelaClient.Worker.Update(registryWorker.GetHostname(), registryWorker)

			// if unable to update the worker, log the error but allow the worker to continue running
			if err != nil {
				logrus.Errorf("unable to update worker %s on the server: %v", registryWorker.GetHostname(), err)
			}
		}
	}()

	// setup the queue
	//
	// https://pkg.go.dev/github.com/go-vela/pkg-queue/queue?tab=doc#New
	w.Queue, err = queue.New(w.Config.Queue)
	if err != nil {
		return err
	}

	// create the errgroup for managing operator subprocesses
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group
	executors := new(errgroup.Group)

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
				// exec operator subprocess to poll and execute builds
				err = w.exec(id)
				if err != nil {
					// log the error received from the executor
					//
					// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Errorf
					logrus.Errorf("failing worker executor: %v", err)

					return err
				}
			}
		})
	}

	// wait for errors from operator subprocesses
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group.Wait
	return executors.Wait()
}
