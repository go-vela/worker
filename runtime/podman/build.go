// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package podman

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/containers/podman/v3/pkg/bindings/pods"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman/v3/pkg/specgen"
	"github.com/go-vela/types/pipeline"
)

// InspectBuild displays details about the pod for the init step.
func (c *client) InspectBuild(ctx context.Context, b *pipeline.Build) ([]byte, error) {
	c.Logger.Tracef("inspecting build pod for pipeline %s", b.ID)

	// create output for inspecting Pod
	output := []byte(fmt.Sprintf("> Inspecting pod for pipeline %s\n", b.ID))

	// send API call to capture the Pod
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/pods#Inspect
	r, err := pods.Inspect(c.Podman, b.ID, &pods.InspectOptions{})
	if err != nil {
		return []byte{}, err
	}

	// convert the Pod type PodInspectReport to bytes with pretty print
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/domain/entities#PodInspectReport
	podOutput, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		return []byte{}, fmt.Errorf("unable to serialize pod data: %w", err)
	}

	// TODO: add newline at end of output
	output = append(output, podOutput...)

	return output, nil
}

// SetupBuild prepares the pipeline build.
func (c *client) SetupBuild(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("setting up pod for build %s", b.ID)

	// create a new Pod Spec
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/domain/entities#PodSpec
	podSpec := new(entities.PodSpec)
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/specgen#PodSpecGenerator
	podGen := specgen.NewPodSpecGenerator()
	podGen.Name = b.ID

	podSpec.PodSpecGen = *podGen

	// store podSpec on client
	c.Pod = podSpec

	// create the pod with the created Spec
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/pods#CreatePodFromSpec
	_, err := pods.CreatePodFromSpec(c.Podman, c.Pod)
	if err != nil {
		return err
	}

	c.Logger.Info("created pod for build %s", b.ID)

	return nil
}

// AssembleBuild finalizes pipeline build setup.
// This is a no-op for podman.
func (c *client) AssembleBuild(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("no-op: assembling build %s", b.ID)

	return nil
}

// RemoveBuild deletes (kill, remove) the pipeline build metadata.
func (c *client) RemoveBuild(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("removing build %s", b.ID)

	// check if the pod exists
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/pods#Exists
	podExists, err := pods.Exists(c.Podman, b.ID, &pods.ExistsOptions{})
	if err != nil {
		return err
	}

	// exit if it doesn't
	if !podExists {
		return nil
	}

	// remove options for removing the Pod
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/pods#RemoveOptions
	rmOpts := new(pods.RemoveOptions).WithForce(true)

	// remove the Pod
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/pods#Remove
	_, err = pods.Remove(c.Podman, b.ID, rmOpts)
	if err != nil {
		return err
	}

	c.Logger.Infof("removed build %s", b.ID)

	return nil
}
