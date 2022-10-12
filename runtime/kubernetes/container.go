// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
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
	pod, err := c.PodTracker.PodLister.
		Pods(c.config.Namespace).
		Get(c.Pod.ObjectMeta.Name)
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

		// avoid a panic if the build ends without terminating all containers
		if cst.State.Terminated == nil {
			// steps that were not executed will still be "running" the pause image as expected.
			if cst.Image == pauseImage || cst.Image == image.Parse(pauseImage) {
				return nil
			}

			return fmt.Errorf("expected container %s to be terminated, got %v", ctn.ID, cst.State)
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
func (c *client) RunContainer(ctx context.Context, ctn *pipeline.Container, b *pipeline.Build, r *library.Repo) error {
	c.Logger.Tracef("running container %s", ctn.ID)
	// parse image from step
	_image, err := image.ParseWithError(ctn.Image)
	if err != nil {
		return err
	}

	// set the pod container image to the parsed step image
	c.Pod.Spec.Containers[c.containersLookup[ctn.ID]].Image = _image

	// send API call to patch the pod with the new container image
	//
	// https://pkg.go.dev/k8s.io/client-go/kubernetes/typed/core/v1?tab=doc#PodInterface
	_, err = c.Kubernetes.CoreV1().Pods(c.config.Namespace).Patch(
		ctx,
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
func (c *client) SetupContainer(ctx context.Context, ctn *pipeline.Container, r *library.Repo) error {
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
		Image:           image.Parse(pauseImage),
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

		// ensure repo is trusted and therefore allowed to run privileged containers
		if r != nil && !r.GetTrusted() {
			return errors.New("cannot use privileged image in a non-trusted repo")
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

	// record the index for this container
	c.containersLookup[ctn.ID] = len(c.Pod.Spec.Containers)

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
	container := &c.Pod.Spec.Containers[c.containersLookup[ctn.ID]]
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

	// create function for periodically capturing
	// the logs from the container with backoff
	logsFunc := func() (bool, error) {
		// create options for capturing the logs from the container
		//
		// https://pkg.go.dev/k8s.io/api/core/v1?tab=doc#PodLogOptions
		opts := &v1.PodLogOptions{
			Container:  ctn.ID,
			Follow:     true,
			Timestamps: false,
		}

		// send API call to capture stream of container logs
		//
		// https://pkg.go.dev/k8s.io/client-go/kubernetes/typed/core/v1?tab=doc#PodExpansion
		// ->
		// https://pkg.go.dev/k8s.io/client-go/rest?tab=doc#Request.Stream
		stream, err := c.Kubernetes.CoreV1().
			Pods(c.config.Namespace).
			GetLogs(c.Pod.ObjectMeta.Name, opts).
			Stream(ctx)
		if err != nil {
			c.Logger.Errorf("error while requesting pod/logs stream for container %s: %v", ctn.ID, err)
			return false, nil
		}

		// create temporary reader to ensure logs are available
		reader := bufio.NewReader(stream)

		// peek at container logs from the stream
		bytes, err := reader.Peek(5)
		if err != nil {
			//nolint:nilerr // ignore nil return
			// skip so we resend API call to capture stream
			return false, nil
		}

		// check if we have container logs from the stream
		if len(bytes) > 0 {
			// set the logs to the reader
			logs = io.NopCloser(reader)
			return true, nil
		}

		// no logs are available
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

	c.Logger.Tracef("capturing logs with exponential backoff for container %s", ctn.ID)
	// perform the function to capture logs with periodic backoff
	//
	// https://pkg.go.dev/k8s.io/apimachinery/pkg/util/wait?tab=doc#ExponentialBackoff
	err := wait.ExponentialBackoffWithContext(ctx, backoff, logsFunc)
	if err != nil {
		c.Logger.Errorf("exponential backoff error while tailing container %s: %v", ctn.ID, err)
		return nil, err
	}

	return logs, nil
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
func (p *podTracker) inspectContainerStatuses(pod *v1.Pod) {
	// check if the pod is in a pending state
	//
	// https://pkg.go.dev/k8s.io/api/core/v1?tab=doc#PodStatus
	if pod.Status.Phase == v1.PodPending {
		p.Logger.Debugf("skipping container status inspection as pod %s is pending", p.TrackedPod)

		// nothing to inspect if pod is in a pending state
		return
	}

	// iterate through each container in the pod
	for _, cst := range pod.Status.ContainerStatuses {
		// get the containerTracker for this container
		tracker, ok := p.Containers[cst.Name]
		if !ok {
			// unknown container (probably a sidecar injected by an admissions controller)
			p.Logger.Debugf("ignoring untracked container %s from pod %s", cst.Name, p.TrackedPod)

			continue
		}

		// cst.State has details about the cst.Image's exit.
		// cst.LastTerminationState has details about the kubernetes/pause image's exit.
		// cst.RestartCount is 1 at exit due to switch from kubernetes/pause to final image.

		// check if the container is in a terminated state
		//
		// https://pkg.go.dev/k8s.io/api/core/v1?tab=doc#ContainerState
		if cst.State.Terminated != nil {
			tracker.terminatedOnce.Do(func() {
				p.Logger.Debugf("container completed: %s in pod %s, %v", cst.Name, p.TrackedPod, cst)

				// let WaitContainer know the container is terminated
				close(tracker.Terminated)
			})
		}
	}
}
