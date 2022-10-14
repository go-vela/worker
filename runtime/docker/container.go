// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"

	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/image"
)

// InspectContainer inspects the pipeline container.
func (c *client) InspectContainer(ctx context.Context, ctn *pipeline.Container) error {
	c.Logger.Tracef("inspecting container %s", ctn.ID)

	// send API call to inspect the container
	//
	// https://godoc.org/github.com/docker/docker/client#Client.ContainerInspect
	container, err := c.Docker.ContainerInspect(ctx, ctn.ID)
	if err != nil {
		return err
	}

	// capture the container exit code
	//
	// https://godoc.org/github.com/docker/docker/api/types#ContainerState
	ctn.ExitCode = container.State.ExitCode

	return nil
}

// RemoveContainer deletes (kill, remove) the pipeline container.
func (c *client) RemoveContainer(ctx context.Context, ctn *pipeline.Container) error {
	c.Logger.Tracef("removing container %s", ctn.ID)

	// send API call to inspect the container
	//
	// https://godoc.org/github.com/docker/docker/client#Client.ContainerInspect
	container, err := c.Docker.ContainerInspect(ctx, ctn.ID)
	if err != nil {
		return err
	}

	// if the container is paused, restarting or running
	//
	// https://godoc.org/github.com/docker/docker/api/types#ContainerState
	if container.State.Paused ||
		container.State.Restarting ||
		container.State.Running {
		// send API call to kill the container
		//
		// https://godoc.org/github.com/docker/docker/client#Client.ContainerKill
		err := c.Docker.ContainerKill(ctx, ctn.ID, "SIGKILL")
		if err != nil {
			return err
		}
	}

	// create options for removing container
	//
	// https://godoc.org/github.com/docker/docker/api/types#ContainerRemoveOptions
	opts := types.ContainerRemoveOptions{
		Force:         true,
		RemoveLinks:   false,
		RemoveVolumes: true,
	}

	// send API call to remove the container
	//
	// https://godoc.org/github.com/docker/docker/client#Client.ContainerRemove
	err = c.Docker.ContainerRemove(ctx, ctn.ID, opts)
	if err != nil {
		return err
	}

	return nil
}

// RunContainer creates and starts the pipeline container.
func (c *client) RunContainer(ctx context.Context, ctn *pipeline.Container, b *pipeline.Build, r *library.Repo) error {
	c.Logger.Tracef("running container %s", ctn.ID)

	// allocate new container config from pipeline container
	containerConf := ctnConfig(ctn)
	// allocate new host config with volume data
	hostConf := hostConfig(c.Logger, b.ID, ctn.Ulimits, c.config.Volumes)
	// allocate new network config with container name
	networkConf := netConfig(b.ID, ctn.Name)

	// -------------------- Start of TODO: --------------------
	//
	// Remove the below code once the mounting issue with Kaniko is
	// resolved to allow mounting private cert bundles with Vela.
	//
	// This code is required due to a known bug in Kaniko:
	//
	// * https://github.com/go-vela/community/issues/253

	// check if the pipeline container image contains
	// the key words "kaniko" and "vela"
	//
	// this is a soft check for the Vela Kaniko plugin
	if strings.Contains(ctn.Image, "kaniko") &&
		strings.Contains(ctn.Image, "vela") {
		// iterate through the list of host mounts provided
		for i, mount := range hostConf.Mounts {
			// check if the source path or target path
			// for the mount contains "/etc/ssl/certs"
			//
			// this is a soft check for mounting private cert bundles
			if strings.Contains(mount.Source, "/etc/ssl/certs") ||
				strings.Contains(mount.Target, "/etc/ssl/certs") {
				// remove the private cert bundle mount from the host config
				hostConf.Mounts = append(hostConf.Mounts[:i], hostConf.Mounts[i+1:]...)
			}
		}
	}
	//
	// -------------------- End of TODO: --------------------

	// check if the container pull policy is on_start
	if strings.EqualFold(ctn.Pull, constants.PullOnStart) {
		// send API call to create the image
		err := c.CreateImage(ctx, ctn)
		if err != nil {
			return err
		}
	}

	// check if the image is allowed to run privileged
	for _, pattern := range c.config.Images {
		privileged, err := image.IsPrivilegedImage(ctn.Image, pattern)
		if err != nil {
			return err
		}

		if privileged {
			// ensure repo is trusted and therefore allowed to run privileged containers
			if c.config.EnforceTrustedRepos && (r == nil || !r.GetTrusted()) {
				return errors.New("repo must be trusted to run privileged images")
			}

			hostConf.Privileged = true
		}
	}

	// send API call to create the container
	//
	// https://godoc.org/github.com/docker/docker/client#Client.ContainerCreate
	_, err := c.Docker.ContainerCreate(
		ctx,
		containerConf,
		hostConf,
		networkConf,
		nil,
		ctn.ID,
	)
	if err != nil {
		return err
	}

	// create options for starting container
	//
	// https://godoc.org/github.com/docker/docker/api/types#ContainerStartOptions
	opts := types.ContainerStartOptions{}

	// send API call to start the container
	//
	// https://godoc.org/github.com/docker/docker/client#Client.ContainerStart
	err = c.Docker.ContainerStart(ctx, ctn.ID, opts)
	if err != nil {
		return err
	}

	return nil
}

