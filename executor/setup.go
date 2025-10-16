// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/storage"
	"github.com/go-vela/worker/executor/linux"
	"github.com/go-vela/worker/executor/local"
	"github.com/go-vela/worker/runtime"
)

// Setup represents the configuration necessary for
// creating a Vela engine capable of integrating
// with a configured executor.
type Setup struct {
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	Logger *logrus.Entry

	// Executor Configuration

	// Mock should only be true for tests.
	Mock bool

	// specifies the executor driver to use
	Driver string
	// specifies the maximum log size
	MaxLogSize uint
	// specifies how long to wait after the build finishes
	// for log streaming to complete
	LogStreamingTimeout time.Duration
	// specifies a list of privileged images to use
	PrivilegedImages []string
	// configuration for enforcing that only trusted repos may run privileged images
	EnforceTrustedRepos bool
	// specifies the executor hostname
	Hostname string
	// specifies the executor version
	Version string
	// API client for sending requests to Vela
	Client *vela.Client
	// engine used for creating runtime resources
	Runtime runtime.Engine

	OutputCtn *pipeline.Container

	// Vela Resource Configuration

	// resource for storing build information in Vela
	Build *api.Build
	// resource for storing pipeline information in Vela
	Pipeline *pipeline.Build
	// id token request token for the build
	RequestToken string
	// storage client for interacting with storage resources
	Storage storage.Storage
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

	// create options for Linux executor
	opts := []linux.Opt{
		linux.WithBuild(s.Build),
		linux.WithMaxLogSize(s.MaxLogSize),
		linux.WithLogStreamingTimeout(s.LogStreamingTimeout),
		linux.WithPrivilegedImages(s.PrivilegedImages),
		linux.WithEnforceTrustedRepos(s.EnforceTrustedRepos),
		linux.WithHostname(s.Hostname),
		linux.WithPipeline(s.Pipeline),
		linux.WithRuntime(s.Runtime),
		linux.WithVelaClient(s.Client),
		linux.WithVersion(s.Version),
		linux.WithLogger(s.Logger),
		linux.WithOutputCtn(s.OutputCtn),
	}

	// Conditionally add storage option
	if s.Storage != nil {
		fmt.Printf("Adding storage to linux executor\n")
		opts = append(opts, linux.WithStorage(s.Storage))
	}

	// create new Linux executor engine
	//
	// https://pkg.go.dev/github.com/go-vela/worker/executor/linux#New
	return linux.New(opts...)

}

// Local creates and returns a Vela engine capable of
// integrating with a local executor.
func (s *Setup) Local() (Engine, error) {
	logrus.Trace("creating local executor client from setup")

	opts := []local.Opt{
		local.WithBuild(s.Build),
		local.WithHostname(s.Hostname),
		local.WithPipeline(s.Pipeline),
		local.WithRuntime(s.Runtime),
		local.WithVelaClient(s.Client),
		local.WithVersion(s.Version),
		local.WithMockStdout(s.Mock),
		local.WithOutputCtn(s.OutputCtn),
	}

	// Conditionally add storage option
	if s.Storage != nil {
		opts = append(opts, local.WithStorage(s.Storage))
	}

	// create new Local executor engine
	//
	// https://pkg.go.dev/github.com/go-vela/worker/executor/local#New
	return local.New(opts...)
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
	if s.Build.Repo == nil {
		return fmt.Errorf("no Vela repo provided in setup")
	}

	// check if a Vela user was provided
	if s.Build.Repo.Owner == nil {
		return fmt.Errorf("no Vela user provided in setup")
	}

	// If storage is provided, ensure it's enabled
	if s.Storage != nil && !s.Storage.StorageEnable() {
		return fmt.Errorf("storage client provided but not enabled in setup")
	}

	// setup is valid
	return nil
}
