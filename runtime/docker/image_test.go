// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"context"
	"testing"

	"github.com/go-vela/types/pipeline"
)

func TestDocker_InspectImage(t *testing.T) {
	// setup types
	p := &pipeline.Container{
		ID:    "container_id",
		Image: "alpine:latest",
	}

	// setup Docker
	c, _ := NewMock()

	// run test
	got, err := c.InspectImage(context.Background(), p)

	if err != nil {
		t.Errorf("InspectImage returned err: %v", err)
	}

	if got == nil {
		t.Errorf("InspectImage is nil, want %v", got)
	}
}

func TestDocker_InspectImage_BadImage(t *testing.T) {
	// setup types
	p := &pipeline.Container{
		ID:    "container_id",
		Image: "alpine:notfound",
	}

	// setup Docker
	c, _ := NewMock()

	// run test
	got, err := c.InspectImage(context.Background(), p)

	if err == nil {
		t.Errorf("InspectImage should have returned err")
	}

	if got != nil {
		t.Errorf("InspectImage is %v, want nil", got)
	}
}

func TestDocker_InspectImage_NoImage(t *testing.T) {
	// setup types
	p := new(pipeline.Container)

	// setup Docker
	c, _ := NewMock()

	// run test
	got, err := c.InspectImage(context.Background(), p)

	if err == nil {
		t.Errorf("InspectImage should have returned err")
	}

	if got != nil {
		t.Errorf("InspectImage is %v, want nil", got)
	}
}
