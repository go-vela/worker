// SPDX-License-Identifier: Apache-2.0

package linux

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	envparse "github.com/hashicorp/go-envparse"
	"github.com/sirupsen/logrus"
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

	// grab outputs
	outputBytes, err := o.client.Runtime.PollOutputsContainer(ctx, ctn, "/vela/outputs/.env")
	if err != nil {
		return nil, nil, err
	}

	reader := bytes.NewReader(outputBytes)

	outputMap, err := envparse.Parse(reader)
	if err != nil {
		logger.Debugf("unable to parse output map: %v", err)
	}

	// grab masked outputs
	maskedBytes, err := o.client.Runtime.PollOutputsContainer(ctx, ctn, "/vela/outputs/masked.env")
	if err != nil {
		return nil, nil, err
	}

	reader = bytes.NewReader(maskedBytes)

	maskMap, err := envparse.Parse(reader)
	if err != nil {
		logger.Debugf("unable to parse masked output map: %v", err)
	}

	return outputMap, maskMap, nil
}

// pollFiles tails the output for sidecar container.
func (o *outputSvc) pollFiles(ctx context.Context, ctn *pipeline.Container, fileList []string) error {
	// exit if outputs container has not been configured
	if len(ctn.Image) == 0 {
		return nil
	}

	// update engine logger with outputs metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
	logger := o.client.Logger.WithField("test-outputs", ctn.Name)

	logger.Debug("tailing container")
	logger.Debugf("fileList: %v", fileList)
	// grab outputs
	fileNames, err := o.client.Runtime.PollFileNames(ctx, ctn, fileList)
	if err != nil {
		return fmt.Errorf("unable to poll file names: %v", err)
	}
	if len(fileNames) != 0 {
		for _, fileName := range fileNames {
			logger.Infof("fileName: %v", fileName)
			reader, size, err := o.client.Runtime.PollFileContent(ctx, ctn, fileName)
			if err != nil {
				logger.Errorf("unable to poll file content: %v", err)
				return err
			}

			err = o.client.Storage.UploadObject(ctx, &api.Object{
				ObjectName: fileName,
				Bucket:     api.Bucket{BucketName: o.client.Storage.GetBucket(ctx)},
				FilePath:   fileName,
			}, reader, size)
			if err != nil {
				logger.Errorf("unable to upload object: %v", err)
				return err
			}

		}
		logger.Infof("fileNames: %v", fileNames)
		return nil
	}
	logger.Debug("no files found")
	return fmt.Errorf("no files found: %v", err)

	//reader := bytes.NewReader(outputBytes)
	//
	//outputMap, err := envparse.Parse(reader)
	//if err != nil {
	//	logger.Debugf("unable to parse output map: %v", err)
	//}
	//
	//// grab masked outputs
	//maskedBytes, err := o.client.Runtime.PollOutputsContainer(ctx, ctn, "/vela/outputs/masked.env")
	//if err != nil {
	//	return nil, nil, err
	//}
	//
	//reader = bytes.NewReader(maskedBytes)
	//
	//maskMap, err := envparse.Parse(reader)
	//if err != nil {
	//	logger.Debugf("unable to parse masked output map: %v", err)
	//}

	//return outputMap, maskMap, nil
}
