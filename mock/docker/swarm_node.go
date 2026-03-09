// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"

	"github.com/moby/moby/client"
)

// NodeService implements all the node
// related functions for the Docker mock.
type NodeService struct{}

// NodeInspectWithRaw is a helper function to simulate
// a mocked call to inspect a node for a Docker swarm
// cluster and return the raw body received from the API.
func (n *NodeService) NodeInspect(_ context.Context, _ string, _ client.NodeInspectOptions) (client.NodeInspectResult, error) {
	return client.NodeInspectResult{}, nil
}

// NodeList is a helper function to simulate
// a mocked call to list the nodes for a
// Docker swarm cluster.
func (n *NodeService) NodeList(_ context.Context, _ client.NodeListOptions) (client.NodeListResult, error) {
	return client.NodeListResult{}, nil
}

// NodeUpdate is a helper function to simulate
// a mocked call to update a node for a
// Docker swarm cluster.
func (n *NodeService) NodeUpdate(_ context.Context, _ string, _ client.NodeUpdateOptions) (client.NodeUpdateResult, error) {
	return client.NodeUpdateResult{}, nil
}

// NodeRemove is a helper function to simulate
// a mocked call to remove a node from a
// Docker swarm cluster.
func (n *NodeService) NodeRemove(_ context.Context, _ string, _ client.NodeRemoveOptions) (client.NodeRemoveResult, error) {
	return client.NodeRemoveResult{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// service satisfies the APIClient interface that
// the Docker client expects.
var _ client.NodeAPIClient = (*NodeService)(nil)
