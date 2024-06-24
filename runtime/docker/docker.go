// SPDX-License-Identifier: Apache-2.0

package docker

import (
	docker "github.com/docker/docker/client"
	"github.com/sirupsen/logrus"

	mock "github.com/go-vela/worker/mock/docker"
)

// Version represents the supported Docker API version for the mock.
//
// The Docker API version is pinned to ensure compatibility between the
// Docker API and client. The goal is to maintain n-1 compatibility.
//
// The maximum supported Docker API version for the client is here:
//
// https://docs.docker.com/engine/api/#api-version-matrix
//
// For example (use the compatibility matrix above for reference):
//
// * the Docker version of v20.10 has a maximum API version of v1.41
// * to maintain n-1, the API version is pinned to v1.40
// .
const Version = "v1.40"

type config struct {
	// specifies a list of privileged images to use for the Docker client
	Images []string
	// specifies a list of host volumes to use for the Docker client
	Volumes []string
	// specifies a list of kernel capabilities to drop for each Docker container
	DropCapabilities []string
}

type client struct {
	config *config
	// https://pkg.go.dev/github.com/docker/docker/client#CommonAPIClient
	Docker docker.CommonAPIClient
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	Logger *logrus.Entry
}

// New returns an Engine implementation that
// integrates with a Docker runtime.
//
//nolint:revive // ignore returning unexported client
func New(opts ...ClientOpt) (*client, error) {
	// create new Docker client
	c := new(client)

	// create new fields
	c.config = new(config)

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#StandardLogger
	logger := logrus.StandardLogger()

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#NewEntry
	c.Logger = logrus.NewEntry(logger)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// create new Docker client from environment
	//
	// https://pkg.go.dev/github.com/docker/docker/client#NewClientWithOpts
	_docker, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return nil, err
	}

	// pin version to ensure we know what Docker API version we're using
	//
	// typically this would be inherited from the host environment
	// but this ensures the version of client being used
	//
	// https://pkg.go.dev/github.com/docker/docker/client#WithVersion
	_ = docker.WithVersion(Version)(_docker)

	// set the Docker client in the runtime client
	c.Docker = _docker

	return c, nil
}

// NewMock returns an Engine implementation that
// integrates with a mock Docker runtime.
//
// This function is intended for running tests only.
//
//nolint:revive // ignore returning unexported client
func NewMock(opts ...ClientOpt) (*client, error) {
	// create new Docker runtime client
	c, err := New(opts...)
	if err != nil {
		return nil, err
	}

	// create Docker client from the mock client
	//
	// https://pkg.go.dev/github.com/go-vela/worker/mock/docker#New
	_docker, err := mock.New()
	if err != nil {
		return nil, err
	}

	// set the Docker client in the runtime client
	c.Docker = _docker

	return c, nil
}
