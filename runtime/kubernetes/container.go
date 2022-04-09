// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/internal/image"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

// InspectContainer inspects the pipeline container.
func (c *client) InspectContainer(ctx context.Context, ctn *pipeline.Container) error {
	c.Logger.Tracef("inspecting container %s", ctn.ID)

	// get the pod from the local cache, which the Informer keeps up-to-date
	pod, err := c.PodTracker.PodLister.Pods(c.config.Namespace).Get(c.Pod.ObjectMeta.Name)
	if err != nil {
		return err
	}

	// iterate through each container in the pod
	for _, cst := range pod.Status.ContainerStatuses {
		// check if the container has a matching ID
		//
		// https://pkg.go.dev/k8s.io/api/core/v1?tab=doc#ContainerStatus
		if !strings.EqualFold(cst.Name, ctn.ID) {
			// skip container if it's not a matching ID
			continue
		}

		// set the step exit code
		ctn.ExitCode = int(cst.State.Terminated.ExitCode)

		break
	}

	return nil
}

// RemoveContainer deletes (kill, remove) the pipeline container.
// This is a no-op for kubernetes. RemoveBuild handles deleting the pod.
func (c *client) RemoveContainer(ctx context.Context, ctn *pipeline.Container) error {
	c.Logger.Tracef("no-op: removing container %s", ctn.ID)

	return nil
}

// RunContainer creates and starts the pipeline container.
func (c *client) RunContainer(ctx context.Context, ctn *pipeline.Container, b *pipeline.Build) error {
	c.Logger.Tracef("running container %s", ctn.ID)
	// parse image from step
	_image, err := image.ParseWithError(ctn.Image)
	if err != nil {
		return err
	}

	// set the pod container image to the parsed step image
	// (-1 to convert to 0-based index, -1 for init which isn't a container)
	c.Pod.Spec.Containers[ctn.Number-2].Image = _image

	// send API call to patch the pod with the new container image
	//
	// https://pkg.go.dev/k8s.io/client-go/kubernetes/typed/core/v1?tab=doc#PodInterface
	// nolint: contextcheck // ignore non-inherited new context
	_, err = c.Kubernetes.CoreV1().Pods(c.config.Namespace).Patch(
		context.Background(),
		c.Pod.ObjectMeta.Name,
		types.StrategicMergePatchType,
		[]byte(fmt.Sprintf(imagePatch, ctn.ID, _image)),
		metav1.PatchOptions{},
	)
	if err != nil {
		return err
	}

	return nil
}

// SetupContainer prepares the image for the pipeline container.
func (c *client) SetupContainer(ctx context.Context, ctn *pipeline.Container) error {
	c.Logger.Tracef("setting up for container %s", ctn.ID)

	// create the container object for the pod
	//
	// https://pkg.go.dev/k8s.io/api/core/v1?tab=doc#Container
	container := v1.Container{
		Name: ctn.ID,
		// create the container with the kubernetes/pause image
		//
		// This is done due to the nature of how containers are
		// executed inside the pod. Kubernetes will attempt to
		// start and run all containers in the pod at once. We
		// want to control the execution of the containers
		// inside the pod so we use the pause image as the
		// default for containers, and then sequentially patch
		// the containers with the proper image.
		//
		// https://hub.docker.com/r/kubernetes/pause
		Image:           image.Parse("kubernetes/pause:latest"),
		Env:             []v1.EnvVar{},
		Stdin:           false,
		StdinOnce:       false,
		TTY:             false,
		WorkingDir:      ctn.Directory,
		SecurityContext: &v1.SecurityContext{},
	}

	// handle the container pull policy (This cannot be updated like the image can)
	switch ctn.Pull {
	case constants.PullAlways:
		// set the pod container pull policy to always
		container.ImagePullPolicy = v1.PullAlways
	case constants.PullNever:
		// set the pod container pull policy to never
		container.ImagePullPolicy = v1.PullNever
	case constants.PullOnStart:
		// set the pod container pull policy to always
		//
		// if the pipeline container image should be pulled on start, than
		// we force Kubernetes to pull the image on start with the always
		// pull policy for the pod container
		container.ImagePullPolicy = v1.PullAlways
	case constants.PullNotPresent:
		fallthrough
	default:
		// default the pod container pull policy to if not present
		container.ImagePullPolicy = v1.PullIfNotPresent
	}

	// fill in the VolumeMounts including workspaceMount
	volumeMounts, err := c.setupVolumeMounts(ctx, ctn)
	if err != nil {
		return err
	}

	container.VolumeMounts = volumeMounts

	// check if the image is allowed to run privileged
	for _, pattern := range c.config.Images {
		privileged, err := image.IsPrivilegedImage(ctn.Image, pattern)
		if err != nil {
			return err
		}

		container.SecurityContext.Privileged = &privileged
	}

	if c.PipelinePodTemplate != nil && c.PipelinePodTemplate.Spec.Container != nil {
		securityContext := c.PipelinePodTemplate.Spec.Container.SecurityContext

		// TODO: add more SecurityContext options (runAsUser, runAsNonRoot, sysctls)
		if securityContext != nil && securityContext.Capabilities != nil {
			container.SecurityContext.Capabilities = securityContext.Capabilities
		}
	}

	// Executor.CreateBuild extends the environment AFTER calling Runtime.SetupBuild.
	// So, configure the environment as late as possible (just before pod creation).

	// check if the entrypoint is provided
	if len(ctn.Entrypoint) > 0 {
		// add entrypoint to container config
		container.Args = ctn.Entrypoint
	}

	// check if the commands are provided
	if len(ctn.Commands) > 0 {
		// add commands to container config
		container.Args = append(container.Args, ctn.Commands...)
	}

	// add the container definition to the pod spec
	//
	// https://pkg.go.dev/k8s.io/api/core/v1?tab=doc#PodSpec
	c.Pod.Spec.Containers = append(c.Pod.Spec.Containers, container)

	return nil
}

