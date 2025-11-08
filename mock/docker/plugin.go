// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// PluginService implements all the plugin
// related functions for the Docker mock.
type PluginService struct{}

// PluginCreate is a helper function to simulate
// a mocked call to create a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginCreate
func (p *PluginService) PluginCreate(_ context.Context, _ io.Reader, _ types.PluginCreateOptions) error {
	return nil
}

// PluginDisable is a helper function to simulate
// a mocked call to disable a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginDisable
func (p *PluginService) PluginDisable(_ context.Context, _ string, _ types.PluginDisableOptions) error {
	return nil
}

// PluginEnable is a helper function to simulate
// a mocked call to enable a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginEnable
func (p *PluginService) PluginEnable(_ context.Context, _ string, _ types.PluginEnableOptions) error {
	return nil
}

// PluginInspectWithRaw is a helper function to simulate
// a mocked call to inspect a Docker plugin and return
// the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginInspectWithRaw
func (p *PluginService) PluginInspectWithRaw(_ context.Context, _ string) (*types.Plugin, []byte, error) {
	return nil, nil, nil
}

// PluginInstall is a helper function to simulate
// a mocked call to install a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginInstall
func (p *PluginService) PluginInstall(_ context.Context, _ string, _ types.PluginInstallOptions) (io.ReadCloser, error) {
	return nil, nil
}

// PluginList is a helper function to simulate
// a mocked call to list Docker plugins.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginList
func (p *PluginService) PluginList(_ context.Context, _ filters.Args) (types.PluginsListResponse, error) {
	return types.PluginsListResponse{}, nil
}

// PluginPush is a helper function to simulate
// a mocked call to push a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginPush
func (p *PluginService) PluginPush(_ context.Context, _ string, _ string) (io.ReadCloser, error) {
	return nil, nil
}

// PluginRemove is a helper function to simulate
// a mocked call to remove a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginRemove
func (p *PluginService) PluginRemove(_ context.Context, _ string, _ types.PluginRemoveOptions) error {
	return nil
}

// PluginSet is a helper function to simulate
// a mocked call to update settings for a
// Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginSet
func (p *PluginService) PluginSet(_ context.Context, _ string, _ []string) error {
	return nil
}

// PluginUpgrade is a helper function to simulate
// a mocked call to upgrade a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginUpgrade
func (p *PluginService) PluginUpgrade(_ context.Context, _ string, _ types.PluginInstallOptions) (io.ReadCloser, error) {
	return nil, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// PluginService satisfies the PluginAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#PluginAPIClient
var _ client.PluginAPIClient = (*PluginService)(nil)
