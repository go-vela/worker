// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// NodeService implements all the node
// related functions for the Docker mock.
type NodeService struct{}

// NodeInspectWithRaw is a helper function to simulate
// a mocked call to inspect a node for a Docker swarm
// cluster and return the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NodeInspectWithRaw
func (n *NodeService) NodeInspectWithRaw(_ context.Context, _ string) (swarm.Node, []byte, error) {
	return swarm.Node{}, nil, nil
}

// NodeList is a helper function to simulate
// a mocked call to list the nodes for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NodeList
func (n *NodeService) NodeList(_ context.Context, _ swarm.NodeListOptions) ([]swarm.Node, error) {
	return nil, nil
}

// NodeRemove is a helper function to simulate
// a mocked call to remove a node from a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NodeRemove
func (n *NodeService) NodeRemove(_ context.Context, _ string, _ swarm.NodeRemoveOptions) error {
	return nil
}

// NodeUpdate is a helper function to simulate
// a mocked call to update a node for a
// Docker swarm cluster.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.NodeUpdate
func (n *NodeService) NodeUpdate(_ context.Context, _ string, _ swarm.Version, _ swarm.NodeSpec) error {
	return nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// NodeService satisfies the NodeAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#NodeAPIClient
var _ client.NodeAPIClient = (*NodeService)(nil)
