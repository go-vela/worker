// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package podman

import (
	"context"

	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/sirupsen/logrus"
)

type config struct {
	// specifies a list of privileged images to use for the Podman client
	Images []string
	// specifies a list of host volumes to use for the Podman client
	Volumes []string
}

type client struct {
	config *config
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings#NewConnection
	Podman context.Context
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	Logger *logrus.Entry
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/domain/entities#PodSpec
	Pod *entities.PodSpec
}

// New returns an Engine implementation that
// integrates with a Podman runtime.
//
// nolint: golint // ignore returning unexported client
func New(opts ...ClientOpt) (*client, error) {
	// create new Docker client
	c := new(client)

	// create new fields
	c.config = new(config)

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#StandardLogger
	logger := logrus.StandardLogger()

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#NewEntry
	c.Logger = logrus.NewEntry(logger)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// create new Podman client from environment
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings#NewConnection
	_podman, err := bindings.NewConnection(context.Background(), "")
	if err != nil {
		return nil, err
	}

	// set the Podman client in the runtime client
	c.Podman = _podman

	return c, nil
}

// NewMock returns an Engine implementation that
// integrates with a mock Docker runtime.
//
// This function is intended for running tests only.
//
// nolint: golint // ignore returning unexported client
// func NewMock(opts ...ClientOpt) (*client, error) {
// 	// create new Docker runtime client
// 	c, err := New(opts...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// create Docker client from the mock client
// 	//
// 	// https://pkg.go.dev/github.com/go-vela/worker/mock/docker#New
// 	_docker, err := mock.New()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// set the Docker client in the runtime client
// 	c.Docker = _docker

// 	return c, nil
// }
