// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"context"
	"testing"

	"github.com/go-vela/types/pipeline"
)

func TestDocker_InfoContainer_Success(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.InfoContainer(context.Background(), &pipeline.Container{
		ID:    "container_id",
		Image: "alpine:latest",
	})

	if got != nil {
		t.Error("InfoContainer should not have returned err: ", got)
	}

	if got != nil {
		t.Errorf("InfoContainer is %v, want nil", got)
	}
}

func TestDocker_InfoContainer_Failure(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.InfoContainer(context.Background(), &pipeline.Container{})

	if got == nil {
		t.Errorf("InfoContainer should have returned err: %+v", got)
	}
}

func TestDocker_RemoveContainer_Success(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.RemoveContainer(context.Background(), &pipeline.Container{
		ID:    "container_id",
		Image: "alpine:latest",
	})

	if got != nil {
		t.Error("RemoveContainer should not have returned err: ", got)
	}

	if got != nil {
		t.Errorf("RemoveContainer is %v, want nil", got)
	}
}

func TestDocker_RemoveContainer_Failure(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.RemoveContainer(context.Background(), &pipeline.Container{})

	if got == nil {
		t.Errorf("RemoveContainer should have returned err: %+v", got)
	}
}

func TestDocker_RunContainer_Success(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.RunContainer(context.Background(),
		&pipeline.Build{
			Version: "1",
			ID:      "__0",
		},
		&pipeline.Container{
			ID:    "container_id",
			Image: "alpine:latest",
		})

	if got != nil {
		t.Error("RunContainer should not have returned err: ", got)
	}

	if got != nil {
		t.Errorf("RunContainer is %v, want nil", got)
	}
}

// TODO: rethink how the mock is being done in the
// router switch. This current gives false positives
func TestDocker_RunContainer_Failure(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.RunContainer(context.Background(),
		&pipeline.Build{},
		&pipeline.Container{})

	// this should be "=="
	if got != nil {
		t.Errorf("RunContainer should have returned err: %+v", got)
	}
}

func TestDocker_SetupContainer_Success(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.SetupContainer(context.Background(), &pipeline.Container{
		ID:    "container_id",
		Image: "alpine:latest",
		Pull:  true,
	})

	if got != nil {
		t.Errorf("SetupContainer is %v, want nil", got)
	}
}

func TestDocker_SetupContainer_Failure(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.SetupContainer(context.Background(), &pipeline.Container{})

	if got == nil {
		t.Errorf("SetupContainer should have returned err: %+v", got)
	}
}

func TestDocker_TailContainer_Success(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	_, got := c.TailContainer(context.Background(), &pipeline.Container{
		ID:    "container_id",
		Image: "alpine:latest",
	})

	if got != nil {
		t.Error("TailContainer should not have returned err: ", got)
	}

	if got != nil {
		t.Errorf("TailContainer is %v, want nil", got)
	}
}

// TODO: rethink how the mock is being done in the
// router switch. This current gives false positives
func TestDocker_TailContainer_Failure(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	_, got := c.TailContainer(context.Background(), &pipeline.Container{})

	// this should be "=="
	if got != nil {
		t.Errorf("TailContainer should have returned err: %+v", got)
	}
}

func TestDocker_WaitContainer_Success(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.WaitContainer(context.Background(), &pipeline.Container{
		ID:    "container_id",
		Image: "alpine:latest",
	})

	if got != nil {
		t.Error("WaitContainer should not have returned err: ", got)
	}

	if got != nil {
		t.Errorf("WaitContainer is %v, want nil", got)
	}
}

func TestDocker_WaitContainer_Failure(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.WaitContainer(context.Background(), &pipeline.Container{})

	if got == nil {
		t.Errorf("WaitContainer should have returned err: %+v", got)
	}
}
