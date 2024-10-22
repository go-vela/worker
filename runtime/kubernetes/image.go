// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
)

const (
	pauseImage = "kubernetes/pause:latest"
	imagePatch = `
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
)

// CreateImage creates the pipeline container image.
func (c *client) CreateImage(ctx context.Context, ctn *pipeline.Container) error {
	c.Logger.Tracef("no-op: creating image for container %s", ctn.ID)

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
	image, err := json.MarshalIndent(
		c.Pod.Spec.Containers[c.containersLookup[ctn.ID]].Image, "", " ",
	)
	if err != nil {
		return output, err
	}

	// add new line to end of bytes
	return append(output, append(image, "\n"...)...), nil
}
