// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"context"

	"github.com/sirupsen/logrus"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	velav1alpha1 "github.com/go-vela/worker/runtime/kubernetes/apis/vela/v1alpha1"
	velaK8sClient "github.com/go-vela/worker/runtime/kubernetes/generated/clientset/versioned"
	fakeVelaK8sClient "github.com/go-vela/worker/runtime/kubernetes/generated/clientset/versioned/fake"
)

type config struct {
	// specifies the config file to use for the Kubernetes client
	File string
	// specifies the namespace to use for the Kubernetes client
	Namespace string
	// specifies a list of privileged images to use for the Kubernetes client
	Images []string
	// specifies a list of host volumes to use for the Kubernetes client
	Volumes []string
	// PipelinePodsTemplateName has the name of the PipelinePodTemplate to retrieve from the Namespace
	PipelinePodsTemplateName string
}

type client struct {
	config *config
	// https://pkg.go.dev/k8s.io/client-go/kubernetes#Interface
	Kubernetes kubernetes.Interface
	// VelaKubernetes is a client for custom Vela CRD-based APIs
	VelaKubernetes velaK8sClient.Interface
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	Logger *logrus.Entry
	// https://pkg.go.dev/k8s.io/api/core/v1#Pod
	Pod *v1.Pod
	// containersLookup maps the container name to its index in Containers
	containersLookup map[string]int
	// PodTracker wraps the Kubernetes client to simplify watching the pod for changes
	PodTracker *podTracker
	// PipelinePodTemplate has default values to be used in Setup* methods
	PipelinePodTemplate *velav1alpha1.PipelinePodTemplate
	// commonVolumeMounts includes workspace mount and any global host mounts (VELA_RUNTIME_VOLUMES)
	commonVolumeMounts []v1.VolumeMount
	// indicates when the pod has been created in kubernetes
	createdPod bool
}

// New returns an Engine implementation that
// integrates with a Kubernetes runtime.
//
//nolint:revive // ignore returning unexported client
func New(opts ...ClientOpt) (*client, error) {
	// create new Kubernetes client
	c := new(client)

	// create new fields
	c.config = new(config)
	c.Pod = new(v1.Pod)
	c.containersLookup = map[string]int{}

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#StandardLogger
	logger := logrus.StandardLogger()

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#NewEntry
	c.Logger = logrus.NewEntry(logger)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// use the current context in kubeconfig
	//
	// when no kube config is provided create InClusterConfig
	// else use out of cluster config option
	var (
		config *rest.Config
		err    error
	)

	if c.config.File == "" {
		// https://pkg.go.dev/k8s.io/client-go/rest?tab=doc#InClusterConfig
		config, err = rest.InClusterConfig()
		if err != nil {
			c.Logger.Error("VELA_RUNTIME_CONFIG not defined and failed to create kubernetes InClusterConfig!")
			return nil, err
		}
	} else {
		// https://pkg.go.dev/k8s.io/client-go/tools/clientcmd?tab=doc#BuildConfigFromFlags
		config, err = clientcmd.BuildConfigFromFlags("", c.config.File)
		if err != nil {
			return nil, err
		}
	}

	// creates Kubernetes client from configuration
	//
	// https://pkg.go.dev/k8s.io/client-go/kubernetes?tab=doc#NewForConfig
	_kubernetes, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// set the Kubernetes client in the runtime client
	c.Kubernetes = _kubernetes

	// creates VelaKubernetes client from configuration
	_velaKubernetes, err := velaK8sClient.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// set the VelaKubernetes client in the runtime client
	c.VelaKubernetes = _velaKubernetes

	return c, nil
}

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

	// The test is responsible for calling c.PodTracker.Start(ctx) if needed.
	// In some cases it is more convenient to call c.(MockKubernetesRuntime).StartPodTracker(ctx)

	return c, nil
}

// MockKubernetesRuntime makes it possible to use the client mocks in other packages.
type MockKubernetesRuntime interface {
	SetupMock() error
	StartPodTracker(context.Context)
	WaitForPodTrackerReady()
	WaitForPodCreate(string, string)
	SimulateResync()
	SimulateStatusUpdate(*v1.Pod, []v1.ContainerStatus) error
}

// SetupMock allows the Kubernetes runtime to perform additional Mock-related config.
// Many tests should call this right after they call runtime.SetupBuild (or executor.CreateBuild).
func (c *client) SetupMock() error {
	// This assumes that c.Pod.ObjectMeta.Namespace and c.Pod.ObjectMeta.Name are filled in.
	return c.PodTracker.setupMockFor(c.Pod)
}

// StartPodTracker tells the podTracker it can start populating the cache.
// This is only here for tests.
func (c *client) StartPodTracker(ctx context.Context) {
	c.PodTracker.Start(ctx)
}

// WaitForPodTrackerReady waits for PodTracker.Ready to be closed (which happens in AssembleBuild).
// This is only here for tests.
func (c *client) WaitForPodTrackerReady() {
	<-c.PodTracker.Ready
}

// WaitForPodCreate waits for PodTracker.Ready to be closed (which happens in AssembleBuild).
// This is only here for tests.
func (c *client) WaitForPodCreate(namespace, name string) {
	created := make(chan struct{})

	c.PodTracker.podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
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

	<-created
}

// SimulateResync simulates an resync where the PodTracker refreshes its cache.
// This is only here for tests.
func (c *client) SimulateResync() {
	// Future: maybe allow passing in either new or old pod
	oldPod := c.Pod.DeepCopy()
	oldPod.SetResourceVersion("older")

	// simulate a re-sync/PodUpdate event
	c.PodTracker.HandlePodUpdate(oldPod, c.Pod)
}

// SimulateUpdate simulates an update event from the k8s API.
// This is only here for tests.
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
