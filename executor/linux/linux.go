// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package linux

import (
	"fmt"
	"os"
	"sync"

	"github.com/go-vela/worker/executor"

	"github.com/go-vela/pkg-runtime/runtime"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"

	"github.com/sirupsen/logrus"
)

type client struct {
	Vela     *vela.Client
	Runtime  runtime.Engine
	Secrets  map[string]*library.Secret
	Hostname string

	// private fields
	logger      *logrus.Entry
	build       *library.Build
	pipeline    *pipeline.Build
	repo        *library.Repo
	services    sync.Map
	serviceLogs sync.Map
	steps       sync.Map
	stepLogs    sync.Map
	user        *library.User
	err         error
}

// New returns an Executor implementation that integrates with a Linux instance.
func New(c *vela.Client, r runtime.Engine) (*client, error) {
	// immediately return if a nil Vela client is provided
	if c == nil {
		return nil, fmt.Errorf("empty Vela client provided to executor")
	}

	// immediately return if a nil runtime Engine is provided
	if r == nil {
		return nil, fmt.Errorf("empty runtime provided to executor")
	}

	// capture the hostname
	h, _ := os.Hostname()

	// create the logger object
	l := logrus.WithFields(logrus.Fields{
		"host": h,
	})

	return &client{
		Vela:        c,
		Runtime:     r,
		Hostname:    h,
		logger:      l,
		services:    sync.Map{},
		serviceLogs: sync.Map{},
		steps:       sync.Map{},
		stepLogs:    sync.Map{},
		err:         nil,
	}, nil
}

// WithBuild sets the library build type in the Engine.
func (c *client) WithBuild(b *library.Build) executor.Engine {
	// set build in engine if one is provided
	if b != nil {
		c.build = b
	}

	return c
}

// WithPipeline sets the pipeline Build type in the Engine.
func (c *client) WithPipeline(p *pipeline.Build) executor.Engine {
	// set pipeline in engine if one is provided
	if p != nil {
		c.pipeline = p
	}

	return c
}

// WithRepo sets the library Repo type in the Engine.
func (c *client) WithRepo(r *library.Repo) executor.Engine {
	// set repo in engine if one is provided
	if r != nil {
		c.repo = r
	}

	return c
}

// WithUser sets the library User type in the Engine.
func (c *client) WithUser(u *library.User) executor.Engine {
	// set user in engine if one is provided
	if u != nil {
		c.user = u
	}

	return c
}

// GetBuild gets the current build in execution.
func (c *client) GetBuild() (*library.Build, error) {
	b := c.build

	// check if the build resource is available
	if b == nil {
		return nil, fmt.Errorf("build resource not found")
	}

	return b, nil
}

// GetPipeline gets the current pipeline in execution.
func (c *client) GetPipeline() (*pipeline.Build, error) {
	p := c.pipeline

	// check if the pipeline resource is available
	if p == nil {
		return nil, fmt.Errorf("pipeline resource not found")
	}

	return p, nil
}

// GetRepo gets the current repo in execution.
func (c *client) GetRepo() (*library.Repo, error) {
	r := c.repo

	// check if the repo resource is available
	if r == nil {
		return nil, fmt.Errorf("repo resource not found")
	}

	return r, nil
}
