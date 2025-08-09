// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/swarm"
)

func TestSecretService_SecretCreate(t *testing.T) {
	s := &SecretService{}

	spec := swarm.SecretSpec{
		Annotations: swarm.Annotations{
			Name: "test-secret",
		},
		Data: []byte("secret-data"),
	}

	response, err := s.SecretCreate(context.Background(), spec)
	if err != nil {
		t.Errorf("SecretCreate() returned error: %v", err)
	}

	// Should return empty SecretCreateResponse struct
	if response.ID != "" {
		t.Errorf("SecretCreate() ID = %v, want empty string", response.ID)
	}
}

func TestSecretService_SecretInspectWithRaw(t *testing.T) {
	s := &SecretService{}

	secret, raw, err := s.SecretInspectWithRaw(context.Background(), "test-secret-id")
	if err != nil {
		t.Errorf("SecretInspectWithRaw() returned error: %v", err)
	}

	// Should return empty Secret struct and nil raw data
	if secret.ID != "" {
		t.Errorf("SecretInspectWithRaw() Secret.ID = %v, want empty string", secret.ID)
	}

	if raw != nil {
		t.Errorf("SecretInspectWithRaw() raw = %v, want nil", raw)
	}
}

func TestSecretService_SecretList(t *testing.T) {
	s := &SecretService{}

	options := swarm.SecretListOptions{}

	secrets, err := s.SecretList(context.Background(), options)
	if err != nil {
		t.Errorf("SecretList() returned error: %v", err)
	}

	// Should return nil slice
	if secrets != nil {
		t.Errorf("SecretList() = %v, want nil", secrets)
	}
}

func TestSecretService_SecretRemove(t *testing.T) {
	s := &SecretService{}

	err := s.SecretRemove(context.Background(), "test-secret-id")
	if err != nil {
		t.Errorf("SecretRemove() returned error: %v", err)
	}
}

func TestSecretService_SecretUpdate(t *testing.T) {
	s := &SecretService{}

	version := swarm.Version{Index: 1}
	spec := swarm.SecretSpec{
		Annotations: swarm.Annotations{
			Name: "test-secret-updated",
		},
		Data: []byte("updated-secret-data"),
	}

	err := s.SecretUpdate(context.Background(), "test-secret-id", version, spec)
	if err != nil {
		t.Errorf("SecretUpdate() returned error: %v", err)
	}
}
