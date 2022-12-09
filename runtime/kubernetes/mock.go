// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	velav1alpha1 "github.com/go-vela/worker/runtime/kubernetes/apis/vela/v1alpha1"
	fakeVelaK8sClient "github.com/go-vela/worker/runtime/kubernetes/generated/clientset/versioned/fake"
)

// NewMock returns an Engine implementation that
// integrates with a Kubernetes runtime.
//
// This function is intended for running tests only.
//
//nolint:revive // ignore returning unexported client
func NewMock(_pod *v1.Pod, opts ...ClientOpt) (*client, error) {
	// create new Kubernetes client
	c := new(client)

	// create new fields
	c.config = new(config)
	c.Pod = new(v1.Pod)

	c.containersLookup = map[string]int{}
	for i, ctn := range _pod.Spec.Containers {
		c.containersLookup[ctn.Name] = i
	}

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#StandardLogger
	logger := logrus.StandardLogger()

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#NewEntry
	c.Logger = logrus.NewEntry(logger)

	// set the Kubernetes namespace in the runtime client
	c.config.Namespace = "test"

	// set the Kubernetes pod in the runtime client
	c.Pod = _pod.DeepCopy()
	c.Pod.SetResourceVersion("0")

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// set the Kubernetes fake client in the runtime client
	//
	// https://pkg.go.dev/k8s.io/client-go/kubernetes/fake?tab=doc#NewSimpleClientset
	c.Kubernetes = fake.NewSimpleClientset(c.Pod)

	// set the VelaKubernetes fake client in the runtime client
	c.VelaKubernetes = fakeVelaK8sClient.NewSimpleClientset(
		&velav1alpha1.PipelinePodsTemplate{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: c.config.Namespace,
				Name:      "mock-pipeline-pods-template",
			},
		},
	)

	// set the PodTracker (normally populated in SetupBuild)
	tracker, err := mockPodTracker(c.Logger, c.Kubernetes, c.Pod)
	if err != nil {
		return c, err
	}

	c.PodTracker = tracker

	// The test is responsible for calling c.PodTracker.Start() if needed

	return c, nil
}

// MockKubernetesRuntime makes it possible to use the client mocks in other packages.
//
// This interface is intended for running tests only.
type MockKubernetesRuntime interface {
	MarkPodTrackerReady()
	SimulateResync(*v1.Pod)
}

// MarkPodTrackerReady signals that PodTracker has been setup with ContainerTrackers.
//
// This function is intended for running tests only.
func (c *client) MarkPodTrackerReady() {
	close(c.PodTracker.Ready)
}

// SimulateResync simulates an resync where the PodTracker refreshes its cache.
// This resync is from oldPod to runtime.Pod. If nil, oldPod defaults to runtime.Pod.
//
// This function is intended for running tests only.
func (c *client) SimulateResync(oldPod *v1.Pod) {
	if oldPod == nil {
		oldPod = c.Pod
	}

	oldPod = oldPod.DeepCopy()
	oldPod.SetResourceVersion("older")

	// simulate a re-sync/PodUpdate event
	c.PodTracker.HandlePodUpdate(oldPod, c.Pod)
}
