// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
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

	// Engine Interface Functions

	// Driver defines a function that outputs
	// the configured runtime driver.
	Driver() string

	// Build Engine Interface Functions

	// InspectBuild defines a function that
	// displays details about the build for the init step.
	InspectBuild(ctx context.Context, b *pipeline.Build) ([]byte, error)
	// SetupBuild defines a function that
	// prepares the pipeline build.
	SetupBuild(context.Context, *pipeline.Build) error
	// AssembleBuild defines a function that
	// finalizes pipeline build setup.
	AssembleBuild(context.Context, *pipeline.Build) error
	// RemoveBuild defines a function that deletes
	// (kill, remove) the pipeline build metadata.
	RemoveBuild(context.Context, *pipeline.Build) error

	// Container Engine Interface Functions

	// InspectContainer defines a function that inspects
	// the pipeline container.
	InspectContainer(context.Context, *pipeline.Container) error
	// RemoveContainer defines a function that deletes
	// (kill, remove) the pipeline container.
	RemoveContainer(context.Context, *pipeline.Container) error
	// RunContainer defines a function that creates
	// and starts the pipeline container.
	RunContainer(context.Context, *pipeline.Container, *pipeline.Build) error
	// SetupContainer defines a function that prepares
	// the image for the pipeline container.
	SetupContainer(context.Context, *pipeline.Container) error
	// TailContainer defines a function that captures
	// the logs on the pipeline container.
	TailContainer(context.Context, context.Context, *pipeline.Container) (io.ReadCloser, error)
	// WaitContainer defines a function that blocks
	// until the pipeline container completes.
	WaitContainer(context.Context, *pipeline.Container) error

	// Image Engine Interface Functions

	// CreateImage defines a function that
	// creates the pipeline container image.
	CreateImage(context.Context, *pipeline.Container) error
	// InspectImage defines a function that
	// inspects the pipeline container image.
	InspectImage(context.Context, *pipeline.Container) ([]byte, error)

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