// setupContainerEnvironment adds env vars to the Pod spec for a container.
// Call this just before pod creation to capture as many env changes as possible.
func (c *client) setupContainerEnvironment(ctn *pipeline.Container) error {
	c.Logger.Tracef("setting up environment for container %s", ctn.ID)

	// get the matching container spec
	// (-1 to convert to 0-based index, -1 for injected init container)
	container := &c.Pod.Spec.Containers[ctn.Number-2]
	if !strings.EqualFold(container.Name, ctn.ID) {
		return fmt.Errorf("wrong container! got %s instead of %s", container.Name, ctn.ID)
	}

	// check if the environment is provided
	if len(ctn.Environment) > 0 {
		// iterate through each element in the container environment
		for k, v := range ctn.Environment {
			// add key/value environment to container config
			container.Env = append(container.Env, v1.EnvVar{Name: k, Value: v})
		}
	}

	return nil
}

// TailContainer captures the logs for the pipeline container.
func (c *client) TailContainer(ctx context.Context, ctn *pipeline.Container) (io.ReadCloser, error) {
	c.Logger.Tracef("tailing output for container %s", ctn.ID)

	// create object to store container logs
	var logs io.ReadCloser

	// get the containerTracker for this container
	containerTracker, ok := c.PodTracker.Containers[ctn.ID]
	if !ok {
		return nil, fmt.Errorf("containerTracker is missing for %s", ctn.ID)
	}

	// wrap the bytes.Reader in an io.NopCloser
	logs = io.NopCloser(containerTracker.Logs())

	logsError := containerTracker.LogsError
	// io.EOF means that all logs have been captured.
	if logsError != nil && logsError != io.EOF && logsError != TruncatedLogs {
		// TODO: modify the executor to accept record partial logs before the failure
		return logs, logsError
	}

	return logs, nil
}

