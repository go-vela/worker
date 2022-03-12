// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/image"
)

const imagePatch = `
{
  "spec": {
    "containers": [
      {
        "name": "%s",
        "image": "%s"
      }
    ]
  }
}
`

// CreateImage creates the pipeline container image.
func (c *client) CreateImage(ctx context.Context, ctn *pipeline.Container) error {
	c.Logger.Tracef("creating image for container %s", ctn.ID)

	// parse/validate image from container
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/image#ParseWithError
	_, err := image.ParseWithError(ctn.Image)
	if err != nil {
		return err
	}

	// Kubernetes does not have an API to make sure it can access an image,
	// so we have to query the appropriate docker registry ourselves.
	// TODO: query docker registry for the image (if possible)
	//       this might require retrieving the pullSecrets from k8s
	// 		 or have the admin add a Vela accessible secret as well.
	return nil
}

// InspectImage inspects the pipeline container image.
func (c *client) InspectImage(ctx context.Context, ctn *pipeline.Container) ([]byte, error) {
	c.Logger.Tracef("inspecting image for container %s", ctn.ID)

	// TODO: consider updating this command
	//
	// create output for inspecting image
	output := []byte(
		fmt.Sprintf("$ kubectl get pod -o=jsonpath='{.spec.containers[%d].image}' %s\n", ctn.Number, ctn.ID),
	)

	// check if the container pull policy is on start
	if strings.EqualFold(ctn.Pull, constants.PullOnStart) {
		return []byte(
			fmt.Sprintf("skipped for container %s due to pull policy %s\n", ctn.ID, ctn.Pull),
		), nil
	}

	// marshal the image information from the container
	// (-1 to convert to 0-based index, -1 for init which isn't a container)
	image, err := json.MarshalIndent(c.Pod.Spec.Containers[ctn.Number-2].Image, "", " ")
	if err != nil {
		return output, err
	}

	// add new line to end of bytes
	return append(output, append(image, "\n"...)...), nil
}
