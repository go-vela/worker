// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/registry"
)

func TestSystemService_DiskUsage(t *testing.T) {
	s := &SystemService{}

	usage, err := s.DiskUsage(context.Background(), types.DiskUsageOptions{})
	if err != nil {
		t.Errorf("DiskUsage() returned error: %v", err)
	}

	// Should return empty DiskUsage struct
	if usage.LayersSize != 0 {
		t.Errorf("DiskUsage() LayersSize = %v, want 0", usage.LayersSize)
	}
}

func TestSystemService_Events(t *testing.T) {
	s := &SystemService{}

	msgChan, errChan := s.Events(context.Background(), events.ListOptions{})

	// Should return nil channels for mock
	if msgChan != nil {
		t.Errorf("Events() message channel = %v, want nil", msgChan)
	}

	if errChan != nil {
		t.Errorf("Events() error channel = %v, want nil", errChan)
	}
}

func TestSystemService_Info(t *testing.T) {
	s := &SystemService{}

	info, err := s.Info(context.Background())
	if err != nil {
		t.Errorf("Info() returned error: %v", err)
	}

	// Should return empty Info struct
	if info.ID != "" {
		t.Errorf("Info() ID = %v, want empty string", info.ID)
	}
}

func TestSystemService_Ping(t *testing.T) {
	s := &SystemService{}

	ping, err := s.Ping(context.Background())
	if err != nil {
		t.Errorf("Ping() returned error: %v", err)
	}

	// Should return empty Ping struct
	if ping.APIVersion != "" {
		t.Errorf("Ping() APIVersion = %v, want empty string", ping.APIVersion)
	}

	if ping.OSType != "" {
		t.Errorf("Ping() OSType = %v, want empty string", ping.OSType)
	}
}

func TestSystemService_RegistryLogin(t *testing.T) {
	s := &SystemService{}

	authConfig := registry.AuthConfig{
		Username: "test",
		Password: "password",
	}

	auth, err := s.RegistryLogin(context.Background(), authConfig)
	if err != nil {
		t.Errorf("RegistryLogin() returned error: %v", err)
	}

	// Should return empty AuthenticateOKBody struct
	if auth.Status != "" {
		t.Errorf("RegistryLogin() Status = %v, want empty string", auth.Status)
	}

	if auth.IdentityToken != "" {
		t.Errorf("RegistryLogin() IdentityToken = %v, want empty string", auth.IdentityToken)
	}
}
