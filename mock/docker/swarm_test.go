// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

func TestSwarmService_SwarmGetUnlockKey(t *testing.T) {
	service := &SwarmService{}

	response, err := service.SwarmGetUnlockKey(context.Background())
	if err != nil {
		t.Errorf("SwarmGetUnlockKey() error = %v, want nil", err)
	}

	if response.UnlockKey != "" {
		t.Errorf("SwarmGetUnlockKey() response.UnlockKey = %v, want empty", response.UnlockKey)
	}
}

func TestSwarmService_SwarmInit(t *testing.T) {
	service := &SwarmService{}
	request := swarm.InitRequest{}

	nodeID, err := service.SwarmInit(context.Background(), request)
	if err != nil {
		t.Errorf("SwarmInit() error = %v, want nil", err)
	}

	if nodeID != "" {
		t.Errorf("SwarmInit() nodeID = %v, want empty", nodeID)
	}
}

func TestSwarmService_SwarmInspect(t *testing.T) {
	service := &SwarmService{}

	swarmInfo, err := service.SwarmInspect(context.Background())
	if err != nil {
		t.Errorf("SwarmInspect() error = %v, want nil", err)
	}

	if swarmInfo.ID != "" {
		t.Errorf("SwarmInspect() swarmInfo.ID = %v, want empty", swarmInfo.ID)
	}
}

func TestSwarmService_SwarmJoin(t *testing.T) {
	service := &SwarmService{}
	request := swarm.JoinRequest{}

	err := service.SwarmJoin(context.Background(), request)
	if err != nil {
		t.Errorf("SwarmJoin() error = %v, want nil", err)
	}
}

func TestSwarmService_SwarmLeave(t *testing.T) {
	service := &SwarmService{}

	err := service.SwarmLeave(context.Background(), false)
	if err != nil {
		t.Errorf("SwarmLeave() error = %v, want nil", err)
	}
}

func TestSwarmService_SwarmUnlock(t *testing.T) {
	service := &SwarmService{}
	request := swarm.UnlockRequest{}

	err := service.SwarmUnlock(context.Background(), request)
	if err != nil {
		t.Errorf("SwarmUnlock() error = %v, want nil", err)
	}
}

func TestSwarmService_SwarmUpdate(t *testing.T) {
	service := &SwarmService{}
	version := swarm.Version{}
	spec := swarm.Spec{}
	flags := swarm.UpdateFlags{}

	err := service.SwarmUpdate(context.Background(), version, spec, flags)
	if err != nil {
		t.Errorf("SwarmUpdate() error = %v, want nil", err)
	}
}

func TestSwarmService_InterfaceCompliance(_ *testing.T) {
	var _ client.SwarmAPIClient = (*SwarmService)(nil)
}
