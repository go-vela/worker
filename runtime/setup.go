// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package runtime

import (
	"fmt"

	"github.com/go-vela/worker/runtime/docker"
	"github.com/go-vela/worker/runtime/kubernetes"

	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

// Setup represents the configuration necessary for
// creating a Vela engine capable of integrating
// with a configured runtime environment.
type Setup struct {
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	Logger *logrus.Entry

	// Runtime Configuration

	// Mock should only be true for tests.
	Mock bool

	// specifies the driver to use for the runtime client
	Driver string
	// specifies the path to a configuration file to use for the runtime client
	ConfigFile string
	// specifies a list of host volumes to use for the runtime client
	HostVolumes []string
	// specifies the namespace to use for the runtime client (only used by kubernetes)
	Namespace string
	// specifies the name of the PipelinePodsTemplate to retrieve from the given namespace (only used by kubernetes)
	PodsTemplateName string
	// specifies the fallback path of a PipelinePodsTemplate in a local YAML file (only used by kubernetes; only used if PodsTemplateName not defined)
	PodsTemplateFile string
	// specifies a list of privileged images to use for the runtime client
	PrivilegedImages []string
	// specifies settings for restrictions on trusted repos
	EnableTrusted bool
}

// Docker creates and returns a Vela engine capable of
// integrating with a Docker runtime environment.
func (s *Setup) Docker() (Engine, error) {
	logrus.Trace("creating docker runtime client from setup")

	logrus.Trace("creating docker runtime using enable trusted %v", s.EnableTrusted)

	opts := []docker.ClientOpt{
		docker.WithHostVolumes(s.HostVolumes),
		docker.WithPrivilegedImages(s.PrivilegedImages),
		docker.WithEnableTrusted(s.EnableTrusted),
		docker.WithLogger(s.Logger),
	}

	if s.Mock {
		// create new mock Docker runtime engine
		//
		// https://pkg.go.dev/github.com/go-vela/worker/runtime/docker?tab=doc#NewMock
		return docker.NewMock(opts...)
	}

	// create new Docker runtime engine
	//
	// https://pkg.go.dev/github.com/go-vela/worker/runtime/docker?tab=doc#New
	return docker.New(opts...)
}

// Kubernetes creates and returns a Vela engine capable of
// integrating with a Kubernetes runtime environment.
func (s *Setup) Kubernetes() (Engine, error) {
	logrus.Trace("creating kubernetes runtime client from setup")

	opts := []kubernetes.ClientOpt{
		kubernetes.WithConfigFile(s.ConfigFile),
		kubernetes.WithHostVolumes(s.HostVolumes),
		kubernetes.WithNamespace(s.Namespace),
		kubernetes.WithPodsTemplate(s.PodsTemplateName, s.PodsTemplateFile),
		kubernetes.WithPrivilegedImages(s.PrivilegedImages),
		kubernetes.WithEnableTrusted(s.EnableTrusted),
		kubernetes.WithLogger(s.Logger),
	}

	if s.Mock {
		// create new mock Kubernetes runtime engine
		//
		// https://pkg.go.dev/github.com/go-vela/worker/runtime/kubernetes?tab=doc#NewMock
		return kubernetes.NewMock(&v1.Pod{}, opts...)
	}

	// create new Kubernetes runtime engine
	//
	// https://pkg.go.dev/github.com/go-vela/worker/runtime/kubernetes?tab=doc#New
	return kubernetes.New(opts...)
}

// Validate verifies the necessary fields for the
// provided configuration are populated correctly.
func (s *Setup) Validate() error {
	logrus.Trace("validating runtime setup for client")

	// check if a runtime driver was provided
	if len(s.Driver) == 0 {
		return fmt.Errorf("no runtime driver provided")
	}

	// process the secret driver being provided
	switch s.Driver {
	case constants.DriverDocker:
		break
	case constants.DriverKubernetes:
		// check if a runtime namespace was provided
		if len(s.Namespace) == 0 {
			return fmt.Errorf("no runtime namespace provided")
		}
	}

	// setup is valid
	return nil
}
