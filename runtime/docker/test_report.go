package docker

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/go-vela/server/compiler/types/pipeline"
	"io"
	"path/filepath"
	"strings"
)

// PollFileNames grabs files name from provided path
// within a container and uploads them to s3
func (c *client) PollFileNames(ctx context.Context, ctn *pipeline.Container, paths []string) ([]string, error) {
	c.Logger.Tracef("gathering test results and attachments from container %s", ctn.ID)

	var fullFilePaths []string
	if len(ctn.Image) == 0 {
		return nil, nil
	}
	// iterate through the steps in the build
	// iterate through the results paths and store them in the map
	for _, path := range paths {
		dir, filename := filepath.Split(path)
		c.Logger.Tracef("searching for file %s in %s", filename, dir)

		execConfig := types.ExecConfig{
			Tty: true,
			//Cmd:          []string{"sh", "-c", fmt.Sprintf("find %s -type f -name %s", dir, filename)},
			Cmd:          []string{"sh", "-c", fmt.Sprintf("find / -type f -path *%s  -print", path)},
			AttachStderr: true,
			AttachStdout: true,
		}

		c.Logger.Infof("executing command: %v", execConfig.Cmd)
		responseExec, err := c.Docker.ContainerExecCreate(ctx, ctn.ID, execConfig)
		if err != nil {
			c.Logger.Errorf("unable to create exec for container: %v", err)
			return nil, err
		}

		hijackedResponse, err := c.Docker.ContainerExecAttach(ctx, responseExec.ID, types.ExecStartCheck{})
		if err != nil {
			c.Logger.Errorf("unable to attach to exec for container: %v", err)
			return nil, err
		}

		defer func() {
			if hijackedResponse.Conn != nil {
				hijackedResponse.Close()
			}
		}()

		outputStdout := new(bytes.Buffer)
		outputStderr := new(bytes.Buffer)

		if hijackedResponse.Reader != nil {
			_, err := stdcopy.StdCopy(outputStdout, outputStderr, hijackedResponse.Reader)
			if err != nil {
				c.Logger.Errorf("unable to copy logs for container: %v", err)
			}
		}

		if outputStderr.Len() > 0 {
			return nil, fmt.Errorf("error: %s", outputStderr.String())
		}

		data := outputStdout.String()
		c.Logger.Infof("found files: %s", data)

		filePaths := strings.Split(data, "\n")
		for _, filePath := range filePaths {
			if filePath != "" {

				fullFilePaths = append(fullFilePaths, strings.TrimSpace(filePath))
				c.Logger.Infof("full file: %s", filePath)
			}
		}
	}
	if len(fullFilePaths) == 0 {
		return nil, fmt.Errorf("no matching files found for any provided paths")
	}

	return fullFilePaths, nil

	// iterate through the steps in the build
	//for _, step := range p.Steps {
	//	if len(step.TestReport.Results) == 0 {
	//		c.Logger.Warnf("no results provided for the step %s", step.ID)
	//		return fmt.Errorf("no results provided for the step %s", step.ID)
	//	}
	//	if len(step.TestReport.Attachments) == 0 {
	//		c.Logger.Warnf("no attachments provided for the step %s", step.ID)
	//		return fmt.Errorf("no attachments provided for the step %s", step.ID)
	//	}
	//	// check if the step has the provided paths from results
	//	for _, result := range step.TestReport.Results {
	//		_, err := os.Stat(result)
	//		if err != nil {
	//			c.Logger.Errorf("unable to find test result %s for step %s", result, step.ID)
	//			continue
	//		}
	//	}
	//}
	//return nil
}

// PollFileContent retrieves the content and size of a file inside a container.
func (c *client) PollFileContent(ctx context.Context, ctn *pipeline.Container, path string) (io.Reader, int64, error) {
	c.Logger.Tracef("gathering test results and attachments from container %s", ctn.ID)

	if len(ctn.Image) == 0 {
		return nil, 0, nil
	}
	cmd := []string{"sh", "-c", fmt.Sprintf("cat %s", path)}
	execConfig := types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: false,
		Tty:          true,
	}

	c.Logger.Infof("executing command for content: %v", execConfig.Cmd)
	execID, err := c.Docker.ContainerExecCreate(ctx, ctn.ID, execConfig)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create exec instance: %w", err)
	}
	resp, err := c.Docker.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to attach to exec instance: %w", err)
	}

	defer func() {
		if resp.Conn != nil {
			resp.Close()
		}
	}()

	outputStdout := new(bytes.Buffer)
	outputStderr := new(bytes.Buffer)

	if resp.Reader != nil {
		_, err := stdcopy.StdCopy(outputStdout, outputStderr, resp.Reader)
		if err != nil {
			c.Logger.Errorf("unable to copy logs for container: %v", err)
		}
	}

	if outputStderr.Len() > 0 {
		return nil, -1, fmt.Errorf("error: %s", outputStderr.String())
	}

	data := outputStdout.Bytes()
	//c.Logger.Infof("data: %v", string(data))

	// convert the data to a reader
	reader := bytes.NewReader(data)
	// get the size of the data
	size := int64(len(data))
	return reader, size, nil
}
