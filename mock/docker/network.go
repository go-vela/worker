// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

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
// https://pkg.go.dev/github.com/docker/docker/client#Client.NetworkConnect
func (n *NetworkService) NetworkConnect(_ context.Context, _, _ string, _ *network.EndpointSettings) error {
	return nil
}

// NetworkCreate is a helper function to simulate
// a mocked call to create a Docker network.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NetworkCreate
func (n *NetworkService) NetworkCreate(_ context.Context, name string, _ network.CreateOptions) (network.CreateResponse, error) {
	// verify a network was provided
	if len(name) == 0 {
		return network.CreateResponse{}, errors.New("no network provided")
	}

	// check if the network is notfound and
	// check if the notfound should be ignored
	if strings.Contains(name, "notfound") &&
		!strings.Contains(name, "ignorenotfound") {
		return network.CreateResponse{},
			errdefs.NotFound(fmt.Errorf("error: no such network: %s", name))
	}

	// check if the network is not-found and
	// check if the not-found should be ignored
	if strings.Contains(name, "not-found") &&
		!strings.Contains(name, "ignore-not-found") {
		return network.CreateResponse{},
			errdefs.NotFound(fmt.Errorf("error: no such network: %s", name))
	}

	// create response object to return
	response := network.CreateResponse{
		ID: stringid.GenerateRandomID(),
	}

	return response, nil
}

// NetworkDisconnect is a helper function to simulate
// a mocked call to disconnect from a Docker network.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NetworkDisconnect
func (n *NetworkService) NetworkDisconnect(_ context.Context, _, _ string, _ bool) error {
	return nil
}

// NetworkInspect is a helper function to simulate
// a mocked call to inspect a Docker network.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NetworkInspect
func (n *NetworkService) NetworkInspect(ctx context.Context, networkID string, options network.InspectOptions) (network.Inspect, error) {
	// verify a network was provided
	if len(networkID) == 0 {
		return network.Inspect{}, errors.New("no network provided")
	}

	// check if the network is notfound
	if strings.Contains(networkID, "notfound") {
		return network.Inspect{},
			errdefs.NotFound(fmt.Errorf("error: no such network: %s", networkID))
	}

	// check if the network is not-found
	if strings.Contains(networkID, "not-found") {
		return network.Inspect{},
			errdefs.NotFound(fmt.Errorf("error: no such network: %s", networkID))
	}

	// create response object to return
	response := network.Inspect{
		Attachable: false,
		ConfigOnly: false,
		Created:    time.Now(),
		Driver:     "host",
		ID:         stringid.GenerateRandomID(),
		Ingress:    false,
		Internal:   false,
		Name:       networkID,
		Scope:      "local",
	}

	return response, nil
}

// NetworkInspectWithRaw is a helper function to simulate
// a mocked call to inspect a Docker network and return
// the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NetworkInspectWithRaw
func (n *NetworkService) NetworkInspectWithRaw(ctx context.Context, networkID string, options network.InspectOptions) (network.Inspect, []byte, error) {
	// verify a network was provided
	if len(networkID) == 0 {
		return network.Inspect{}, nil, errors.New("no network provided")
	}

	// create response object to return
	response := network.Inspect{
		Attachable: false,
		ConfigOnly: false,
		Created:    time.Now(),
		Driver:     "host",
		ID:         stringid.GenerateRandomID(),
		Ingress:    false,
		Internal:   false,
		Name:       networkID,
		Scope:      "local",
	}

	// marshal response into raw bytes
	b, err := json.Marshal(response)
	if err != nil {
		return network.Inspect{}, nil, err
	}

	return response, b, nil
}

// NetworkList is a helper function to simulate
// a mocked call to list Docker networks.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NetworkList
func (n *NetworkService) NetworkList(ctx context.Context, options network.ListOptions) ([]network.Summary, error) {
	return nil, nil
}

// NetworkRemove is a helper function to simulate
// a mocked call to remove Docker a network.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NetworkRemove
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
// https://pkg.go.dev/github.com/docker/docker/client#Client.NetworksPrune
func (n *NetworkService) NetworksPrune(ctx context.Context, pruneFilter filters.Args) (network.PruneReport, error) {
	return network.PruneReport{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// NetworkService satisfies the NetworkAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#NetworkAPIClient
var _ client.NetworkAPIClient = (*NetworkService)(nil)
