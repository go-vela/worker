// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	velav1alpha1 "github.com/go-vela/worker/runtime/kubernetes/apis/vela/v1alpha1"
	velaK8sClient "github.com/go-vela/worker/runtime/kubernetes/generated/clientset/versioned"
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
