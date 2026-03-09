// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"

	"github.com/go-vela/server/compiler/types/pipeline"
)

// CreateNetwork creates the pipeline network.
func (c *client) CreateNetwork(_ context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("creating network for pipeline %s", b.ID)

	// create the network for the pod
	//
	// This is done due to the nature of how networking works inside the
	// pod. Each container inside the pod shares the same network IP and
	// port space. This allows them to communicate with each other via
	// localhost. However, to keep the runtime behavior consistent,
	// Vela adds DNS entries for each container that requires it.
	//
	// More info:
	//   * https://kubernetes.io/docs/concepts/workloads/pods/pod/
	//   * https://kubernetes.io/docs/concepts/services-networking/add-entries-to-pod-etc-hosts-with-host-aliases/
	//
	// https://pkg.go.dev/k8s.io/api/core/v1#HostAlias
	network := v1.HostAlias{
		IP:        "127.0.0.1",
		Hostnames: []string{},
	}

	// iterate through all services in the pipeline
	for _, service := range b.Services {
		// create the host entry for the pod container aliases
		host := fmt.Sprintf("%s.local", service.Name)

		// add the host entry to the pod container aliases
		network.Hostnames = append(network.Hostnames, host)
	}

	// iterate through all steps in the pipeline
	for _, step := range b.Steps {
		// skip all steps not running in detached mode
		if !step.Detach {
			continue
		}

		// create the host entry for the pod container aliases
		host := fmt.Sprintf("%s.local", step.Name)

		// add the host entry to the pod container aliases
		network.Hostnames = append(network.Hostnames, host)
	}

	// iterate through all stages in the pipeline
	for _, stage := range b.Stages {
		// iterate through all steps in each stage
		for _, step := range stage.Steps {
			// skip all steps not running in detached mode
			if !step.Detach {
				continue
			}

			// create the host entry for the pod container aliases
			host := fmt.Sprintf("%s.local", step.Name)

			// add the host entry to the pod container aliases
			network.Hostnames = append(network.Hostnames, host)
		}
	}

	// add the network definition to the pod spec
	//
	// https://pkg.go.dev/k8s.io/api/core/v1#PodSpec
	c.Pod.Spec.HostAliases = append(c.Pod.Spec.HostAliases, network)

	return nil
}

// InspectNetwork inspects the pipeline network.
func (c *client) InspectNetwork(_ context.Context, b *pipeline.Build) ([]byte, error) {
	c.Logger.Tracef("inspecting network for pipeline %s", b.ID)

	// TODO: consider updating this command
	//
	// create output for inspecting volume
	output :=
		fmt.Appendf(nil, "$ kubectl get pod -o=jsonpath='{.spec.hostAliases}' %s\n", b.ID)

	// marshal the network information from the pod
	network, err := json.MarshalIndent(c.Pod.Spec.HostAliases, "", " ")
	if err != nil {
		return output, err
	}

	return append(output, append(network, "\n"...)...), nil
}

// RemoveNetwork deletes the pipeline network.
//
// Currently, this is comparable to a no-op because in Kubernetes the
// network lives and dies with the pod it's attached to. However, Vela
// uses it to cleanup the network definition for the pod.
func (c *client) RemoveNetwork(_ context.Context, b *pipeline.Build) error {
	c.Logger.Tracef("removing network for pipeline %s", b.ID)

	// remove the network definition from the pod spec
	//
	// https://pkg.go.dev/k8s.io/api/core/v1#PodSpec
	c.Pod.Spec.HostAliases = []v1.HostAlias{}

	return nil
}
