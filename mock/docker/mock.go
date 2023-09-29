// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"net"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type mock struct {
	// Services
	ConfigService
	ContainerService
	DistributionService
	ImageService
	NetworkService
	NodeService
	PluginService
	SecretService
	ServiceService
	SwarmService
	SystemService
	VolumeService

	// Docker API version for the mock
	Version string
}

// ClientVersion is a helper function to return
// the version string associated with the mock.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ClientVersion
func (m *mock) ClientVersion() string {
	return m.Version
}

// Close is a helper function to simulate
// closing the transport client for the mock.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.Close
func (m *mock) Close() error {
	return nil
}

// DaemonHost is a helper function to simulate
// returning the host address used by the client.
func (m *mock) DaemonHost() string {
	return ""
}

// DialSession is a helper function to simulate
// returning a connection that can be used
// for communication with daemon.
func (m *mock) DialSession(ctx context.Context, proto string, meta map[string][]string) (net.Conn, error) {
	return nil, nil
}

// DialHijack is a helper function to simulate
// returning a hijacked connection with
// negotiated protocol proto.
func (m *mock) DialHijack(ctx context.Context, url, proto string, meta map[string][]string) (net.Conn, error) {
	return nil, nil
}

// Dialer is a helper function to simulate
// returning a dialer for the raw stream
// connection, with HTTP/1.1 header, that can
// be used for proxying the daemon connection.
func (m *mock) Dialer() func(context.Context) (net.Conn, error) {
	return func(context.Context) (net.Conn, error) { return nil, nil }
}

// HTTPClient is a helper function to simulate
// returning a copy of the HTTP client bound
// to the server.
func (m *mock) HTTPClient() *http.Client {
	return http.DefaultClient
}

// NegotiateAPIVersion is a helper function to simulate
// a mocked call to query the API and update the client
// version to match the API version.
func (m *mock) NegotiateAPIVersion(ctx context.Context) {}

// NegotiateAPIVersionPing is a helper function to simulate
// a mocked call to update the client version to match
// the ping version if it's less than the default version.
func (m *mock) NegotiateAPIVersionPing(types.Ping) {}

// ServerVersion is a helper function to simulate
// a mocked call to return information on the
// Docker client and server host.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#Client.ServerVersion
func (m *mock) ServerVersion(ctx context.Context) (types.Version, error) {
	return types.Version{}, nil
}

// WARNING: DO NOT REMOVE THIS UNDER ANY CIRCUMSTANCES
//
// This line serves as a quick and efficient way to ensure
// that our mock satisfies the Go interface that the
// Docker client expects.
//
// https://pkg.go.dev/github.com/docker/docker/client?tab=doc#CommonAPIClient
var _ client.CommonAPIClient = (*mock)(nil)
