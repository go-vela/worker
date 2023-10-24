// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/worker/runtime/kubernetes/apis/vela/v1alpha1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	// The k8s libraries have some quirks around yaml marshaling (see opts.go).
	// So, just use the same library for all kubernetes-related YAML.
	"sigs.k8s.io/yaml"
)

// InspectBuild displays details about the pod for the init step.
func (c *client) InspectBuild(ctx context.Context, b *pipeline.Build) ([]byte, error) {
	c.Logger.Tracef("inspecting build pod for pipeline %s", b.ID)

	output := []byte(fmt.Sprintf("> Inspecting pod for pipeline %s\n", b.ID))

	// TODO: The environment gets populated in AssembleBuild, after InspectBuild runs.
	//       But, we should make sure that secrets can't be leaked here anyway.
	buildOutput, err := yaml.Marshal(c.Pod)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to serialize pod: %w", err)
	}

	output = append(output, buildOutput...)

	// TODO: make other k8s Inspect* funcs no-ops (prefer this method):
	// 	     InspectVolume, InspectImage, InspectNetwork
	return output, nil
}

// SetupBuild prepares the pod metadata for the pipeline build.
func (c *client) SetupBuild(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("setting up for build %s", b.ID)

	if c.PipelinePodTemplate == nil {
		if len(c.config.PipelinePodsTemplateName) > 0 {
			podsTemplateResponse, err := c.VelaKubernetes.VelaV1alpha1().
				PipelinePodsTemplates(c.config.Namespace).
				Get(ctx, c.config.PipelinePodsTemplateName, metav1.GetOptions{})
			if err != nil {
				return err
			}

			// save the PipelinePodTemplate to use later in SetupContainer and other Setup methods
			c.PipelinePodTemplate = &podsTemplateResponse.Spec.Template
		} else {
			c.PipelinePodTemplate = &v1alpha1.PipelinePodTemplate{}
		}
	}

	// These labels will be used to call k8s watch APIs.
	labels := map[string]string{"pipeline": b.ID}

	if c.PipelinePodTemplate.Metadata.Labels != nil {
		// merge the template labels into the worker-defined labels.
		for k, v := range c.PipelinePodTemplate.Metadata.Labels {
			// do not allow overwriting any of the worker-defined labels.
			if _, ok := labels[k]; ok {
				continue
			}

			labels[k] = v
		}
	}

	// create the object metadata for the pod
	//
	// https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1?tab=doc#ObjectMeta
	c.Pod.ObjectMeta = metav1.ObjectMeta{
		Name:        b.ID,
		Namespace:   c.config.Namespace, // this is used by the podTracker
		Labels:      labels,
		Annotations: c.PipelinePodTemplate.Metadata.Annotations,
	}

	// TODO: Vela admin defined worker-specific: AutomountServiceAccountToken

	if c.PipelinePodTemplate.Spec.NodeSelector != nil {
		c.Pod.Spec.NodeSelector = c.PipelinePodTemplate.Spec.NodeSelector
	}

	if c.PipelinePodTemplate.Spec.Tolerations != nil {
		c.Pod.Spec.Tolerations = c.PipelinePodTemplate.Spec.Tolerations
	}

	if c.PipelinePodTemplate.Spec.Affinity != nil {
		c.Pod.Spec.Affinity = c.PipelinePodTemplate.Spec.Affinity
	}

	// create the restart policy for the pod
	//
	// https://pkg.go.dev/k8s.io/api/core/v1?tab=doc#RestartPolicy
	c.Pod.Spec.RestartPolicy = v1.RestartPolicyNever

	if len(c.PipelinePodTemplate.Spec.DNSPolicy) > 0 {
		c.Pod.Spec.DNSPolicy = c.PipelinePodTemplate.Spec.DNSPolicy
	}

	if c.PipelinePodTemplate.Spec.DNSConfig != nil {
		c.Pod.Spec.DNSConfig = c.PipelinePodTemplate.Spec.DNSConfig
	}

	if c.PipelinePodTemplate.Spec.SecurityContext != nil {
		if c.Pod.Spec.SecurityContext == nil {
			c.Pod.Spec.SecurityContext = &v1.PodSecurityContext{}
		}

		if c.PipelinePodTemplate.Spec.SecurityContext.RunAsNonRoot != nil {
			c.Pod.Spec.SecurityContext.RunAsNonRoot = c.PipelinePodTemplate.Spec.SecurityContext.RunAsNonRoot
		}

		if c.PipelinePodTemplate.Spec.SecurityContext.Sysctls != nil {
			c.Pod.Spec.SecurityContext.Sysctls = c.PipelinePodTemplate.Spec.SecurityContext.Sysctls
		}
	}

	// initialize the PodTracker now that we have a Pod for it to track
	tracker, err := newPodTracker(c.Logger, c.Kubernetes, c.Pod, time.Second*30)
	if err != nil {
		return err
	}

	c.PodTracker = tracker

	return nil
}

