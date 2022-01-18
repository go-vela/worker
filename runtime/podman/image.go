// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package podman

import (
	"context"
	"fmt"
	"strings"

	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/image"
)

// CreateImage creates the pipeline container image.
func (c *client) CreateImage(ctx context.Context, ctn *pipeline.Container) error {
	c.Logger.Tracef("creating image for container %s", ctn.ID)

	// parse image from container
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/image#ParseWithError
	_image, err := image.ParseWithError(ctn.Image)
	if err != nil {
		return err
	}

	// create options for pulling image
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/images#PullOptions
	opts := &images.PullOptions{}

	// send API call to pull the image for the container
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/images#Pull
	_, err = images.Pull(c.Podman, _image, opts)
	if err != nil {
		return err
	}

	return nil
}

// InspectImage inspects the pipeline container image.
func (c *client) InspectImage(ctx context.Context, ctn *pipeline.Container) ([]byte, error) {
	c.Logger.Tracef("inspecting image for container %s", ctn.ID)

	// create output for inspecting image
	output := []byte(
		fmt.Sprintf("$ podman image inspect %s\n", ctn.Image),
	)

	// parse image from container
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/image#ParseWithError
	_image, err := image.ParseWithError(ctn.Image)
	if err != nil {
		return output, err
	}

	// check if the container pull policy is on start
	if strings.EqualFold(ctn.Pull, constants.PullOnStart) {
		return []byte(
			fmt.Sprintf("skipped for container %s due to pull policy %s\n", ctn.ID, ctn.Pull),
		), nil
	}

	// send API call to inspect the image
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/images#GetImage
	i, err := images.GetImage(c.Podman, _image, &images.GetOptions{})
	if err != nil {
		return output, err
	}

	// add new line to end of bytes
	return append(output, []byte(i.ID+"\n")...), nil
}
