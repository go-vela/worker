// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func TestPluginService_PluginCreate(t *testing.T) {
	service := &PluginService{}
	opts := types.PluginCreateOptions{}

	err := service.PluginCreate(context.Background(), nil, opts)

	if err != nil {
		t.Errorf("PluginCreate() error = %v, want nil", err)
	}
}

func TestPluginService_PluginDisable(t *testing.T) {
	service := &PluginService{}
	opts := types.PluginDisableOptions{}

	err := service.PluginDisable(context.Background(), "test-plugin", opts)

	if err != nil {
		t.Errorf("PluginDisable() error = %v, want nil", err)
	}
}

func TestPluginService_PluginEnable(t *testing.T) {
	service := &PluginService{}
	opts := types.PluginEnableOptions{}

	err := service.PluginEnable(context.Background(), "test-plugin", opts)

	if err != nil {
		t.Errorf("PluginEnable() error = %v, want nil", err)
	}
}

func TestPluginService_PluginInspectWithRaw(t *testing.T) {
	service := &PluginService{}

	plugin, raw, err := service.PluginInspectWithRaw(context.Background(), "test-plugin")

	if err != nil {
		t.Errorf("PluginInspectWithRaw() error = %v, want nil", err)
	}

	if plugin != nil {
		t.Errorf("PluginInspectWithRaw() plugin = %v, want nil", plugin)
	}

	if raw != nil {
		t.Errorf("PluginInspectWithRaw() raw = %v, want nil", raw)
	}
}

func TestPluginService_PluginInstall(t *testing.T) {
	service := &PluginService{}
	opts := types.PluginInstallOptions{}

	response, err := service.PluginInstall(context.Background(), "test-plugin", opts)

	if err != nil {
		t.Errorf("PluginInstall() error = %v, want nil", err)
	}

	if response != nil {
		t.Errorf("PluginInstall() response = %v, want nil", response)
	}
}

func TestPluginService_PluginList(t *testing.T) {
	service := &PluginService{}
	filters := filters.Args{}

	plugins, err := service.PluginList(context.Background(), filters)

	if err != nil {
		t.Errorf("PluginList() error = %v, want nil", err)
	}

	if len(plugins) != 0 {
		t.Errorf("PluginList() plugins = %v, want empty slice", plugins)
	}
}

func TestPluginService_PluginPush(t *testing.T) {
	service := &PluginService{}

	response, err := service.PluginPush(context.Background(), "test-plugin", "registry-auth")

	if err != nil {
		t.Errorf("PluginPush() error = %v, want nil", err)
	}

	if response != nil {
		t.Errorf("PluginPush() response = %v, want nil", response)
	}
}

func TestPluginService_PluginRemove(t *testing.T) {
	service := &PluginService{}
	opts := types.PluginRemoveOptions{}

	err := service.PluginRemove(context.Background(), "test-plugin", opts)

	if err != nil {
		t.Errorf("PluginRemove() error = %v, want nil", err)
	}
}

func TestPluginService_PluginSet(t *testing.T) {
	service := &PluginService{}
	args := []string{"key=value"}

	err := service.PluginSet(context.Background(), "test-plugin", args)

	if err != nil {
		t.Errorf("PluginSet() error = %v, want nil", err)
	}
}

func TestPluginService_PluginUpgrade(t *testing.T) {
	service := &PluginService{}
	opts := types.PluginInstallOptions{}

	response, err := service.PluginUpgrade(context.Background(), "test-plugin", opts)

	if err != nil {
		t.Errorf("PluginUpgrade() error = %v, want nil", err)
	}

	if response != nil {
		t.Errorf("PluginUpgrade() response = %v, want nil", response)
	}
}

func TestPluginService_InterfaceCompliance(t *testing.T) {
	var _ client.PluginAPIClient = (*PluginService)(nil)
}