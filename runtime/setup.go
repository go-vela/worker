// SPDX-License-Identifier: Apache-2.0

package runtime

import (
	"fmt"
	"strings"

	"github.com/docker/docker/oci/caps"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/worker/runtime/docker"
	"github.com/go-vela/worker/runtime/kubernetes"
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
	// specifies a list of kernel capabilities to drop from container (only used by Docker)
	DropCapabilities []string
	//
	ContainerPlatform string
}

// Docker creates and returns a Vela engine capable of
// integrating with a Docker runtime environment.
func (s *Setup) Docker() (Engine, error) {
	logrus.Trace("creating docker runtime client from setup")

	opts := []docker.ClientOpt{
		docker.WithHostVolumes(s.HostVolumes),
		docker.WithPrivilegedImages(s.PrivilegedImages),
		docker.WithLogger(s.Logger),
		docker.WithDropCapabilities(s.DropCapabilities),
		docker.WithContainerPlatform(s.ContainerPlatform),
	}

	if s.Mock {
		// create new mock Docker runtime engine
		//
		// https://pkg.go.dev/github.com/go-vela/worker/runtime/docker#NewMock
		return docker.NewMock(opts...)
	}

	// create new Docker runtime engine
	//
	// https://pkg.go.dev/github.com/go-vela/worker/runtime/docker#New
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
		kubernetes.WithLogger(s.Logger),
	}

	if s.Mock {
		// create new mock Kubernetes runtime engine
		//
		// https://pkg.go.dev/github.com/go-vela/worker/runtime/kubernetes#NewMock
		return kubernetes.NewMock(&v1.Pod{}, opts...)
	}

	// create new Kubernetes runtime engine
	//
	// https://pkg.go.dev/github.com/go-vela/worker/runtime/kubernetes#New
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
		// check to make sure drop capabilities is configured correctly
		if len(s.DropCapabilities) != 0 {
			for _, configCap := range s.DropCapabilities {
				valid := false

				for _, validCap := range caps.DefaultCapabilities() {
					if strings.EqualFold(configCap, validCap) {
						valid = true
						break
					}
				}

				if !valid {
					return fmt.Errorf("invalid capability %s provided in RUNTIME_DROP_CAPABILITIES", configCap)
				}
			}
		}
	case constants.DriverKubernetes:
		// check if a runtime namespace was provided
		if len(s.Namespace) == 0 {
			return fmt.Errorf("no runtime namespace provided")
		}
	}

	// setup is valid
	return nil
}
