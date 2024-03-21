// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	dockerImageTypes "github.com/docker/docker/api/types/image"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/image"
	"github.com/sirupsen/logrus"
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
	// https://pkg.go.dev/github.com/docker/docker/api/types/image#PullOptions
	opts := dockerImageTypes.PullOptions{}

	// send API call to pull the image for the container
	//
	// https://pkg.go.dev/github.com/docker/docker/client#Client.ImagePull
	reader, err := c.Docker.ImagePull(ctx, _image, opts)
	if err != nil {
		return err
	}
	defer reader.Close()

	// check if logrus is set up with trace level
	if logrus.GetLevel() == logrus.TraceLevel {
		// copy output from image pull to standard output
		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			return err
		}
	} else {
		// discard output from image pull
		_, err = io.Copy(io.Discard, reader)
		if err != nil {
			return err
		}
	}

	return nil
}

// InspectImage inspects the pipeline container image.
func (c *client) InspectImage(ctx context.Context, ctn *pipeline.Container) ([]byte, error) {
	c.Logger.Tracef("inspecting image for container %s", ctn.ID)

	// create output for inspecting image
	output := []byte(
		fmt.Sprintf("$ docker image inspect %s\n", ctn.Image),
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
	// https://pkg.go.dev/github.com/docker/docker/client#Client.ImageInspectWithRaw
	i, _, err := c.Docker.ImageInspectWithRaw(ctx, _image)
	if err != nil {
		return output, err
	}

	// add new line to end of bytes
	return append(output, []byte(i.ID+"\n")...), nil
}