// SetupContainer prepares the image for the pipeline container.
func (c *client) SetupContainer(ctx context.Context, ctn *pipeline.Container, r *library.Repo) error {
	c.Logger.Tracef("setting up for container %s", ctn.ID)

	// handle the container pull policy
	switch ctn.Pull {
	case constants.PullAlways:
		// send API call to create the image
		return c.CreateImage(ctx, ctn)
	case constants.PullNotPresent:
		// handled further down in this function
		break
	case constants.PullNever:
		fallthrough
	case constants.PullOnStart:
		fallthrough
	default:
		c.Logger.Tracef("skipping setup for container %s due to pull policy %s", ctn.ID, ctn.Pull)

		return nil
	}

	// parse image from container
	//
	// https://pkg.go.dev/github.com/go-vela/worker/internal/image#ParseWithError
	_image, err := image.ParseWithError(ctn.Image)
	if err != nil {
		return err
	}

	// check if the container image exists on the host
	//
	// https://godoc.org/github.com/docker/docker/client#Client.ImageInspectWithRaw
	_, _, err = c.Docker.ImageInspectWithRaw(ctx, _image)
	if err == nil {
		return nil
	}

	// if the container image does not exist on the host
	// we attempt to capture it for executing the pipeline
	//
	// https://godoc.org/github.com/docker/docker/client#IsErrNotFound
	if docker.IsErrNotFound(err) {
		// send API call to create the image
		return c.CreateImage(ctx, ctn)
	}

	return err
}

// TailContainer captures the logs for the pipeline container.
func (c *client) TailContainer(ctx context.Context, ctn *pipeline.Container) (io.ReadCloser, error) {
	c.Logger.Tracef("tailing output for container %s", ctn.ID)

	// create options for capturing container logs
	//
	// https://godoc.org/github.com/docker/docker/api/types#ContainerLogsOptions
	opts := types.ContainerLogsOptions{
		Follow:     true,
		ShowStdout: true,
		ShowStderr: true,
		Details:    false,
		Timestamps: false,
	}

	// send API call to capture the container logs
	//
	// https://godoc.org/github.com/docker/docker/client#Client.ContainerLogs
	logs, err := c.Docker.ContainerLogs(ctx, ctn.ID, opts)
	if err != nil {
		return nil, err
	}

	// create in-memory pipe for capturing logs
	rc, wc := io.Pipe()

	// capture all stdout and stderr logs
	go func() {
		c.Logger.Tracef("copying logs for container %s", ctn.ID)

		// copy container stdout and stderr logs to our in-memory pipe
		//
		// https://godoc.org/github.com/docker/docker/pkg/stdcopy#StdCopy
		_, err := stdcopy.StdCopy(wc, wc, logs)
		if err != nil {
			c.Logger.Errorf("unable to copy logs for container: %v", err)
		}

		// close logs buffer
		logs.Close()

		// close in-memory pipe write closer
		wc.Close()
	}()

	return rc, nil
}

// WaitContainer blocks until the pipeline container completes.
func (c *client) WaitContainer(ctx context.Context, ctn *pipeline.Container) error {
	c.Logger.Tracef("waiting for container %s", ctn.ID)

	// send API call to wait for the container completion
	//
	// https://godoc.org/github.com/docker/docker/client#Client.ContainerWait
	wait, errC := c.Docker.ContainerWait(ctx, ctn.ID, container.WaitConditionNotRunning)

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
	// create container config object
	//
	// https://godoc.org/github.com/docker/docker/api/types/container#Config
	config := &container.Config{
		Image:        image.Parse(ctn.Image),
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

	// check if the user is present
	if len(ctn.User) > 0 {
		// add user to container config
		config.User = ctn.User
	}

	return config
}
