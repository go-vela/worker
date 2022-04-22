// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package local

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/internal/service"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

// create a service logging pattern.
const servicePattern = "[service: %s]"

// CreateService configures the service for execution.
func (c *client) CreateService(ctx context.Context, ctn *pipeline.Container) error {
	// setup the runtime container
	err := c.Runtime.SetupContainer(ctx, ctn)
	if err != nil {
		return err
	}

	// update the service container environment
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/service#Environment
	err = service.Environment(ctn, c.build, c.repo, nil, c.Version)
	if err != nil {
		return err
	}

	// substitute container configuration
	//
	// https://pkg.go.dev/github.com/go-vela/types/pipeline#Container.Substitute
	err = ctn.Substitute()
	if err != nil {
		return err
	}

	return nil
}

// PlanService prepares the service for execution.
func (c *client) PlanService(ctx context.Context, ctn *pipeline.Container) error {
	// update the engine service object
	_service := new(library.Service)
	_service.SetName(ctn.Name)
	_service.SetNumber(ctn.Number)
	_service.SetImage(ctn.Image)
	_service.SetStatus(constants.StatusRunning)
	_service.SetStarted(time.Now().UTC().Unix())
	_service.SetHost(c.build.GetHost())
	_service.SetRuntime(c.build.GetRuntime())
	_service.SetDistribution(c.build.GetDistribution())

	// update the service container environment
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/service#Environment
	err := service.Environment(ctn, c.build, c.repo, _service, c.Version)
	if err != nil {
		return err
	}

	// add a service to a map
	c.services.Store(ctn.ID, _service)

	return nil
}

// ExecService runs a service.
func (c *client) ExecService(ctx context.Context, ctn *pipeline.Container) error {
	// load the service from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/service#Load
	_service, err := service.Load(ctn, &c.services)
	if err != nil {
		return err
	}

	// defer taking a snapshot of the service
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/service#Snapshot
	defer func() { service.Snapshot(ctn, c.build, nil, nil, nil, _service) }()

	// run the runtime container
	err = c.Runtime.RunContainer(ctx, ctn, c.pipeline)
	if err != nil {
		return err
	}

	// trigger StreamService goroutine with logging context
	c.streamRequests <- message.StreamRequest{
		Key:       "service",
		Stream:    c.StreamService,
		Container: ctn,
	}

	return nil
}

// StreamService tails the output for a service.
func (c *client) StreamService(ctx context.Context, ctn *pipeline.Container) error {
	// tail the runtime container
	rc, err := c.Runtime.TailContainer(ctx, ctn)
	if err != nil {
		return err
	}
	defer rc.Close()

	// create a service pattern for log output
	_pattern := fmt.Sprintf(servicePattern, ctn.Name)

	// create new scanner from the container output
	scanner := bufio.NewScanner(rc)

	// scan entire container output
	for scanner.Scan() {
		// ensure we output to stdout
		fmt.Fprintln(os.Stdout, _pattern, scanner.Text())
	}

	return scanner.Err()
}

// DestroyService cleans up services after execution.
func (c *client) DestroyService(ctx context.Context, ctn *pipeline.Container) error {
	// load the service from the client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/service#Load
	_service, err := service.Load(ctn, &c.services)
	if err != nil {
		// create the service from the container
		//
		// https://pkg.go.dev/github.com/go-vela/types/library#ServiceFromContainerEnvironment
		_service = library.ServiceFromContainerEnvironment(ctn)
	}

	// defer an upload of the service
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/service#Upload
	defer func() { service.Upload(ctn, c.build, nil, nil, nil, _service) }()

	// inspect the runtime container
	err = c.Runtime.InspectContainer(ctx, ctn)
	if err != nil {
		return err
	}

	// remove the runtime container
	err = c.Runtime.RemoveContainer(ctx, ctn)
	if err != nil {
		return err
	}

	return nil
}
