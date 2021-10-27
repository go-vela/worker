// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executor

import (
	"fmt"
	"strings"

	"github.com/go-vela/sdk-go/vela"

	"github.com/go-vela/worker/executor/linux"
	"github.com/go-vela/worker/executor/local"

	"github.com/go-vela/worker/runtime"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"

	"github.com/sirupsen/logrus"
)

// Setup represents the configuration necessary for
// creating a Vela engine capable of integrating
// with a configured executor.
type Setup struct {
	// Executor Configuration

	// specifies the executor driver to use
	Driver string
	// specifies the executor hostname
	Hostname string
	// specifies the executor version
	Version string
	// API client for sending requests to Vela
	Client *vela.Client
	// engine used for creating runtime resources
	Runtime runtime.Engine

	// Vela Resource Configuration

	// resource for storing build information in Vela
	Build *library.Build
	// resource for storing pipeline information in Vela
	Pipeline *pipeline.Build
	// resource for storing repo information in Vela
	Repo *library.Repo
	// resource for storing user information in Vela
	User *library.User
}

// Darwin creates and returns a Vela engine capable of
// integrating with a Darwin executor.
func (s *Setup) Darwin() (Engine, error) {
	logrus.Trace("creating darwin executor client from setup")

	return nil, fmt.Errorf("unsupported executor driver: %s", constants.DriverDarwin)
}

// Linux creates and returns a Vela engine capable of
// integrating with a Linux executor.
func (s *Setup) Linux() (Engine, error) {
	logrus.Trace("creating linux executor client from setup")

	// create new Linux executor engine
	//
	// https://pkg.go.dev/github.com/go-vela/worker/executor/linux?tab=doc#New
	return linux.New(
		linux.WithBuild(s.Build),
		linux.WithHostname(s.Hostname),
		linux.WithPipeline(s.Pipeline),
		linux.WithRepo(s.Repo),
		linux.WithRuntime(s.Runtime),
		linux.WithUser(s.User),
		linux.WithVelaClient(s.Client),
		linux.WithVersion(s.Version),
	)
}

// Local creates and returns a Vela engine capable of
// integrating with a local executor.
func (s *Setup) Local() (Engine, error) {
	logrus.Trace("creating local executor client from setup")

	// create new Local executor engine
	//
	// https://pkg.go.dev/github.com/go-vela/worker/executor/local?tab=doc#New
	return local.New(
		local.WithBuild(s.Build),
		local.WithHostname(s.Hostname),
		local.WithPipeline(s.Pipeline),
		local.WithRepo(s.Repo),
		local.WithRuntime(s.Runtime),
		local.WithUser(s.User),
		local.WithVelaClient(s.Client),
		local.WithVersion(s.Version),
	)
}

// Windows creates and returns a Vela engine capable of
// integrating with a Windows executor.
func (s *Setup) Windows() (Engine, error) {
	logrus.Trace("creating windows executor client from setup")

	return nil, fmt.Errorf("unsupported executor driver: %s", constants.DriverWindows)
}

// Validate verifies the necessary fields for the
// provided configuration are populated correctly.
func (s *Setup) Validate() error {
	logrus.Trace("validating executor setup for client")

	// check if an executor driver was provided
	if len(s.Driver) == 0 {
		return fmt.Errorf("no executor driver provided in setup")
	}

	// check if a Vela pipeline was provided
	if s.Pipeline == nil {
		return fmt.Errorf("no Vela pipeline provided in setup")
	}

	// check if a runtime engine was provided
	if s.Runtime == nil {
		return fmt.Errorf("no runtime engine provided in setup")
	}

	// check if the local driver is provided
	if strings.EqualFold(constants.DriverLocal, s.Driver) {
		// all other fields are not required
		// for the local executor
		return nil
	}

	// check if a Vela client was provided
	if s.Client == nil {
		return fmt.Errorf("no Vela client provided in setup")
	}

	// check if a Vela build was provided
	if s.Build == nil {
		return fmt.Errorf("no Vela build provided in setup")
	}

	// check if a Vela repo was provided
	if s.Repo == nil {
		return fmt.Errorf("no Vela repo provided in setup")
	}

	// check if a Vela user was provided
	if s.User == nil {
		return fmt.Errorf("no Vela user provided in setup")
	}

	// setup is valid
	return nil
}
