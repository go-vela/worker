// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

func TestConfigService_ConfigCreate(t *testing.T) {
	service := &ConfigService{}
	spec := swarm.ConfigSpec{}

	response, err := service.ConfigCreate(context.Background(), spec)
	if err != nil {
		t.Errorf("ConfigCreate() error = %v, want nil", err)
	}

	if response.ID != "" {
		t.Errorf("ConfigCreate() response.ID = %v, want empty", response.ID)
	}
}

func TestConfigService_ConfigInspectWithRaw(t *testing.T) {
	service := &ConfigService{}

	config, raw, err := service.ConfigInspectWithRaw(context.Background(), "test-config")
	if err != nil {
		t.Errorf("ConfigInspectWithRaw() error = %v, want nil", err)
	}

	if config.ID != "" {
		t.Errorf("ConfigInspectWithRaw() config.ID = %v, want empty", config.ID)
	}

	if raw != nil {
		t.Errorf("ConfigInspectWithRaw() raw = %v, want nil", raw)
	}
}

func TestConfigService_ConfigList(t *testing.T) {
	service := &ConfigService{}
	opts := swarm.ConfigListOptions{}

	configs, err := service.ConfigList(context.Background(), opts)
	if err != nil {
		t.Errorf("ConfigList() error = %v, want nil", err)
	}

	if configs != nil {
		t.Errorf("ConfigList() = %v, want nil", configs)
	}
}

func TestConfigService_ConfigRemove(t *testing.T) {
	service := &ConfigService{}

	err := service.ConfigRemove(context.Background(), "test-config")
	if err != nil {
		t.Errorf("ConfigRemove() error = %v, want nil", err)
	}
}

func TestConfigService_ConfigUpdate(t *testing.T) {
	service := &ConfigService{}
	version := swarm.Version{}
	spec := swarm.ConfigSpec{}

	err := service.ConfigUpdate(context.Background(), "test-config", version, spec)
	if err != nil {
		t.Errorf("ConfigUpdate() error = %v, want nil", err)
	}
}

func TestConfigService_InterfaceCompliance(_ *testing.T) {
	var _ client.ConfigAPIClient = (*ConfigService)(nil)
}
