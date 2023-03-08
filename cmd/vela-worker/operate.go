// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"context"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/types/library"
	"time"

	"github.com/sirupsen/logrus"

	"golang.org/x/sync/errgroup"
)

// operate is a helper function to initiate all
// subprocesses for the operator to poll the
// queue and execute Vela pipelines.
func (w *Worker) operate(ctx context.Context) error {
	var err error

	// create the errgroup for managing operator subprocesses
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group
	executors, gctx := errgroup.WithContext(ctx)
	w.VelaClient, err = setupClient(w.Config.Server, "")
	if err != nil {
		return err
	}
	// Define the database representation of the worker
	// and register itself in the database
	registryWorker := new(library.Worker)
	registryWorker.SetHostname(w.Config.API.Address.Hostname())
	registryWorker.SetAddress(w.Config.API.Address.String())
	registryWorker.SetRoutes(w.Config.Queue.Routes)
	registryWorker.SetActive(true)
	registryWorker.SetLastCheckedIn(time.Now().UTC().Unix())
	registryWorker.SetBuildLimit(int64(w.Config.Build.Limit))
	if len(w.Config.Server.RegistrationToken) > 0 {
		logrus.Infof("Found seeded token! Attempting to validate!")
		w.Valid, err = w.VelaClient.Authentication.ValidateAndSetToken(w.Config.Server.RegistrationToken)
		if err != nil {
			logrus.Infof("seeded token %s", w.Config.Server.RegistrationToken)
			logrus.Errorf("Expired seeded token %s", err)
		}
	}
	// if seeded token is available, check if valid
	// if valid, set w.Valid to true and do a checkIn
	// if !valid, set w.Valid to false
	// while !valid, wait for new token to be provided

	// spawn goroutine for phoning home
	executors.Go(func() error {
		for {
			if !w.Valid {
				// if no seeded token then wait for token before continuing
				logrus.Info("verifying token is present in channel")
				// wait for token
				token := <-w.AuthToken
				// continue operation like normal
				logrus.Info("token present, continuing operation")
				w.Valid, _ = w.VelaClient.Authentication.ValidateAndSetToken(token)
				logrus.Infof("value of Valid %b", w.Valid)
			}
			if w.Valid {
				select {
				case <-gctx.Done():
					logrus.Info("Completed looping on worker registration")
					return nil
				default:
					logrus.Info("Checking worker!")
					// set checking time to now and call the server
					registryWorker.SetLastCheckedIn(time.Now().UTC().Unix())
					//registryWorker.GetLastCheckedIn()
					logrus.Infof("velaClient auth %s", w.VelaClient.Authentication.HasTokenAuth())
					// register or update the worker
					//nolint:contextcheck // ignore passing context
					// if unable to update the worker, log the error but allow the worker to continue running
					w.CheckedIn, err = w.checkIn(registryWorker)
					if err != nil {
						logrus.Errorf("unable to update worker %s on the server: %v", registryWorker.GetHostname(), err)
						logrus.Info("waiting for registration token")

					}
					// send true/false over to let user know whether registration was a success
					w.Success <- w.CheckedIn

					// Send a bool into channel to avoid double registration
					if w.CheckedIn {
						w.Registered <- w.CheckedIn
					} //else {
					//	// clean Registered channel for registering
					//	<-w.Registered
					//}
					// sleep for the configured time
					time.Sleep(w.Config.CheckIn)
				}
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
				logrus.Info("begins queue exec")
				if !w.CheckedIn {
					time.Sleep(5 * time.Second)
					logrus.Info("worker not checked in, skipping queue read")
					continue
				}
				select {
				case <-gctx.Done():
					logrus.WithFields(logrus.Fields{
						"id": id,
					}).Info("Completed looping on worker executor")
					return nil
				default:
					// exec operator subprocess to poll and execute builds
					// (do not pass the context to avoid errors in one
					// executor+build inadvertently canceling other builds)
					//nolint:contextcheck // ignore passing context
					err = w.exec(id)
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

	// reset w.Valid
	//w.Valid = false
	// wait for errors from operator subprocesses
	//
	// https://pkg.go.dev/golang.org/x/sync/errgroup?tab=doc#Group.Wait
	return executors.Wait()
}
