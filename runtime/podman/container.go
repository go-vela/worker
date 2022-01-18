// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package podman

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/containers/podman/v3/libpod/define"
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/containers/podman/v3/pkg/specgen"
	"github.com/go-vela/types/constants"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/image"
)

// InspectContainer inspects the pipeline container.
func (c *client) InspectContainer(ctx context.Context, ctn *pipeline.Container) error {
	c.Logger.Tracef("inspecting container %s", ctn.ID)

	// send API call to inspect the container
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/containers#Inspect
	container, err := containers.Inspect(c.Podman, normalizeContainerName(ctn.Name), &containers.InspectOptions{})
	if err != nil {
		return err
	}

	// capture the container exit code
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/libpod/define#InspectContainerState
	ctn.ExitCode = int(container.State.ExitCode)

	return nil
}

// RemoveContainer deletes (kill, remove) the pipeline container.
// This is a no-op for podman. RemoveBuild handles deleting the Pod
// and all resources within.
func (c *client) RemoveContainer(ctx context.Context, ctn *pipeline.Container) error {
	c.Logger.Tracef("no-op: removing container %s", ctn.ID)

	return nil
}

// RunContainer creates and starts the pipeline container.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) RunContainer(ctx context.Context, ctn *pipeline.Container, b *pipeline.Build) error {
	c.Logger.Tracef("running container %s", ctn.ID)

	// allocate new container config from pipeline container and configuration
	spec := specConf(ctn, b.ID, c.Logger)

	// allocate storage configuration for container
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/specgen#ContainerStorageConfig
	spec.ContainerStorageConfig = storageConfig(ctn, b.ID, c.config.Volumes, c.Logger)

	// validate the generated spec
	if err := spec.Validate(); err != nil {
		return err
	}

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
		for i, mount := range spec.Mounts {
			// check if the source path or target path
			// for the mount contains "/etc/ssl/certs"
			//
			// this is a soft check for mounting private cert bundles
			if strings.Contains(mount.Source, "/etc/ssl/certs") ||
				strings.Contains(mount.Destination, "/etc/ssl/certs") {
				// remove the private cert bundle mount from the host config
				spec.Mounts = append(spec.Mounts[:i], spec.Mounts[i+1:]...)
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
	// TODO: privileged should not be possible if Podman
	// is running in user space, so we need to add a check here
	for _, pattern := range c.config.Images {
		privileged, err := image.IsPrivilegedImage(ctn.Image, pattern)
		if err != nil {
			return err
		}

		spec.Privileged = privileged
	}

	// send API call to create the container
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/containers#CreateWithSpec
	_, err := containers.CreateWithSpec(c.Podman, spec, nil)
	if err != nil {
		return err
	}

	// send API call to start the container
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/containers#Start
	err = containers.Start(c.Podman, normalizeContainerName(ctn.Name), &containers.StartOptions{})
	if err != nil {
		return err
	}

	return nil
}

// SetupContainer prepares the image for the pipeline container.
func (c *client) SetupContainer(ctx context.Context, ctn *pipeline.Container) error {
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
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/images#Exists
	imageExists, err := images.Exists(c.Podman, _image, &images.ExistsOptions{})
	if err != nil {
		return err
	}

	// if the container image does
	// not exist on the host, we create it
	if !imageExists {
		// send API call to create the image
		return c.CreateImage(ctx, ctn)
	}

	return err
}

// TailContainer captures the logs for the pipeline container.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) TailContainer(ctx context.Context, ctn *pipeline.Container) (io.ReadCloser, error) {
	c.Logger.Tracef("tailing output for container %s", ctn.ID)

	// create options for capturing container logs
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/containers#LogOptions
	opts := new(containers.LogOptions).
		WithFollow(true).
		WithStderr(true).
		WithStdout(true).
		WithTimestamps(false)

	// create in-memory pipe for capturing logs
	rc, wc := io.Pipe()

	// create channels that will receive the log output
	stdOutChan := make(chan string)
	stdErrChan := make(chan string)

	// spawn a go routine to start capturing logs
	go func() {
		// send API call to capture the container logs
		//
		// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/containers#Logs
		err := containers.Logs(c.Podman, normalizeContainerName(ctn.Name), opts, stdOutChan, stdErrChan)
		if err != nil {
			return
		}

		// close PipeWriter
		wc.Close()

		// close channels
		close(stdErrChan)
		close(stdOutChan)
	}()

	// spawn a go routing to receive the stdout/stderr output
	go func() {
		for {
			select {
			case line, ok := <-stdOutChan:
				c.Logger.Tracef("Log: %s", line)

				// write current log to PipeWriter
				wc.Write([]byte(fmt.Sprintf("%s\n", line)))
				if !ok {
					stdOutChan = nil
				}
			case line, ok := <-stdErrChan:
				c.Logger.Tracef("Log: %s", line)

				// write current log to PipeWriter
				wc.Write([]byte(fmt.Sprintf("%s\n", line)))
				if !ok {
					stdErrChan = nil
				}
			}

			// if we receive an error exit
			if stdOutChan == nil && stdErrChan != nil {
				c.Logger.Tracef("Log: error occured, closing")

				return
			}
		}
	}()

	return rc, nil
}

// WaitContainer blocks until the pipeline container completes.
func (c *client) WaitContainer(ctx context.Context, ctn *pipeline.Container) error {
	c.Logger.Tracef("waiting for container %s", ctn.ID)

	// create wait options for the wait call
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/containers#WaitOptions
	waitOpts := new(containers.WaitOptions).
		WithCondition([]define.ContainerStatus{define.ContainerStateExited})

	// send API call to wait for the container completion
	//
	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/bindings/containers#Wait
	_, err := containers.Wait(c.Podman, normalizeContainerName(ctn.Name), waitOpts)

	return err
}

// ctnConfig is a helper function to
// generate the container config.
func specConf(ctn *pipeline.Container, id string, logger *logrus.Entry) *specgen.SpecGenerator {
	// create a new spec
	spec := specgen.NewSpecGenerator(image.Parse(ctn.Image), false)

	// the pod id that the container will join
	spec.Pod = id
	spec.Stdin = false
	spec.Terminal = false
	spec.Name = normalizeContainerName(ctn.Name)
	spec.Env = ctn.Environment

	// check if the entrypoint is provided
	if len(ctn.Entrypoint) > 0 {
		// add the entrypoint to container config
		spec.Entrypoint = ctn.Entrypoint
	}

	// check if the commands are provided
	if len(ctn.Commands) > 0 {
		// add commands to container config
		spec.Command = ctn.Commands
	}

	// check if the user is provided
	if len(ctn.User) > 0 {
		// add user to container config
		spec.User = ctn.User
	}

	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/specgen#LogConfig
	spec.LogConfiguration = &specgen.LogConfig{
		Driver: "json-file",
	}

	// https://pkg.go.dev/github.com/containers/podman/v3/pkg/specgen#ContainerResourceConfig
	for _, v := range ctn.Ulimits {
		spec.Rlimits = append(spec.Rlimits, specs.POSIXRlimit{
			Type: v.Name,
			Hard: uint64(v.Hard),
			Soft: uint64(v.Soft),
		})
	}

	return spec
}

// normalizeContainerName is a helper function to provide
// a cleaned up container name that can be used as the name
// for the container being spun up, since that will be its
// network addressable name by default inside the pod.
// TODO: see if we can keep using the more unique ID
// instead of Name and create an Alias instead.
func normalizeContainerName(name string) string {
	return strings.ReplaceAll(name, " ", "_")
}
