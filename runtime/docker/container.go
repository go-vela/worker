// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/go-vela/types/pipeline"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"

	"github.com/sirupsen/logrus"
)

// InspectContainer inspects the pipeline container.
func (c *client) InspectContainer(ctx context.Context, ctn *pipeline.Container) ([]byte, error) {
	logrus.Tracef("Inspecting container for step %s", ctn.ID)

	// send API call to inspect the container
	container, err := c.Runtime.ContainerInspect(ctx, ctn.ID)
	if err != nil {
		return nil, err
	}

	// set the exit code
	ctn.ExitCode = container.State.ExitCode

	return []byte(container.Image + "\n"), nil
}

// RemoveContainer deletes (kill, remove) the pipeline container.
func (c *client) RemoveContainer(ctx context.Context, ctn *pipeline.Container) error {
	logrus.Tracef("Removing container for step %s", ctn.ID)

	// send API call to inspect the container
	container, err := c.Runtime.ContainerInspect(ctx, ctn.ID)
	if err != nil {
		return err
	}

	// if the container is paused, restarting or running
	if container.State.Paused ||
		container.State.Restarting ||
		container.State.Running {
		// send API call to kill the container
		err := c.Runtime.ContainerKill(ctx, ctn.ID, "SIGKILL")
		if err != nil {
			return err
		}
	}

	// create options for removing container
	opts := types.ContainerRemoveOptions{
		Force:         true,
		RemoveLinks:   false,
		RemoveVolumes: true,
	}

	// send API call to remove the container
	err = c.Runtime.ContainerRemove(ctx, ctn.ID, opts)
	if err != nil {
		return err
	}

	return nil
}

// RunContainer creates and start the pipeline container.
func (c *client) RunContainer(ctx context.Context, b *pipeline.Build, ctn *pipeline.Container) error {
	// create container configuration
	ctnConf := ctnConfig(ctn)
	// create host configuration
	hostConf := hostConfig(b.ID)
	// create network configuration
	netConf := netConfig(b.ID, ctn.Name)

	logrus.Tracef("Creating container for step %s", b.ID)

	// send API call to create the container
	container, err := c.Runtime.ContainerCreate(
		ctx,
		ctnConf,
		hostConf,
		netConf,
		ctn.ID,
	)
	if err != nil {
		return err
	}

	logrus.Tracef("Starting container for step %s", b.ID)

	// create options for starting container
	opts := types.ContainerStartOptions{}

	// send API call to start the container
	err = c.Runtime.ContainerStart(ctx, container.ID, opts)
	if err != nil {
		return err
	}

	return nil
}

// SetupContainer pulls the image for the pipeline container.
func (c *client) SetupContainer(ctx context.Context, ctn *pipeline.Container) error {
	logrus.Tracef("Parsing image %s", ctn.Image)

	// parse image from container
	image, err := parseImage(ctn.Image)
	if err != nil {
		return err
	}

	// check if the container should be updated
	if ctn.Pull {
		logrus.Tracef("Pulling configured image %s", image)
		// create options for pulling image
		opts := types.ImagePullOptions{}

		// send API call to pull the image for the container
		reader, err := c.Runtime.ImagePull(ctx, image, opts)
		if err != nil {
			return err
		}
		defer reader.Close()

		// copy output from image pull to standard output
		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			return err
		}

		return nil
	}

	// check if the container image exists on the host
	_, _, err = c.Runtime.ImageInspectWithRaw(ctx, image)
	if err == nil {
		return nil
	}

	// if the container image does not exist on the host
	// we attempt to capture it for executing the pipeline
	if docker.IsErrNotFound(err) {
		logrus.Tracef("Pulling unfound image %s", image)

		// create options for pulling image
		opts := types.ImagePullOptions{}

		// send API call to pull the image for the container
		reader, err := c.Runtime.ImagePull(ctx, image, opts)
		if err != nil {
			return err
		}
		defer reader.Close()

		// copy output from image pull to standard output
		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			return err
		}

		return nil
	}

	return err
}

// TailContainer captures the logs for the pipeline container.
func (c *client) TailContainer(ctx context.Context, ctn *pipeline.Container) (io.ReadCloser, error) {
	logrus.Tracef("Capturing container logs for step %s", ctn.ID)

	// create options for capturing container logs
	opts := types.ContainerLogsOptions{
		Follow:     true,
		ShowStdout: true,
		ShowStderr: true,
		Details:    false,
		Timestamps: false,
	}

	// send API call to capture the container logs
	logs, err := c.Runtime.ContainerLogs(ctx, ctn.ID, opts)
	if err != nil {
		return nil, err
	}

	// create in-memory pipe for capturing logs
	rc, wc := io.Pipe()

	logrus.Tracef("Copying container logs for step %s", ctn.ID)

	// capture all stdout and stderr logs
	go func() {
		stdcopy.StdCopy(wc, wc, logs)
		logs.Close()
		wc.Close()
		rc.Close()
	}()

	return rc, nil
}

// WaitContainer blocks until the pipeline container completes.
func (c *client) WaitContainer(ctx context.Context, ctn *pipeline.Container) error {
	logrus.Tracef("Waiting for container for step %s", ctn.ID)

	// send API call to wait for the container completion
	wait, errC := c.Runtime.ContainerWait(ctx, ctn.ID, container.WaitConditionNotRunning)
	select {
	case <-wait:
	case err := <-errC:
		return err
	}

	return nil
}

// ctnConfig is a helper function to
// generate the container config.
func ctnConfig(ctn *pipeline.Container) *container.Config {
	logrus.Tracef("Creating container configuration for step %s", ctn.ID)

	// parse image from container
	image, err := parseImage(ctn.Image)
	if err != nil {
		logrus.Errorf("unable to parse image: %w", err)
	}

	// create container config object
	config := &container.Config{
		Image:        image,
		WorkingDir:   ctn.Directory,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		OpenStdin:    false,
		StdinOnce:    false,
		ArgsEscaped:  false,
	}

	// check if the environment is provided
	if len(ctn.Environment) > 0 {
		// iterate through each element in the container environment
		for k, v := range ctn.Environment {
			// add key/value environment to container config
			config.Env = append(config.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// check if the entrypoint is provided
	if len(ctn.Entrypoint) > 0 {
		// add entrypoint to container config
		config.Entrypoint = ctn.Entrypoint
	}

	// check if the commands are provided
	if len(ctn.Commands) > 0 {
		// add commands to container config
		config.Cmd = ctn.Commands
	}

	return config
}

// hostConfig is a helper function to generate
// the host config for a container.
func hostConfig(id string) *container.HostConfig {
	return &container.HostConfig{
		LogConfig: container.LogConfig{
			Type: "json-file",
		},
		Privileged: false,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: id,
				Target: "/home",
			},
		},
	}
}

// parseImage is a helper function to parse
// the image for the provided container.
func parseImage(s string) (string, error) {
	// create fully qualified reference
	image, err := reference.ParseNormalizedNamed(s)
	if err != nil {
		return "", err
	}

	// add latest tag to image if no tag was provided
	return reference.TagNameOnly(image).String(), nil
}
