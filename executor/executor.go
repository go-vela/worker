// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executor

import (
	"context"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

// Engine represents the interface for Vela integrating
// with the different supported operating systems.
type Engine interface {

	// API interface functions

	// GetBuild defines a function for the API
	// that gets the current build in execution.
	GetBuild() (*library.Build, error)
	// GetRepo defines a function for the API
	// that gets the current repo in execution.
	GetRepo() (*library.Repo, error)
	// GetPipeline defines a function for the API
	// that gets the current pipeline in execution.
	GetPipeline() (*pipeline.Build, error)
	// KillBuild defines a function for the API
	// that kills the current build in execution.
	KillBuild() (*library.Build, error)

	// Secrets Engine Interface Functions

	// PullSecret defines a function that pulls
	// the secrets for a given pipeline.
	PullSecret(context.Context) error

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

	// Stage Engine Interface Functions

	// CreateStage defines a function that
	// configures the stage for execution.
	CreateStage(context.Context, *pipeline.Stage) error
	// PlanStage defines a function that
	// prepares the stage for execution.
	PlanStage(context.Context, *pipeline.Stage, map[string]chan error) error
	// ExecStage defines a function that
	// runs a stage.
	ExecStage(context.Context, *pipeline.Stage, map[string]chan error) error
	// DestroyStage defines a function that
	// cleans up the stage after execution.
	DestroyStage(context.Context, *pipeline.Stage) error

	// Build Engine interface functions

	// CreateBuild defines a function that
	// configures the build for execution.
	CreateBuild(context.Context) error
	// PlanBuild defines a function that
	// prepares the build for execution.
	PlanBuild(context.Context) error
	// ExecBuild defines a function that
	// runs a pipeline for a build.
	ExecBuild(context.Context) error
	// DestroyBuild defines a function that
	// cleans up the build after execution.
	DestroyBuild(context.Context) error

	// With Engine interface functions

	// WithBuild defines a function that sets
	// the library Build type in the Engine.
	WithBuild(*library.Build) Engine
	// WithPipeline defines a function that sets
	// the pipeline Build type in the Engine.
	WithPipeline(*pipeline.Build) Engine
	// WithRepo defines a function that sets
	// the library Repo type in the Engine.
	WithRepo(*library.Repo) Engine
	// WithUser defines a function that sets
	// the library User type in the Engine.
	WithUser(*library.User) Engine
}
