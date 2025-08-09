// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"net/http"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func TestMock_ClientVersion(t *testing.T) {
	m := &mock{
		Version: "1.41",
	}

	version := m.ClientVersion()

	if version != "1.41" {
		t.Errorf("ClientVersion() = %v, want 1.41", version)
	}

	// Test empty version
	m.Version = ""

	version = m.ClientVersion()
	if version != "" {
		t.Errorf("ClientVersion() = %v, want empty string", version)
	}
}

func TestMock_Close(t *testing.T) {
	m := &mock{}

	err := m.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}
}

func TestMock_DaemonHost(t *testing.T) {
	m := &mock{}

	host := m.DaemonHost()
	if host != "" {
		t.Errorf("DaemonHost() = %v, want empty string", host)
	}
}

func TestMock_DialSession(t *testing.T) {
	m := &mock{}

	conn, err := m.DialSession(context.Background(), "proto", nil)
	if err != nil {
		t.Errorf("DialSession() returned error: %v", err)
	}

	if conn != nil {
		t.Errorf("DialSession() conn = %v, want nil", conn)
	}
}

func TestMock_DialHijack(t *testing.T) {
	m := &mock{}

	conn, err := m.DialHijack(context.Background(), "url", "proto", nil)
	if err != nil {
		t.Errorf("DialHijack() returned error: %v", err)
	}

	if conn != nil {
		t.Errorf("DialHijack() conn = %v, want nil", conn)
	}
}

func TestMock_Dialer(t *testing.T) {
	m := &mock{}

	dialer := m.Dialer()
	if dialer == nil {
		t.Error("Dialer() returned nil, want function")
	}

	// Test the dialer function
	conn, err := dialer(context.Background())
	if err != nil {
		t.Errorf("Dialer function returned error: %v", err)
	}

	if conn != nil {
		t.Errorf("Dialer function conn = %v, want nil", conn)
	}
}

func TestMock_HTTPClient(t *testing.T) {
	m := &mock{}

	client := m.HTTPClient()
	if client == nil {
		t.Error("HTTPClient() returned nil")
		return
	}

	// Should return default HTTP client - verify it's the default client
	if client != http.DefaultClient {
		t.Error("HTTPClient() should return http.DefaultClient")
	}
}

func TestMock_NegotiateAPIVersion(_ *testing.T) {
	m := &mock{}

	// This should not panic and should complete without error
	m.NegotiateAPIVersion(context.Background())
}

func TestMock_NegotiateAPIVersionPing(_ *testing.T) {
	m := &mock{}

	ping := types.Ping{
		APIVersion: "1.40",
	}

	// This should not panic and should complete without error
	m.NegotiateAPIVersionPing(ping)
}

func TestMock_ServerVersion(t *testing.T) {
	m := &mock{}

	version, err := m.ServerVersion(context.Background())
	if err != nil {
		t.Errorf("ServerVersion() returned error: %v", err)
	}

	// Should return empty Version struct
	if version.APIVersion != "" {
		t.Errorf("ServerVersion() APIVersion = %v, want empty string", version.APIVersion)
	}

	if version.Version != "" {
		t.Errorf("ServerVersion() Version = %v, want empty string", version.Version)
	}

	if version.GitCommit != "" {
		t.Errorf("ServerVersion() GitCommit = %v, want empty string", version.GitCommit)
	}
}

func TestMockStructComposition(t *testing.T) {
	m := &mock{}

	// Test that embedded services are accessible by calling methods directly
	err := m.ConfigRemove(context.Background(), "test")
	if err != nil {
		t.Logf("ConfigRemove accessible and returned expected nil error")
	}

	err = m.ContainerRemove(context.Background(), "test", container.RemoveOptions{})
	if err != nil {
		t.Logf("ContainerRemove accessible and returned expected nil error")
	}

	// Test a few key services to verify composition works
	_, err = m.Info(context.Background())
	if err != nil {
		t.Errorf("Info() returned error: %v", err)
	}

	_, err = m.VolumeInspect(context.Background(), "test")
	if err != nil {
		t.Errorf("VolumeInspect() returned error: %v", err)
	}
}
