// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"github.com/go-vela/worker/runtime/docker/testdata/mock"
	docker "github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

const dockerVersion = "1.38"

type client struct {
	Runtime *docker.Client
}

// New returns an Engine implementation that
// integrates with a Docker runtime.
func New() (*client, error) {
	// create Docker client from environment
	r, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}
	// pin version to prevent "client version <version> is too new." errors
	// typically this would be inherited from the host env but this will ensure
	// we know what version of the Docker API we're using
	docker.WithVersion(dockerVersion)(r)

	// create the client object
	c := &client{
		Runtime: r,
	}

	return c, nil
}

// NewMock returns an Engine implementation that
// integrates with a mock Docker runtime.
//
// This function is intended for running tests only.
func NewMock() (*client, error) {
	// create mock client
	mock := mock.Client(mock.Router)

	// create Docker client from the mock client
	r, err := docker.NewClient("tcp://127.0.0.1:2333", dockerVersion, mock, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	// create the client object
	c := &client{
		Runtime: r,
	}

	return c, nil
}
