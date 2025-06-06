// SPDX-License-Identifier: Apache-2.0

package local

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	envparse "github.com/hashicorp/go-envparse"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/compiler/types/pipeline"
)

// outputSvc handles communication with the outputs container during the build.
type outputSvc svc

// create configures the outputs container for execution.
func (o *outputSvc) create(ctx context.Context, ctn *pipeline.Container, timeout int64) error {
	// exit if outputs container has not been configured
	if len(ctn.Image) == 0 {
		return nil
	}

	// Encode script content to Base64
	script := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("mkdir /vela/outputs\nsleep %d\n", timeout)),
	)

	// set the entrypoint for the ctn
	ctn.Entrypoint = []string{"/bin/sh", "-c"}

	// set the commands for the ctn
	ctn.Commands = []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"}

	// set the environment variables for the ctn
	ctn.Environment["HOME"] = "/root"
	ctn.Environment["SHELL"] = "/bin/sh"
	ctn.Environment["VELA_BUILD_SCRIPT"] = script

	// setup the runtime container
	err := o.client.Runtime.SetupContainer(ctx, ctn)
	if err != nil {
		return err
	}

	return nil
}

// destroy cleans up outputs container after execution.
func (o *outputSvc) destroy(ctx context.Context, ctn *pipeline.Container) error {
	// exit if outputs container has not been configured
	if len(ctn.Image) == 0 {
		return nil
	}

	// inspect the runtime container
	err := o.client.Runtime.InspectContainer(ctx, ctn)
	if err != nil {
		return err
	}

	// remove the runtime container
	err = o.client.Runtime.RemoveContainer(ctx, ctn)
	if err != nil {
		return err
	}

	return nil
}

// exec runs the outputs sidecar container for a pipeline.
func (o *outputSvc) exec(ctx context.Context, _outputs *pipeline.Container) error {
	// exit if outputs container has not been configured
	if len(_outputs.Image) == 0 {
		return nil
	}

	// run the runtime container
	err := o.client.Runtime.RunContainer(ctx, _outputs, o.client.pipeline)
	if err != nil {
		return err
	}

	// inspect the runtime container
	err = o.client.Runtime.InspectContainer(ctx, _outputs)
	if err != nil {
		return err
	}

	return nil
}

// poll tails the output for sidecar container.
func (o *outputSvc) poll(ctx context.Context, ctn *pipeline.Container) (map[string]string, map[string]string, error) {
	// exit if outputs container has not been configured
	if len(ctn.Image) == 0 {
		return nil, nil, nil
	}

	// grab outputs
	outputBytes, err := o.client.Runtime.PollOutputsContainer(ctx, ctn, "/vela/outputs/.env")
	if err != nil {
		return nil, nil, err
	}

	reader := bytes.NewReader(outputBytes)

	outputMap, err := envparse.Parse(reader)
	if err != nil {
		logrus.Debugf("unable to parse output map: %v", err)
	}

	// grab masked outputs
	maskedBytes, err := o.client.Runtime.PollOutputsContainer(ctx, ctn, "/vela/outputs/masked.env")
	if err != nil {
		return nil, nil, err
	}

	reader = bytes.NewReader(maskedBytes)

	maskMap, err := envparse.Parse(reader)
	if err != nil {
		logrus.Debugf("unable to parse masked output map: %v", err)
	}

	return outputMap, maskMap, nil
}
