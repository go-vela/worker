// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/swarm"
)

func TestNodeService_NodeInspectWithRaw(t *testing.T) {
	n := &NodeService{}

	node, raw, err := n.NodeInspectWithRaw(context.Background(), "test-node-id")
	if err != nil {
		t.Errorf("NodeInspectWithRaw() returned error: %v", err)
	}

	// Should return empty Node struct and nil raw data
	if node.ID != "" {
		t.Errorf("NodeInspectWithRaw() Node.ID = %v, want empty string", node.ID)
	}

	if raw != nil {
		t.Errorf("NodeInspectWithRaw() raw = %v, want nil", raw)
	}
}

func TestNodeService_NodeList(t *testing.T) {
	n := &NodeService{}

	options := swarm.NodeListOptions{}

	nodes, err := n.NodeList(context.Background(), options)
	if err != nil {
		t.Errorf("NodeList() returned error: %v", err)
	}

	// Should return nil slice
	if nodes != nil {
		t.Errorf("NodeList() = %v, want nil", nodes)
	}
}

func TestNodeService_NodeRemove(t *testing.T) {
	n := &NodeService{}

	options := swarm.NodeRemoveOptions{
		Force: true,
	}

	err := n.NodeRemove(context.Background(), "test-node-id", options)
	if err != nil {
		t.Errorf("NodeRemove() returned error: %v", err)
	}
}

func TestNodeService_NodeUpdate(t *testing.T) {
	n := &NodeService{}

	version := swarm.Version{Index: 1}
	spec := swarm.NodeSpec{
		Role: swarm.NodeRoleWorker,
	}

	err := n.NodeUpdate(context.Background(), "test-node-id", version, spec)
	if err != nil {
		t.Errorf("NodeUpdate() returned error: %v", err)
	}
}
