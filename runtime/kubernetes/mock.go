// SPDX-License-Identifier: Apache-2.0

package kubernetes

// Everything in this file should only be used in test code.
// It is exported for use in tests of other packages.

import (
	"context"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"

	velav1alpha1 "github.com/go-vela/worker/runtime/kubernetes/apis/vela/v1alpha1"
	fakeVelaK8sClient "github.com/go-vela/worker/runtime/kubernetes/generated/clientset/versioned/fake"
)

// NewMock returns an Engine implementation that
// integrates with a Kubernetes runtime.
//
// This function is intended for running tests only.
//

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
	// https://pkg.go.dev/github.com/sirupsen/logrus#StandardLogger
	logger := logrus.StandardLogger()

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#NewEntry
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
	// https://pkg.go.dev/k8s.io/client-go/kubernetes/fake#NewSimpleClientset
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

	// The test is responsible for calling c.PodTracker.Start(ctx) if needed.
	// In some cases it is more convenient to call c.(MockKubernetesRuntime).StartPodTracker(ctx)

	return c, nil
}

// MockKubernetesRuntime makes it possible to use the client mocks in other packages.
//
// This interface is intended for running tests only.
type MockKubernetesRuntime interface {
	SetupMock() error
	MarkPodTrackerReady()
	StartPodTracker(context.Context)
	WaitForPodTrackerReady()
	WaitForPodCreate(string, string)
	SimulateResync(*v1.Pod)
	SimulateStatusUpdate(*v1.Pod, []v1.ContainerStatus) error
}

// SetupMock allows the Kubernetes runtime to perform additional Mock-related config.
// Many tests should call this right after they call runtime.SetupBuild (or executor.CreateBuild).
//
// This function is intended for running tests only.
func (c *client) SetupMock() error {
	// This assumes that c.Pod.ObjectMeta.Namespace and c.Pod.ObjectMeta.Name are filled in.
	return c.PodTracker.setupMockFor(c.Pod)
}

// MarkPodTrackerReady signals that PodTracker has been setup with ContainerTrackers.
//
// This function is intended for running tests only.
func (c *client) MarkPodTrackerReady() {
	close(c.PodTracker.Ready)
}

// StartPodTracker tells the podTracker it can start populating the cache.
//
// This function is intended for running tests only.
func (c *client) StartPodTracker(ctx context.Context) {
	c.PodTracker.Start(ctx)
}

// WaitForPodTrackerReady waits for PodTracker.Ready to be closed (which happens in AssembleBuild).
//
// This function is intended for running tests only.
func (c *client) WaitForPodTrackerReady() {
	<-c.PodTracker.Ready
}

// WaitForPodCreate waits for PodTracker.Ready to be closed (which happens in AssembleBuild).
//
// This function is intended for running tests only.
func (c *client) WaitForPodCreate(namespace, name string) {
	created := make(chan struct{})

	_, err := c.PodTracker.podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			select {
			case <-created:
				// not interested in any other create events.
				return
			default:
				break
			}

			var (
				pod *v1.Pod
				ok  bool
			)

			if pod, ok = obj.(*v1.Pod); !ok {
				return
			}

			if pod.GetNamespace() == namespace && pod.GetName() == name {
				close(created)
			}
		},
	})
	if err != nil {
		return
	}

	<-created
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

// SimulateUpdate simulates an update event from the k8s API.
//
// This function is intended for running tests only.
func (c *client) SimulateStatusUpdate(pod *v1.Pod, containerStatuses []v1.ContainerStatus) error {
	// We have to have a full copy here because the k8s client Mock
	// replaces the pod it is storing, it does not just update the status.
	updatedPod := pod.DeepCopy()
	updatedPod.Status.ContainerStatuses = containerStatuses

	_, err := c.Kubernetes.CoreV1().Pods(pod.GetNamespace()).
		UpdateStatus(
			context.Background(),
			updatedPod,
			metav1.UpdateOptions{},
		)

	return err
}
