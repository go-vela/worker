// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"maps"

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

	// set up outputs logger
	logger := o.client.Logger.WithField("outputs", "outputs")

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

	logger.Debug("setting up outputs container")
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

	// update engine logger with outputs metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := o.client.Logger.WithField("outputs", ctn.Name)

	logger.Debug("inspecting outputs container")
	// inspect the runtime container
	err := o.client.Runtime.InspectContainer(ctx, ctn)
	if err != nil {
		return err
	}

	logger.Debug("removing outputs container")
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

	logrus.Debug("running outputs container")
	// run the runtime container
	err := o.client.Runtime.RunContainer(ctx, _outputs, o.client.pipeline)
	if err != nil {
		return err
	}

	logrus.Debug("inspecting outputs container")
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

	// update engine logger with outputs metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := o.client.Logger.WithField("outputs", ctn.Name)

	logger.Debug("tailing container")

	var (
		filePaths = []string{
			"/vela/outputs/.env",
			"/vela/outputs/masked.env",
			"/vela/outputs/base64.env",
			"/vela/outputs/masked.base64.env",
		}

		outputMap = make(map[string]string)
		maskMap   = make(map[string]string)
	)

	for _, p := range filePaths {
		outputBytes, err := o.client.Runtime.PollOutputsContainer(ctx, ctn, p)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to poll outputs container %s: %w", ctn.Name, err)
		}

		reader := bytes.NewReader(outputBytes)

		switch p {
		case "/vela/outputs/.env":
			parsed, err := envparse.Parse(reader)
			if err != nil {
				logger.Debugf("unable to parse output map: %v", err)
			}

			// add to output map
			maps.Copy(outputMap, parsed)

		case "/vela/outputs/masked.env":
			parsed, err := envparse.Parse(reader)
			if err != nil {
				logger.Debugf("unable to parse masked output map: %v", err)
			}

			// add to mask map
			maps.Copy(maskMap, parsed)

		case "/vela/outputs/base64.env":
			parsed, err := envparse.Parse(reader)
			if err != nil {
				logger.Debugf("unable to parse base64 output map: %v", err)
			}

			for k, v := range parsed {
				// decode the base64 value
				decodedValue, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					logger.Debugf("unable to decode base64 value for key %s: %v", k, err)
					continue
				}

				parsed[k] = string(decodedValue)
			}

			// add to output map
			maps.Copy(outputMap, parsed)

		case "/vela/outputs/masked.base64.env":
			parsed, err := envparse.Parse(reader)
			if err != nil {
				logger.Debugf("unable to parse masked base64 output map: %v", err)
			}

			for k, v := range parsed {
				// decode the base64 value
				decodedValue, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					logger.Debugf("unable to decode base64 value for key %s: %v", k, err)
					continue
				}

				parsed[k] = string(decodedValue)
			}

			// add to mask map
			maps.Copy(maskMap, parsed)
		}
	}

	return outputMap, maskMap, nil
}
