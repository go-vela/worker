// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"fmt"

	"github.com/sirupsen/logrus"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/testing"
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
	// maxLogSize is the max log size enforced by the executor
	maxLogSize uint
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
// nolint: revive // ignore returning unexported client
func New(opts ...ClientOpt) (*client, error) {
	// create new Kubernetes client
	c := new(client)

	// create new fields
	c.config = new(config)
	c.Pod = new(v1.Pod)

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
// nolint: revive // ignore returning unexported client
func NewMock(_pod *v1.Pod, opts ...ClientOpt) (*client, error) {
	// create new Kubernetes client
	c := new(client)

	// create new fields
	c.config = new(config)
	c.Pod = new(v1.Pod)

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

	// make the Kubernetes fake client in the runtime client
	//
	// https://pkg.go.dev/k8s.io/client-go/kubernetes/fake?tab=doc#NewSimpleClientset
	fakeClientset := fake.NewSimpleClientset(c.Pod)

	// work around bug in default ObjectReactor: github.com/kubernetes/client-go/issues/873
	fakeClientset.PrependReactor("get", "pods/log",
		func(action testing.Action) (handled bool, ret runtime.Object, err error) {
			// handled=true to avoid calling the default * reactor which is buggy, and
			// ret=nil as it is unused in k8s.io/client-go/kubernetes/typed/v1/fake.*FakePods.GetLogs
			// where it is returned from c.Fake.Invokes() .
			return true, nil, fmt.Errorf("no reaction implemented for verb:get resource:pods/log")
		},
	)

	// set the Kubernetes fake client in the runtime client
	c.Kubernetes = fakeClientset

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
