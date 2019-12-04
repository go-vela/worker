// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package runtime

import (
	"context"
	"io"

	"github.com/go-vela/types/pipeline"
)

// Engine represents the interface for Vela integrating
// with the different supported Runtime environments.
type Engine interface {

	// Container Engine Interface Functions

	// InspectContainer defines a function that inspects
	// the pipeline container.
	InspectContainer(context.Context, *pipeline.Container) error
	// RemoveContainer defines a function that deletes
	// (kill, remove) the pipeline container.
	RemoveContainer(context.Context, *pipeline.Container) error
	// RunContainer defines a function that creates
	// and start the pipeline container.
	RunContainer(context.Context, *pipeline.Build, *pipeline.Container) error
	// SetupContainer defines a function that pulls
	// the image for the pipeline container.
	SetupContainer(context.Context, *pipeline.Container) error
	// TailContainer defines a function that captures
	// the logs on the pipeline container.
	TailContainer(context.Context, *pipeline.Container) (io.ReadCloser, error)
	// WaitContainer defines a function that blocks
	// until the pipeline container completes.
	WaitContainer(context.Context, *pipeline.Container) error

	// Network Engine Interface Functions

	// CreateNetwork defines a function that
	// creates the pipeline network.
	CreateNetwork(context.Context, *pipeline.Build) error
	// InspectNetwork defines a function that
	// inspects the pipeline network.
	InspectNetwork(context.Context, *pipeline.Build) ([]byte, error)
	// RemoveNetwork defines a function that
	// deletes the pipeline network.
	RemoveNetwork(context.Context, *pipeline.Build) error

	// Volume Engine Interface Functions

	// CreateVolume defines a function that
	// creates the pipeline volume.
	CreateVolume(context.Context, *pipeline.Build) error
	// InspectVolume defines a function that
	// inspects the pipeline volume.
	InspectVolume(context.Context, *pipeline.Build) ([]byte, error)
	// RemoveVolume defines a function that
	// deletes the pipeline volume.
	RemoveVolume(context.Context, *pipeline.Build) error
}
