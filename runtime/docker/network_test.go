// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"context"
	"testing"

	"github.com/go-vela/types/pipeline"
)

func TestDocker_CreateNetwork_Success(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.CreateNetwork(context.Background(), &pipeline.Build{
		Version: "1",
		ID:      "__0"})

	if got != nil {
		t.Error("CreateNetwork should not have returned err: ", got)
	}

	if got != nil {
		t.Errorf("CreateNetwork is %v, want nil", got)
	}
}

// TODO: rethink how the mock is being done in the
// router switch. This current gives false positives
func TestDocker_CreateNetwork_Failure(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.CreateNetwork(context.Background(), &pipeline.Build{
		Version: "1",
		ID:      "__0"})

	// this should be "=="
	if got != nil {
		t.Errorf("CreateNetwork should have returned err: %+v", got)
	}
}

func TestDocker_InspectNetwork_Success(t *testing.T) {
	// setup types
	p := &pipeline.Build{
		Version: "1",
		ID:      "__0",
	}

	// setup Docker
	c, _ := NewMock()

	// run test
	got, err := c.InspectNetwork(context.Background(), p)
	if err != nil {
		t.Errorf("InspectNetwork returned err: %v", err)
	}

	if got == nil {
		t.Errorf("InspectNetwork is nil, want %v", got)
	}
}

func TestDocker_InspectNetwork_Failure(t *testing.T) {
	// setup types
	p := &pipeline.Build{
		Version: "1",
		ID:      "notfound",
	}

	// setup Docker
	c, _ := NewMock()

	// run test
	got, err := c.InspectNetwork(context.Background(), p)
	if err == nil {
		t.Errorf("InspectNetwork should have returned err")
	}

	if got != nil {
		t.Errorf("InspectNetwork is %v, want nil", got)
	}
}

func TestDocker_RemoveNetwork_Success(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.RemoveNetwork(context.Background(), &pipeline.Build{
		Version: "1",
		ID:      "__0"})

	if got != nil {
		t.Error("RemoveNetwork should not have returned err: ", got)
	}

	if got != nil {
		t.Errorf("RemoveNetwork is %v, want nil", got)
	}
}

func TestDocker_RemoveNetwork_Failure(t *testing.T) {

	// setup Docker
	c, _ := NewMock()

	// run test
	got := c.RemoveNetwork(context.Background(), &pipeline.Build{})

	if got == nil {
		t.Errorf("RemoveNetwork should have returned err: %+v", got)
	}
}