// streamContainerLogs streams the logs to a cache up to a maxLogSize, restarting the stream as needed.
// streamContainerLogs is designed to run in its own goroutine.
func (p podTracker) streamContainerLogs(ctx context.Context, ctnTracker *containerTracker, maxLogSize uint) {
	// create function for periodically capturing
	// the logs from the container with backoff
	logsFunc := func() (bool, error) {
		// create options for capturing the logs from the container
		//
		// https://pkg.go.dev/k8s.io/api/core/v1?tab=doc#PodLogOptions
		opts := &v1.PodLogOptions{
			Container:  ctnTracker.Name,
			Follow:     true,
			Timestamps: false,
		}

		// send API call to capture stream of container logs
		//
		// https://pkg.go.dev/k8s.io/client-go/kubernetes/typed/core/v1?tab=doc#PodExpansion
		// ->
		// https://pkg.go.dev/k8s.io/client-go/rest?tab=doc#Request.Stream
		stream, err := p.Kubernetes.CoreV1().
			Pods(p.Namespace).
			GetLogs(p.Name, opts).
			Stream(context.Background())
		if err != nil {
			p.Logger.Errorf("failed to stream logs for %s, %v", p.TrackedPod, err)

			// retry the API call
			return false, nil
		}

		defer stream.Close()

		// create new reader from the container output
		reader := bufio.NewReader(stream)

		// this loop is loosely based on github.com/stern/stern.Tail.ConsumeRequest()
		for {
			// for each line in stream
			line, err := reader.ReadBytes('\n')
			if len(line) != 0 {
				// cache the log line
				ctnTracker.logs = append(ctnTracker.logs, line...)
			}

			if err != nil {
				// save err even if its io.EOF as EOF indicates all logs were read.
				ctnTracker.LogsError = err

				if err != io.EOF {
					p.Logger.Errorf("error while streaming logs for %s, %v", p.TrackedPod, err)

					// we did not reach the end of the logs so let's try again.
					// If this proves problematic, we might need to de-dup log lines
					// which might require using opts.Timestamps or opts.SinceSeconds.
					return false, nil
				}

				// hooray! we reached io.EOF (the end of the logs)
				break
			}

			// there are more logs to read
			// check whether we've reached the maximum log size
			if maxLogSize > 0 && uint(len(ctnTracker.logs)) >= maxLogSize {
				p.Logger.Trace("maximum log size reached")

				ctnTracker.LogsError = TruncatedLogs
				ctnTracker.logs = append(ctnTracker.logs, []byte("LOGS TRUNCATED: Vela Runtime MaxLogSize exceeded.\n")...)
				break
			}
		}

		// check if we have container logs from the stream
		if len(ctnTracker.logs) > 0 {
			// no more logs to stream
			return true, nil
		}

		// no logs are available, so try again
		return false, nil
	}

	// create backoff object for capturing the logs
	// from the container with periodic backoff
	//
	// https://pkg.go.dev/k8s.io/apimachinery/pkg/util/wait?tab=doc#Backoff
	backoff := wait.Backoff{
		Duration: 1 * time.Second,
		Factor:   2.0,
		Jitter:   0.25,
		Steps:    10,
		Cap:      2 * time.Minute,
	}

	p.Logger.Tracef("capturing logs with exponential backoff for container %s", ctnTracker.Name)
	// perform the function to capture logs with periodic backoff
	//
	// https://pkg.go.dev/k8s.io/apimachinery/pkg/util/wait?tab=doc#ExponentialBackoff
	err := wait.ExponentialBackoffWithContext(ctx, backoff, logsFunc)
	if err != nil {
		p.Logger.Errorf("exponential backoff error while streaming logs for %s, %v", ctnTracker.Name, err)

		if ctnTracker.LogsError == nil {
			ctnTracker.LogsError = err
		}
	}
}

// WaitContainer blocks until the pipeline container completes.
func (c *client) WaitContainer(ctx context.Context, ctn *pipeline.Container) error {
	c.Logger.Tracef("waiting for container %s", ctn.ID)

	// get the containerTracker for this container
	tracker, ok := c.PodTracker.Containers[ctn.ID]
	if !ok {
		return fmt.Errorf("containerTracker is missing for %s", ctn.ID)
	}

	// wait for the container terminated signal
	<-tracker.Terminated

	return nil
}

// inspectContainerStatuses signals when a container reaches a terminal state.
func (p podTracker) inspectContainerStatuses(pod *v1.Pod) {
	// check if the pod is in a pending state
	//
	// https://pkg.go.dev/k8s.io/api/core/v1?tab=doc#PodStatus
	if pod.Status.Phase == v1.PodPending {
		// nothing to inspect if pod is in a pending state
		return
	}

	// iterate through each container in the pod
	for _, cst := range pod.Status.ContainerStatuses {
		// get the containerTracker for this container
		tracker, ok := p.Containers[cst.Name]
		if !ok {
			// unknown container
			continue
		}

		// check if the container is in a terminated state
		//
		// https://pkg.go.dev/k8s.io/api/core/v1?tab=doc#ContainerState
		if cst.State.Terminated != nil {
			// && len(cst.State.Terminated.Reason) > 0 {
			// WaitContainer used to check Terminated.Reason as well.
			// if that is still needed, then we can add that check here
			// or retrieve the pod with something like this in WaitContainer:
			// c.PodTracker.PodLister.Pods(c.config.Namespace).Get(c.Pod.GetName())
			tracker.terminatedOnce.Do(func() {
				// let WaitContainer know the container is terminated
				close(tracker.Terminated)
			})
		}
	}
}
