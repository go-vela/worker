// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	vol "github.com/go-vela/worker/internal/volume"
)

// CreateVolume creates the pipeline volume.
func (c *client) CreateVolume(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("creating volume for pipeline %s", b.ID)

	// create the workspace volume for the pod
	//
	// This is done due to the nature of how volumes works inside
	// the pod. Each container inside the pod can access and use
	// the same volume. This allows them to share this volume
	// throughout the life of the pod. However, to keep the
	// runtime behavior consistent, Vela uses an emtpyDir volume
	// because that volume only exists for the life
	// of the pod.
	//
	// More info:
	//   * https://kubernetes.io/docs/concepts/workloads/pods/pod/
	//   * https://kubernetes.io/docs/concepts/storage/volumes/#emptydir
	//
	// https://pkg.go.dev/k8s.io/api/core/v1#Volume
	workspaceVolume := v1.Volume{
		Name: b.ID,
		VolumeSource: v1.VolumeSource{
			EmptyDir: &v1.EmptyDirVolumeSource{},
		},
	}

	// create the workspace volumeMount for the pod
	//
	// https://pkg.go.dev/k8s.io/api/core/v1#VolumeMount
	workspaceVolumeMount := v1.VolumeMount{
		Name:      b.ID,
		MountPath: constants.WorkspaceMount,
	}

	// add the volume definition to the pod spec
	//
	// https://pkg.go.dev/k8s.io/api/core/v1#PodSpec
	c.Pod.Spec.Volumes = append(c.Pod.Spec.Volumes, workspaceVolume)

	// save the volumeMount to add to each of the containers in the pod spec later
	c.commonVolumeMounts = append(c.commonVolumeMounts, workspaceVolumeMount)

	// check if global host volumes were provided (VELA_RUNTIME_VOLUMES)
	if len(c.config.Volumes) > 0 {
		// iterate through all volumes provided
		for k, v := range c.config.Volumes {
			// parse the volume provided
			_volume := vol.Parse(v)
			_volumeName := fmt.Sprintf("%s_%d", b.ID, k)

			// add the volume to the set of pod volumes
			c.Pod.Spec.Volumes = append(c.Pod.Spec.Volumes, v1.Volume{
				Name: _volumeName,
				VolumeSource: v1.VolumeSource{
					HostPath: &v1.HostPathVolumeSource{
						Path: _volume.Source,
					},
				},
			})

			// save the volumeMounts for later addition to each container's mounts
			c.commonVolumeMounts = append(c.commonVolumeMounts, v1.VolumeMount{
				Name:      _volumeName,
				MountPath: _volume.Destination,
			})
		}
	}

	// TODO: extend c.config.Volumes to include container-specific volumes (container.Volumes)

	return nil
}

// InspectVolume inspects the pipeline volume.
func (c *client) InspectVolume(ctx context.Context, b *pipeline.Build) ([]byte, error) {
	c.Logger.Tracef("inspecting volume for pipeline %s", b.ID)

	// TODO: consider updating this command
	//
	// create output for inspecting volume
	output := []byte(
		fmt.Sprintf("$ kubectl get pod -o=jsonpath='{.spec.volumes}' %s\n", b.ID),
	)

	// marshal the volume information from the pod
	volume, err := json.MarshalIndent(c.Pod.Spec.Volumes, "", " ")
	if err != nil {
		return nil, err
	}

	return append(output, append(volume, "\n"...)...), nil
}

// RemoveVolume deletes the pipeline volume.
//
// Currently, this is comparable to a no-op because in Kubernetes the
// volume lives and dies with the pod it's attached to. However, Vela
// uses it to cleanup the volume definition for the pod.
func (c *client) RemoveVolume(ctx context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("removing volume for pipeline %s", b.ID)

	// remove the volume definition from the pod spec
	//
	// https://pkg.go.dev/k8s.io/api/core/v1#PodSpec
	c.Pod.Spec.Volumes = []v1.Volume{}
	c.commonVolumeMounts = []v1.VolumeMount{}

	return nil
}

// setupVolumeMounts generates the VolumeMounts for a given container.
//
//nolint:unparam // keep signature similar to Engine interface methods despite unused ctx and err
func (c *client) setupVolumeMounts(ctx context.Context, ctn *pipeline.Container) (
	volumeMounts []v1.VolumeMount,
	err error,
) {
	c.Logger.Tracef("setting up VolumeMounts for container %s", ctn.ID)

	// add workspace mount and any global host mounts (VELA_RUNTIME_VOLUMES)
	volumeMounts = append(volumeMounts, c.commonVolumeMounts...)

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
		for i, mount := range volumeMounts {
			// check if the path for the mount contains "/etc/ssl/certs"
			//
			// this is a soft check for mounting private cert bundles
			if strings.Contains(mount.MountPath, "/etc/ssl/certs") {
				// remove the private cert bundle mount from the host config
				volumeMounts = append(volumeMounts[:i], volumeMounts[i+1:]...)
			}
		}
	}
	//
	// -------------------- End of TODO: --------------------

	// TODO: extend volumeMounts based on ctn.Volumes

	return volumeMounts, nil
}
