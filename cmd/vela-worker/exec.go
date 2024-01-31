// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/runtime"
	"github.com/go-vela/worker/version"
	"github.com/sirupsen/logrus"
)

// exec is a helper function to poll the queue
// and execute Vela pipelines for the Worker.
//
//nolint:nilerr,funlen // ignore returning nil - don't want to crash worker
func (w *Worker) exec(index int, config *library.Worker) error {
	var err error

	// setup the version
	v := version.New()

	// get worker from database
	worker, _, err := w.VelaClient.Worker.Get(w.Config.API.Address.Hostname())
	if err != nil {
		logrus.Errorf("unable to retrieve worker from server: %s", err)

		return err
	}

	// capture an item from the queue
	item, err := w.Queue.Pop(context.Background(), worker.GetRoutes())
	if err != nil {
		logrus.Errorf("queue pop failed: %v", err)

		// returning immediately on queue pop fail will attempt
		// to pop in quick succession, so we honor the configured timeout
		time.Sleep(w.Config.Queue.Timeout)

		// returning nil to avoid unregistering the worker on pop failure;
		// sometimes queue could be unavailable due to blip or maintenance
		return nil
	}

	if item == nil {
		return nil
	}

	// retrieve a build token from the server to setup the execBuildClient
	bt, resp, err := w.VelaClient.Build.GetBuildToken(item.Repo.GetOrg(), item.Repo.GetName(), item.Build.GetNumber())
	if err != nil {
		logrus.Errorf("unable to retrieve build token: %s", err)

		// build is not in pending state — user canceled build while it was in queue. Pop, discard, move on.
		if resp != nil && resp.StatusCode == http.StatusConflict {
			return nil
		}

		// something else is amiss (auth, server down, etc.) — shut down worker, will have to re-register if registration enabled.
		return err
	}

	// set up build client with build token as auth
	execBuildClient, err := setupClient(w.Config.Server, bt.GetToken())
	if err != nil {
		return err
	}

	// request build executable containing pipeline.Build data using exec client
	execBuildExecutable, _, err := execBuildClient.Build.GetBuildExecutable(item.Repo.GetOrg(), item.Repo.GetName(), item.Build.GetNumber())
	if err != nil {
		return err
	}

	// get the build pipeline from the build executable
	p := new(pipeline.Build)

	err = json.Unmarshal(execBuildExecutable.GetData(), p)
	if err != nil {
		return err
	}

	// set the outputs container ID
	w.Config.Executor.OutputCtn.ID = fmt.Sprintf("outputs_%s", p.ID)

	// create logger with extra metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#WithFields
	logger := logrus.WithFields(logrus.Fields{
		"build":    item.Build.GetNumber(),
		"executor": w.Config.Executor.Driver,
		"host":     w.Config.API.Address.Hostname(),
		"repo":     item.Repo.GetFullName(),
		"runtime":  w.Config.Runtime.Driver,
		"user":     item.User.GetName(),
		"version":  v.Semantic(),
	})

	// lock and append the build to the RunningBuildIDs list
	w.RunningBuildIDsMutex.Lock()

	w.RunningBuildIDs = append(w.RunningBuildIDs, strconv.FormatInt(item.Build.GetID(), 10))

	config.SetRunningBuildIDs(w.RunningBuildIDs)

	w.RunningBuildIDsMutex.Unlock()

	// set worker status
	updateStatus := w.getWorkerStatusFromConfig(config)
	config.SetStatus(updateStatus)
	config.SetLastStatusUpdateAt(time.Now().Unix())
	config.SetLastBuildStartedAt(time.Now().Unix())

	// update worker in the database
	_, _, err = w.VelaClient.Worker.Update(config.GetHostname(), config)
	if err != nil {
		logger.Errorf("unable to update worker: %v", err)
	}

	// handle stale item queued before a Vela upgrade or downgrade.
	if item.ItemVersion != types.ItemVersion {
		// If the ItemVersion is older or newer than what we expect, then it might
		// not be safe to process the build. Fail the build and loop to the next item.
		// TODO: Ask the server to re-compile and requeue the build instead of failing it.
		logrus.Errorf("Failing stale queued build due to wrong item version: want %d, got %d", types.ItemVersion, item.ItemVersion)

		build := item.Build
		build.SetError("Unable to process stale build (queued before Vela upgrade/downgrade).")
		build.SetStatus(constants.StatusError)
		build.SetFinished(time.Now().UTC().Unix())

		_, _, err := execBuildClient.Build.Update(item.Repo.GetOrg(), item.Repo.GetName(), build)
		if err != nil {
			logrus.Errorf("Unable to set build status to %s: %s", constants.StatusFailure, err)
			return err
		}

		return nil
	}

	// setup the runtime
	//
	// https://pkg.go.dev/github.com/go-vela/worker/runtime#New
	w.Runtime, err = runtime.New(&runtime.Setup{
		Logger:           logger,
		Mock:             w.Config.Mock,
		Driver:           w.Config.Runtime.Driver,
		ConfigFile:       w.Config.Runtime.ConfigFile,
		HostVolumes:      w.Config.Runtime.HostVolumes,
		Namespace:        w.Config.Runtime.Namespace,
		PodsTemplateName: w.Config.Runtime.PodsTemplateName,
		PodsTemplateFile: w.Config.Runtime.PodsTemplateFile,
		PrivilegedImages: w.Config.Runtime.PrivilegedImages,
		DropCapabilities: w.Config.Runtime.DropCapabilities,
	})
	if err != nil {
		return err
	}

	// setup the executor
	//
	// https://godoc.org/github.com/go-vela/worker/executor#New
	_executor, err := executor.New(&executor.Setup{
		Logger:              logger,
		Mock:                w.Config.Mock,
		Driver:              w.Config.Executor.Driver,
		MaxLogSize:          w.Config.Executor.MaxLogSize,
		LogStreamingTimeout: w.Config.Executor.LogStreamingTimeout,
		EnforceTrustedRepos: w.Config.Executor.EnforceTrustedRepos,
		PrivilegedImages:    w.Config.Runtime.PrivilegedImages,
		Client:              execBuildClient,
		Hostname:            w.Config.API.Address.Hostname(),
		Runtime:             w.Runtime,
		Build:               item.Build,
		Pipeline:            p.Sanitize(w.Config.Runtime.Driver),
		Repo:                item.Repo,
		User:                item.User,
		Version:             v.Semantic(),
		OutputCtn:           w.Config.Executor.OutputCtn,
	})

	// add the executor to the worker
	w.Executors[index] = _executor

	// This WaitGroup delays calling DestroyBuild until the StreamBuild goroutine finishes.
	var wg sync.WaitGroup

	// this gets deferred first so that DestroyBuild runs AFTER the
	// new contexts (buildCtx and timeoutCtx) have been canceled
	defer func() {
		// if exec() exits before starting StreamBuild, this returns immediately.
		wg.Wait()

		logger.Info("destroying build")

		// destroy the build with the executor (pass a background
		// context to guarantee all build resources are destroyed).
		err = _executor.DestroyBuild(context.Background())
		if err != nil {
			logger.Errorf("unable to destroy build: %v", err)
		}

		logger.Info("completed build")

		// lock and remove the build from the RunningBuildIDs list
		w.RunningBuildIDsMutex.Lock()

		for i, v := range w.RunningBuildIDs {
			if v == strconv.FormatInt(item.Build.GetID(), 10) {
				w.RunningBuildIDs = append(w.RunningBuildIDs[:i], w.RunningBuildIDs[i+1:]...)
			}
		}

		config.SetRunningBuildIDs(w.RunningBuildIDs)

		w.RunningBuildIDsMutex.Unlock()

		// set worker status
		updateStatus := w.getWorkerStatusFromConfig(config)
		config.SetStatus(updateStatus)
		config.SetLastStatusUpdateAt(time.Now().Unix())
		config.SetLastBuildFinishedAt(time.Now().Unix())

		// update worker in the database
		_, _, err := w.VelaClient.Worker.Update(config.GetHostname(), config)
		if err != nil {
			logger.Errorf("unable to update worker: %v", err)
		}
	}()

	// capture the configured build timeout
	t := w.Config.Build.Timeout
	// check if the repository has a custom timeout
	if item.Repo.GetTimeout() > 0 {
		// update timeout variable to repository custom timeout
		t = time.Duration(item.Repo.GetTimeout()) * time.Minute
	}

	// create a build context (from a background context
	// so that other builds can't inadvertently cancel this build)
	buildCtx, done := context.WithCancel(context.Background())
	defer done()

	// add to the background context with a timeout
	// built in for ensuring a build doesn't run forever
	timeoutCtx, timeout := context.WithTimeout(buildCtx, t)
	defer timeout()

	logger.Info("creating build")
	// create the build with the executor
	err = _executor.CreateBuild(timeoutCtx)
	if err != nil {
		logger.Errorf("unable to create build: %v", err)
		return nil
	}

	logger.Info("planning build")
	// plan the build with the executor
	err = _executor.PlanBuild(timeoutCtx)
	if err != nil {
		logger.Errorf("unable to plan build: %v", err)
		return nil
	}

	logger.Info("assembling build")
	// assemble the build with the executor
	err = _executor.AssembleBuild(timeoutCtx)
	if err != nil {
		logger.Errorf("unable to assemble build: %v", err)
		return nil
	}

	// add StreamBuild goroutine to WaitGroup
	wg.Add(1)

	// log/event streaming uses buildCtx so that it is not subject to the timeout.
	go func() {
		defer wg.Done()
		logger.Info("streaming build logs")
		// execute the build with the executor
		err = _executor.StreamBuild(buildCtx)
		if err != nil {
			logger.Errorf("unable to stream build logs: %v", err)
		}
	}()

	logger.Info("executing build")
	// execute the build with the executor
	err = _executor.ExecBuild(timeoutCtx)
	if err != nil {
		logger.Errorf("unable to execute build: %v", err)
		return nil
	}

	return nil
}

// getWorkerStatusFromConfig is a helper function
// to determine the appropriate worker status.
func (w *Worker) getWorkerStatusFromConfig(config *library.Worker) string {
	switch rb := len(config.GetRunningBuildIDs()); {
	case rb == 0:
		return constants.WorkerStatusIdle
	case rb < w.Config.Build.Limit:
		return constants.WorkerStatusAvailable
	case rb == w.Config.Build.Limit:
		return constants.WorkerStatusBusy
	default:
		return constants.WorkerStatusError
	}
}
