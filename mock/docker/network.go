// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"errors"
	"strings"
	"time"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"github.com/moby/moby/client/pkg/stringid"
)

// NetworkService implements all the network
// related functions for the Docker mock.
type NetworkService struct{}

// NetworkCreate is a helper function to simulate
// a mocked call to create a Docker network.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NetworkCreate
func (n *NetworkService) NetworkCreate(_ context.Context, name string, _ client.NetworkCreateOptions) (client.NetworkCreateResult, error) {
	// verify a network was provided
	if len(name) == 0 {
		return client.NetworkCreateResult{}, errors.New("no network provided")
	}

	// check if the network is notfound and
	// check if the notfound should be ignored
	if strings.Contains(name, "notfound") &&
		!strings.Contains(name, "ignorenotfound") {
		return client.NetworkCreateResult{}, cerrdefs.ErrNotFound
	}

	// check if the network is not-found and
	// check if the not-found should be ignored
	if strings.Contains(name, "not-found") &&
		!strings.Contains(name, "ignore-not-found") {
		return client.NetworkCreateResult{}, cerrdefs.ErrNotFound
	}

	// create response object to return
	response := client.NetworkCreateResult{
		ID: stringid.GenerateRandomID(),
	}

	return response, nil
}

// NetworkInspect is a helper function to simulate
// a mocked call to inspect a Docker network.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NetworkInspect
func (n *NetworkService) NetworkInspect(_ context.Context, networkID string, _ client.NetworkInspectOptions) (client.NetworkInspectResult, error) {
	// verify a network was provided
	if len(networkID) == 0 {
		return client.NetworkInspectResult{}, errors.New("no network provided")
	}

	// check if the network is notfound
	if strings.Contains(networkID, "notfound") {
		return client.NetworkInspectResult{}, cerrdefs.ErrNotFound
	}

	// check if the network is not-found
	if strings.Contains(networkID, "not-found") {
		return client.NetworkInspectResult{}, cerrdefs.ErrNotFound
	}

	// create response object to return
	response := client.NetworkInspectResult{
		Network: network.Inspect{
			Network: network.Network{
				Attachable: false,
				ConfigOnly: false,
				Created:    time.Now(),
				Driver:     "host",
				ID:         stringid.GenerateRandomID(),
				Ingress:    false,
				Internal:   false,
				Name:       networkID,
				Scope:      "local",
			},
		},
	}

	return response, nil
}

// NetworkList is a helper function to simulate
// a mocked call to list Docker networks.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NetworkList
func (n *NetworkService) NetworkList(_ context.Context, _ client.NetworkListOptions) (client.NetworkListResult, error) {
	return client.NetworkListResult{}, nil
}

// NetworkRemove is a helper function to simulate
// a mocked call to remove Docker a network.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NetworkRemove
func (n *NetworkService) NetworkRemove(_ context.Context, network string, _ client.NetworkRemoveOptions) (client.NetworkRemoveResult, error) {
	// verify a network was provided
	if len(network) == 0 {
		return client.NetworkRemoveResult{}, errors.New("no network provided")
	}

	return client.NetworkRemoveResult{}, nil
}

// NetworksPrune is a helper function to simulate
// a mocked call to prune Docker networks.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NetworksPrune
func (n *NetworkService) NetworkPrune(_ context.Context, _ client.NetworkPruneOptions) (client.NetworkPruneResult, error) {
	return client.NetworkPruneResult{}, nil
}

// NetworkConnect is a helper function to simulate
// a mocked call to connect to a Docker network.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NetworkConnect
func (n *NetworkService) NetworkConnect(_ context.Context, _ string, _ client.NetworkConnectOptions) (client.NetworkConnectResult, error) {
	return client.NetworkConnectResult{}, nil
}

// NetworkDisconnect is a helper function to simulate
// a mocked call to disconnect from a Docker network.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NetworkDisconnect
func (n *NetworkService) NetworkDisconnect(_ context.Context, _ string, _ client.NetworkDisconnectOptions) (client.NetworkDisconnectResult, error) {
	return client.NetworkDisconnectResult{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// NetworkService satisfies the NetworkAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#NetworkAPIClient
var _ client.NetworkAPIClient = (*NetworkService)(nil)