// StreamBuild initializes log/event streaming for build.
func (c *client) StreamBuild(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("streaming build %s", b.ID)

	select {
	case <-ctx.Done():
		// bail out, as build timed out or was canceled.
		return nil
	case <-c.PodTracker.Ready:
		// AssembleBuild signaled that the PodTracker is ready.
		break
	}

	// Populate the PodTracker caches before creating the pipeline pod
	c.PodTracker.Start(ctx)

	return nil
}

// AssembleBuild finalizes the pipeline build setup.
// This creates the pod in kubernetes for the pipeline build.
// After creation, image is the only container field we can edit in kubernetes.
// So, all environment, volume, and other container metadata must be setup
// before running AssembleBuild.
func (c *client) AssembleBuild(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("assembling build %s", b.ID)

	var err error

	// last minute Environment setup
	for _, _service := range b.Services {
		err = c.setupContainerEnvironment(_service)
		if err != nil {
			return err
		}
	}

	for _, _stage := range b.Stages {
		// TODO: remove hardcoded reference
		if _stage.Name == "init" {
			continue
		}

		for _, _step := range _stage.Steps {
			err = c.setupContainerEnvironment(_step)
			if err != nil {
				return err
			}
		}
	}

	for _, _step := range b.Steps {
		// TODO: remove hardcoded reference
		if _step.Name == "init" {
			continue
		}

		err = c.setupContainerEnvironment(_step)
		if err != nil {
			return err
		}
	}

	for _, _secret := range b.Secrets {
		if _secret.Origin.Empty() {
			continue
		}

		err = c.setupContainerEnvironment(_secret.Origin)
		if err != nil {
			return err
		}
	}

	// setup containerTrackers now that all containers are defined.
	c.PodTracker.TrackContainers(c.Pod.Spec.Containers)

	// send signal to StreamBuild now that PodTracker is ready to be started.
	close(c.PodTracker.Ready)

	// wait for the PodTracker caches to populate before creating the pipeline pod.
	if ok := cache.WaitForCacheSync(ctx.Done(), c.PodTracker.PodSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	// If the api call to create the pod fails, the pod might
	// partially exist. So, set this first to make sure all
	// remnants get deleted.
	c.createdPod = true

	c.Logger.Infof("creating pod %s", c.Pod.ObjectMeta.Name)
	// send API call to create the pod
	//
	// https://pkg.go.dev/k8s.io/client-go/kubernetes/typed/core/v1?tab=doc#PodInterface
	_, err = c.Kubernetes.CoreV1().
		Pods(c.config.Namespace).
		Create(ctx, c.Pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// RemoveBuild deletes (kill, remove) the pipeline build metadata.
// This deletes the kubernetes pod.
func (c *client) RemoveBuild(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("removing build %s", b.ID)

	// PodTracker gets created in SetupBuild before pod is created
	defer func() {
		// check for nil as RemoveBuild may get called multiple times
		if c.PodTracker != nil {
			c.PodTracker.Stop()
			c.PodTracker = nil
		}
	}()

	if !c.createdPod {
		// nothing to do
		return nil
	}

	// create variables for the delete options
	//
	// This is necessary because the delete options
	// expect all values to be passed by reference.
	var (
		period = int64(0)
		policy = metav1.DeletePropagationForeground
	)

	// create options for removing the pod
	//
	// https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1?tab=doc#DeleteOptions
	opts := metav1.DeleteOptions{
		GracePeriodSeconds: &period,
		// https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1?tab=doc#DeletionPropagation
		PropagationPolicy: &policy,
	}

	c.Logger.Infof("removing pod %s", c.Pod.ObjectMeta.Name)
	// send API call to delete the pod
	err := c.Kubernetes.CoreV1().
		Pods(c.config.Namespace).
		Delete(ctx, c.Pod.ObjectMeta.Name, opts)
	if err != nil {
		return err
	}

	c.Pod = &v1.Pod{}
	c.createdPod = false

	return nil
}
