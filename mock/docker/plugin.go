// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"io"

	"github.com/moby/moby/client"
)

// PluginService implements all the plugin
// related functions for the Docker mock.
type PluginService struct{}

// PluginCreate is a helper function to simulate
// a mocked call to create a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginCreate
func (p *PluginService) PluginCreate(_ context.Context, _ io.Reader, _ client.PluginCreateOptions) (client.PluginCreateResult, error) {
	return client.PluginCreateResult{}, nil
}

// PluginInstall is a helper function to simulate
// a mocked call to install a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginInstall
func (p *PluginService) PluginInstall(_ context.Context, _ string, _ client.PluginInstallOptions) (client.PluginInstallResult, error) {
	return client.PluginInstallResult{}, nil
}

// PluginInspectWithRaw is a helper function to simulate
// a mocked call to inspect a Docker plugin and return
// the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginInspectWithRaw
func (p *PluginService) PluginInspect(_ context.Context, _ string, _ client.PluginInspectOptions) (client.PluginInspectResult, error) {
	return client.PluginInspectResult{}, nil
}

// PluginList is a helper function to simulate
// a mocked call to list Docker plugins.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginList
func (p *PluginService) PluginList(_ context.Context, _ client.PluginListOptions) (client.PluginListResult, error) {
	return client.PluginListResult{}, nil
}

// PluginRemove is a helper function to simulate
// a mocked call to remove a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginRemove
func (p *PluginService) PluginRemove(_ context.Context, _ string, _ client.PluginRemoveOptions) (client.PluginRemoveResult, error) {
	return client.PluginRemoveResult{}, nil
}

// PluginEnable is a helper function to simulate
// a mocked call to enable a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginEnable
func (p *PluginService) PluginEnable(_ context.Context, _ string, _ client.PluginEnableOptions) (client.PluginEnableResult, error) {
	return client.PluginEnableResult{}, nil
}

// PluginDisable is a helper function to simulate
// a mocked call to disable a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginDisable
func (p *PluginService) PluginDisable(_ context.Context, _ string, _ client.PluginDisableOptions) (client.PluginDisableResult, error) {
	return client.PluginDisableResult{}, nil
}

// PluginUpgrade is a helper function to simulate
// a mocked call to upgrade a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginUpgrade
func (p *PluginService) PluginUpgrade(_ context.Context, _ string, _ client.PluginUpgradeOptions) (client.PluginUpgradeResult, error) {
	return nil, nil
}

// PluginPush is a helper function to simulate
// a mocked call to push a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginPush
func (p *PluginService) PluginPush(_ context.Context, _ string, _ client.PluginPushOptions) (client.PluginPushResult, error) {
	return client.PluginPushResult{}, nil
}

// PluginSet is a helper function to simulate
// a mocked call to update settings for a
// Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.PluginSet
func (p *PluginService) PluginSet(_ context.Context, _ string, _ client.PluginSetOptions) (client.PluginSetResult, error) {
	return client.PluginSetResult{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// PluginService satisfies the PluginAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#PluginAPIClient
var _ client.PluginAPIClient = (*PluginService)(nil)
