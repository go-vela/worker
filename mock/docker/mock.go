// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"net"
	"net/http"

	"github.com/moby/moby/client"
)

type mock struct {
	// Services
	CheckpointService
	ConfigService
	ContainerService
	DistributionService
	ExecService
	ImageService
	ImageBuildService
	NetworkService
	NodeService
	PluginService
	RegistryService
	SecretService
	ServiceService
	SwarmService
	SystemService
	TaskService
	VolumeService

	// Docker API version for the mock
	Version string
}

// ClientVersion is a helper function to return
// the version string associated with the mock.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ClientVersion
func (m *mock) ClientVersion() string {
	return m.Version
}

// DaemonHost is a helper function to simulate
// returning the host address used by the client.
func (m *mock) DaemonHost() string {
	return ""
}

// ServerVersion is a helper function to simulate
// a mocked call to return information on the
// Docker client and server host.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.ServerVersion
func (m *mock) ServerVersion(_ context.Context, _ client.ServerVersionOptions) (client.ServerVersionResult, error) {
	return client.ServerVersionResult{}, nil
}

// Dialer is a helper function to simulate
// returning a dialer for the raw stream
// connection, with HTTP/1.1 header, that can
// be used for proxying the daemon connection.
func (m *mock) Dialer() func(context.Context) (net.Conn, error) {
	return func(context.Context) (net.Conn, error) { return nil, nil }
}

// DialHijack is a helper function to simulate
// returning a hijacked connection with
// negotiated protocol proto.
func (m *mock) DialHijack(_ context.Context, _ string, _ string, _ map[string][]string) (net.Conn, error) {
	return nil, nil
}

// Close is a helper function to simulate
// closing the transport client for the mock.
//
// https://pkg.go.dev/github.com/docker/docker/client#Client.Close
func (m *mock) Close() error {
	return nil
}

// HTTPClient is a helper function to simulate
// returning a copy of the HTTP client bound
// to the server.
func (m *mock) HTTPClient() *http.Client {
	return http.DefaultClient
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure
// that our mock satisfies the Go interface that the
// Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client#CommonAPIClient
var _ client.APIClient = (*mock)(nil)
