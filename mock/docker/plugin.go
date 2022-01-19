// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.PluginCreate
func (p *PluginService) PluginCreate(ctx context.Context, createContext io.Reader, options types.PluginCreateOptions) error {
	return nil
}

// PluginDisable is a helper function to simulate
// a mocked call to disable a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.PluginDisable
func (p *PluginService) PluginDisable(ctx context.Context, name string, options types.PluginDisableOptions) error {
	return nil
}

// PluginEnable is a helper function to simulate
// a mocked call to enable a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.PluginEnable
func (p *PluginService) PluginEnable(ctx context.Context, name string, options types.PluginEnableOptions) error {
	return nil
}

// PluginInspectWithRaw is a helper function to simulate
// a mocked call to inspect a Docker plugin and return
// the raw body received from the API.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.PluginInspectWithRaw
func (p *PluginService) PluginInspectWithRaw(ctx context.Context, name string) (*types.Plugin, []byte, error) {
	return nil, nil, nil
}

// PluginInstall is a helper function to simulate
// a mocked call to install a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.PluginInstall
func (p *PluginService) PluginInstall(ctx context.Context, name string, options types.PluginInstallOptions) (io.ReadCloser, error) {
	return nil, nil
}

// PluginList is a helper function to simulate
// a mocked call to list Docker plugins.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.PluginList
func (p *PluginService) PluginList(ctx context.Context, filter filters.Args) (types.PluginsListResponse, error) {
	return types.PluginsListResponse{}, nil
}

// PluginPush is a helper function to simulate
// a mocked call to push a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.PluginPush
func (p *PluginService) PluginPush(ctx context.Context, name string, registryAuth string) (io.ReadCloser, error) {
	return nil, nil
}

// PluginRemove is a helper function to simulate
// a mocked call to remove a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.PluginRemove
func (p *PluginService) PluginRemove(ctx context.Context, name string, options types.PluginRemoveOptions) error {
	return nil
}

// PluginSet is a helper function to simulate
// a mocked call to update settings for a
// Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.PluginSet
func (p *PluginService) PluginSet(ctx context.Context, name string, args []string) error {
	return nil
}

// PluginUpgrade is a helper function to simulate
// a mocked call to upgrade a Docker plugin.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.PluginUpgrade
func (p *PluginService) PluginUpgrade(ctx context.Context, name string, options types.PluginInstallOptions) (io.ReadCloser, error) {
	return nil, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure that our
// PluginService satisfies the PluginAPIClient interface that
// the Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#PluginAPIClient
var _ client.PluginAPIClient = (*PluginService)(nil)
