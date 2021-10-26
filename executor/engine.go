// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executor

import (
	"context"
	"sync"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

// Engine represents the interface for Vela integrating
// with the different supported operating systems.
type Engine interface {

	// Engine Interface Functions

	// Driver defines a function that outputs
	// the configured executor driver.
	Driver() string

	// API interface functions

	// GetBuild defines a function for the API
	// that gets the current build in execution.
	GetBuild() (*library.Build, error)
	// GetPipeline defines a function for the API
	// that gets the current pipeline in execution.
	GetPipeline() (*pipeline.Build, error)
	// GetRepo defines a function for the API
	// that gets the current repo in execution.
	GetRepo() (*library.Repo, error)
	// CancelBuild defines a function for the API
	// that Cancels the current build in execution.
	CancelBuild() (*library.Build, error)

	// Build Engine interface functions

	// CreateBuild defines a function that
	// configures the build for execution.
	CreateBuild(context.Context) error
	// PlanBuild defines a function that
	// handles the resource initialization process
	// for the build.
	PlanBuild(context.Context) error
	// AssembleBuild defines a function that
	// prepares the containers within a build
	// for execution.
	AssembleBuild(context.Context) error
	// ExecBuild defines a function that
	// runs a pipeline for a build.
	ExecBuild(context.Context) error
	// DestroyBuild defines a function that
	// cleans up the build after execution.
	DestroyBuild(context.Context) error

	// Service Engine Interface Functions

	// CreateService defines a function that
	// configures the service for execution.
	CreateService(context.Context, *pipeline.Container) error
	// PlanService defines a function that
	// prepares the service for execution.
	PlanService(context.Context, *pipeline.Container) error
	// ExecService defines a function that
	// runs a service.
	ExecService(context.Context, *pipeline.Container) error
	// StreamService defines a function that
	// tails the output for a service.
	StreamService(context.Context, *pipeline.Container) error
	// DestroyService defines a function that
	// cleans up the service after execution.
	DestroyService(context.Context, *pipeline.Container) error

	// Stage Engine Interface Functions

	// CreateStage defines a function that
	// configures the stage for execution.
	CreateStage(context.Context, *pipeline.Stage) error
	// PlanStage defines a function that
	// prepares the stage for execution.
	PlanStage(context.Context, *pipeline.Stage, *sync.Map) error
	// ExecStage defines a function that
	// runs a stage.
	ExecStage(context.Context, *pipeline.Stage, *sync.Map) error
	// DestroyStage defines a function that
	// cleans up the stage after execution.
	DestroyStage(context.Context, *pipeline.Stage) error

	// Step Engine Interface Functions

	// CreateStep defines a function that
	// configures the step for execution.
	CreateStep(context.Context, *pipeline.Container) error
	// PlanStep defines a function that
	// prepares the step for execution.
	PlanStep(context.Context, *pipeline.Container) error
	// ExecStep defines a function that
	// runs a step.
	ExecStep(context.Context, *pipeline.Container) error
	// StreamStep defines a function that
	// tails the output for a step.
	StreamStep(context.Context, *pipeline.Container) error
	// DestroyStep defines a function that
	// cleans up the step after execution.
	DestroyStep(context.Context, *pipeline.Container) error
}
