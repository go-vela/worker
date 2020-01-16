// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"context"
	"testing"

	"github.com/go-vela/types/pipeline"
)

func TestDocker_CreateVolume_Success(t *testing.T) {
	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.CreateVolume(context.Background(), &pipeline.Build{
		Version: "1",
		ID:      "__0"})

	if got != nil {
		t.Error("CreateVolume should not have returned err: ", got)
	}

	if got != nil {
		t.Errorf("CreateVolume is %v, want nil", got)
	}
}

// TODO: rethink how the mock is being done in the
// router switch. This current gives false positives
func TestDocker_CreateVolume_Failure(t *testing.T) {
	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.CreateVolume(context.Background(), &pipeline.Build{
		Version: "1",
		ID:      "__0"})

	// this should be "=="
	if got != nil {
		t.Errorf("CreateVolume should have returned err: %+v", got)
	}
}

func TestDocker_InspectVolume_Success(t *testing.T) {
	// setup types
	p := &pipeline.Build{
		Version: "1",
		ID:      "__0",
	}

	// setup Docker
	c, _ := NewMock()

	// run test
	got, err := c.InspectVolume(context.Background(), p)
	if err != nil {
		t.Errorf("InspectVolume returned err: %v", got)
	}

	if got == nil {
		t.Errorf("InspectVolume is nil, want %v", got)
	}
}

func TestDocker_InspectVolume_Failure(t *testing.T) {
	// setup types
	p := &pipeline.Build{
		Version: "1",
		ID:      "notfound",
	}

	// setup Docker
	c, _ := NewMock()

	// run test
	got, err := c.InspectVolume(context.Background(), p)
	if err == nil {
		t.Errorf("InspectVolume should have returned err")
	}

	if got != nil {
		t.Errorf("InspectVolume is %v, want nil", got)
	}
}

func TestDocker_RemoveVolume_Success(t *testing.T) {
	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.RemoveVolume(context.Background(), &pipeline.Build{
		Version: "1",
		ID:      "__0"})

	if got != nil {
		t.Error("RemoveVolume should not have returned err: ", got)
	}

	if got != nil {
		t.Errorf("RemoveVolume is %v, want nil", got)
	}
}

func TestDocker_RemoveVolume_Failure(t *testing.T) {
	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.RemoveVolume(context.Background(), &pipeline.Build{})

	if got == nil {
		t.Errorf("RemoveVolume should have returned err: %+v", got)
	}
}
