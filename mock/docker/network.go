// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/stringid"
)

// NetworkService implements all the network
// related functions for the Docker mock.
type NetworkService struct{}

// NetworkConnect is a helper function to simulate
// a mocked call to connect to a Docker network.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.NetworkConnect
func (n *NetworkService) NetworkConnect(ctx context.Context, network, container string, config *network.EndpointSettings) error {
	return nil
}

// NetworkCreate is a helper function to simulate
// a mocked call to create a Docker network.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.NetworkCreate
func (n *NetworkService) NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error) {
	// verify a network was provided
	if len(name) == 0 {
		return types.NetworkCreateResponse{}, errors.New("no network provided")
	}

	// check if the network is notfound and
	// check if the notfound should be ignored
	if strings.Contains(name, "notfound") &&
		!strings.Contains(name, "ignorenotfound") {
		return types.NetworkCreateResponse{},
			// nolint:golint,stylecheck // messsage is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such network: %s", name))
	}

	// check if the network is not-found and
	// check if the not-found should be ignored
	if strings.Contains(name, "not-found") &&
		!strings.Contains(name, "ignore-not-found") {
		return types.NetworkCreateResponse{},
			// nolint:golint,stylecheck // messsage is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such network: %s", name))
	}

	// create response object to return
	response := types.NetworkCreateResponse{
		ID: stringid.GenerateRandomID(),
	}

	return response, nil
}

// NetworkDisconnect is a helper function to simulate
// a mocked call to disconnect from a Docker network.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.NetworkDisconnect
func (n *NetworkService) NetworkDisconnect(ctx context.Context, network, container string, force bool) error {
	return nil
}

// NetworkInspect is a helper function to simulate
// a mocked call to inspect a Docker network.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.NetworkInspect
func (n *NetworkService) NetworkInspect(ctx context.Context, network string, options types.NetworkInspectOptions) (types.NetworkResource, error) {
	// verify a network was provided
	if len(network) == 0 {
		return types.NetworkResource{}, errors.New("no network provided")
	}

	// check if the network is notfound
	if strings.Contains(network, "notfound") {
		return types.NetworkResource{},
			// nolint:golint,stylecheck // messsage is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such network: %s", network))
	}

	// check if the network is not-found
	if strings.Contains(network, "not-found") {
		return types.NetworkResource{},
			// nolint:golint,stylecheck // messsage is capitalized to match Docker messages
			errdefs.NotFound(fmt.Errorf("Error: No such network: %s", network))
	}

	// create response object to return
	response := types.NetworkResource{
		Attachable: false,
		ConfigOnly: false,
		Created:    time.Now(),
		Driver:     "host",
		ID:         stringid.GenerateRandomID(),
		Ingress:    false,
		Internal:   false,
		Name:       network,
		Scope:      "local",
	}

	return response, nil
}

// NetworkInspectWithRaw is a helper function to simulate
// a mocked call to inspect a Docker network and return
// the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.NetworkInspectWithRaw
func (n *NetworkService) NetworkInspectWithRaw(ctx context.Context, network string, options types.NetworkInspectOptions) (types.NetworkResource, []byte, error) {
	// verify a network was provided
	if len(network) == 0 {
		return types.NetworkResource{}, nil, errors.New("no network provided")
	}

	// create response object to return
	response := types.NetworkResource{
		Attachable: false,
		ConfigOnly: false,
		Created:    time.Now(),
		Driver:     "host",
		ID:         stringid.GenerateRandomID(),
		Ingress:    false,
		Internal:   false,
		Name:       network,
		Scope:      "local",
	}

	// marshal response into raw bytes
	b, err := json.Marshal(response)
	if err != nil {
		return types.NetworkResource{}, nil, err
	}

	return response, b, nil
}

// NetworkList is a helper function to simulate
// a mocked call to list Docker networks.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.NetworkList
func (n *NetworkService) NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error) {
	return nil, nil
}

// NetworkRemove is a helper function to simulate
// a mocked call to remove Docker a network.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.NetworkRemove
func (n *NetworkService) NetworkRemove(ctx context.Context, network string) error {
	// verify a network was provided
	if len(network) == 0 {
		return errors.New("no network provided")
	}

	return nil
}

// NetworksPrune is a helper function to simulate
// a mocked call to prune Docker networks.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.NetworksPrune
func (n *NetworkService) NetworksPrune(ctx context.Context, pruneFilter filters.Args) (types.NetworksPruneReport, error) {
	return types.NetworksPruneReport{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// NetworkService satisfies the NetworkAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#NetworkAPIClient
var _ client.NetworkAPIClient = (*NetworkService)(nil)
