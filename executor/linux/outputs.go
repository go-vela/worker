// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
)

// outputSvc handles communication with the outputs container during the build.
type outputSvc svc

// traceScript is a helper script that is added to the build script
// to trace a command.
const traceScript = `
echo $ %s
%s
`

// create configures the outputs plugin for execution.
func (o *outputSvc) create(ctx context.Context, ctn *pipeline.Container, timeout int64) error {
	// exit if outputs container has not been configured
	if len(ctn.Image) == 0 {
		return nil
	}

	// update engine logger with secret metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := o.client.Logger.WithField("outputs", "outputs")

	// generate script from commands
	script := generateScriptPosix([]string{fmt.Sprintf("sleep %d", timeout)})

	// set the entrypoint for the ctn
	ctn.Entrypoint = []string{"/bin/sh", "-c"}

	// set the commands for the ctn
	ctn.Commands = []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"}

	// set the environment variables for the ctn
	ctn.Environment["HOME"] = "/root"
	ctn.Environment["SHELL"] = "/bin/sh"
	ctn.Environment["VELA_BUILD_SCRIPT"] = script
	ctn.Environment["VELA_DISTRIBUTION"] = o.client.build.GetDistribution()
	ctn.Environment["BUILD_HOST"] = o.client.build.GetHost()
	ctn.Environment["VELA_HOST"] = o.client.build.GetHost()
	ctn.Environment["VELA_RUNTIME"] = o.client.build.GetRuntime()
	ctn.Environment["VELA_VERSION"] = o.client.Version

	logger.Debug("setting up container")
	// setup the runtime container
	err := o.client.Runtime.SetupContainer(ctx, ctn)
	if err != nil {
		return err
	}

	return nil
}

// destroy cleans up secret plugin after execution.
func (o *outputSvc) destroy(ctx context.Context, ctn *pipeline.Container) error {
	// exit if outputs container has not been configured
	if len(ctn.Image) == 0 {
		return nil
	}

	// update engine logger with secret metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := o.client.Logger.WithField("secret", ctn.Name)

	logger.Debug("inspecting container")
	// inspect the runtime container
	err := o.client.Runtime.InspectContainer(ctx, ctn)
	if err != nil {
		return err
	}

	logger.Debug("removing container")
	// remove the runtime container
	err = o.client.Runtime.RemoveContainer(ctx, ctn)
	if err != nil {
		return err
	}

	return nil
}

// generateScriptPosix is a helper function that generates a build script
// for a linux container using the given commands.
func generateScriptPosix(commands []string) string {
	var buf bytes.Buffer

	// iterate through each command provided
	for _, command := range commands {
		// safely escape entire command
		escaped := fmt.Sprintf("%q", command)

		// safely escape trace character
		escaped = strings.Replace(escaped, "$", `\$`, -1)

		// write escaped lines to buffer
		buf.WriteString(fmt.Sprintf(
			traceScript,
			escaped,
			command,
		))
	}

	// create build script with netrc and buffer information
	script := buf.String()

	return base64.StdEncoding.EncodeToString([]byte(script))
}

// exec runs a secret plugins for a pipeline.
func (o *outputSvc) exec(ctx context.Context, _outputs *pipeline.Container) error {
	// exit if outputs container has not been configured
	if len(_outputs.Image) == 0 {
		return nil
	}

	logrus.Debug("running container")
	// run the runtime container
	err := o.client.Runtime.RunContainer(ctx, _outputs, o.client.pipeline)
	if err != nil {
		return err
	}

	logrus.Debug("inspecting container")
	// inspect the runtime container
	err = o.client.Runtime.InspectContainer(ctx, _outputs)
	if err != nil {
		return err
	}

	return nil
}

// poll tails the output for a secret plugin.
func (o *outputSvc) poll(ctx context.Context, ctn *pipeline.Container) (map[string]string, map[string]string, error) {
	// exit if outputs container has not been configured
	if len(ctn.Image) == 0 {
		return nil, nil, nil
	}

	// update engine logger with secret metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := o.client.Logger.WithField("secret", ctn.Name)

	logger.Debug("tailing container")

	// grab outputs
	outputBytes, err := o.client.Runtime.PollOutputsContainer(ctx, ctn, "/vela/outputs.env")
	if err != nil {
		return nil, nil, err
	}

	// grab masked outputs
	maskedBytes, err := o.client.Runtime.PollOutputsContainer(ctx, ctn, "/vela/masked_outputs.env")
	if err != nil {
		return nil, nil, err
	}

	return toMap(outputBytes), toMap(maskedBytes), nil
}

// toMap is a helper function that turns raw docker exec output bytes into a map
// by splitting on carriage returns + newlines and once more on `=`.
func toMap(input []byte) map[string]string {
	str := string(input[8:]) // Ignore first 8 bytes of Docker output.

	logrus.Infof("string to split: %s", str)
	lines := strings.Split(str, "\r\n")

	m := make(map[string]string)

	for _, line := range lines {
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			s := parts[1]
			if !strings.Contains(parts[1], "\\\\n") {
				s = strings.Replace(parts[1], "\\n", "\\\n", -1)
			}

			m[parts[0]] = s
		}
	}

	return m
}
