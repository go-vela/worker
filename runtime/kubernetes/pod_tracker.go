// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	kubeinformers "k8s.io/client-go/informers"
	informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

// podTracker contains Informers used to watch and synchronize local k8s caches
// This is similar to a typical Kubernetes controller (eg like k8s.io/sample-controller.Controller)
type podTracker struct {
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	Logger *logrus.Entry
	// TrackedPod is the Namespace/Name of the tracked pod
	TrackedPod string

	// informerFactory is used to create Informers and Listers
	informerFactory kubeinformers.SharedInformerFactory
	// podInformer watches the given pod, caches the results, and makes them available in podLister
	podInformer informers.PodInformer

	// PodLister helps list Pods. All objects returned here must be treated as read-only.
	PodLister listers.PodLister
	// PodSynced is a function that can be used to determine if an informer has synced.
	// This is useful for determining if caches have synced.
	PodSynced cache.InformerSynced
}

// AddPodInformerEventHandler adds an event handler to the cache.SharedInformer for the Pod.
// Events to a single handler are delivered sequentially, but there is no coordination
// between different handlers.
// Make sure to add the ResourceEventHandler with this before running Start.
func (p *podTracker) AddPodInformerEventHandler(handler cache.ResourceEventHandler) {
	p.podInformer.Informer().AddEventHandler(handler)
}

// Start kicks off the API calls to start populating the cache.
// There is no need to run this in a separate goroutine (ie go podTracker.Start(ctx)).
func (p *podTracker) Start(ctx context.Context) {
	p.Logger.Tracef("starting PodTracker for pod %s", p.TrackedPod)

	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	p.informerFactory.Start(ctx.Done())
}

// newPodTracker initializes a podTracker with a given clientset for a given pod.
func newPodTracker(log *logrus.Entry, clientset kubernetes.Interface, pod *v1.Pod, defaultResync time.Duration) (*podTracker, error) {
	trackedPod := pod.ObjectMeta.Namespace + "/" + pod.ObjectMeta.Name
	if pod.ObjectMeta.Name == "" || pod.ObjectMeta.Namespace == "" {
		return nil, fmt.Errorf("newPodTracker expects pod to have Name and Namespace, got %s", trackedPod)
	}

	log.Tracef("creating PodTracker for pod %s", trackedPod)

	// create label selector for watching the pod
	selector, err := labels.NewRequirement(
		"pipeline",
		selection.Equals,
		[]string{fields.EscapeValue(pod.ObjectMeta.Name)},
	)
	if err != nil {
		return nil, err
	}

	// create filtered Informer factory which is commonly used for k8s controllers
	informerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(
		clientset,
		defaultResync,
		kubeinformers.WithNamespace(pod.ObjectMeta.Namespace),
		kubeinformers.WithTweakListOptions(func(listOptions *metav1.ListOptions) {
			listOptions.LabelSelector = selector.String()
		}),
	)
	podInformer := informerFactory.Core().V1().Pods()

	// initialize podTracker
	tracker := podTracker{
		Logger:          log,
		TrackedPod:      trackedPod,
		informerFactory: informerFactory,
		podInformer:     podInformer,
		PodLister:       podInformer.Lister(),
		PodSynced:       podInformer.Informer().HasSynced,
	}

	return &tracker, nil
}

// mockPodTracker returns a new podTracker with the given pod pre-loaded in the cache.
func mockPodTracker(log *logrus.Entry, clientset kubernetes.Interface, pod *v1.Pod) (*podTracker, error) {
	// Make sure test pods are valid before passing to PodTracker (ie support &v1.Pod{}).
	if pod.ObjectMeta.Name == "" {
		pod.ObjectMeta.Name = "test-pod"
	}

	if pod.ObjectMeta.Namespace == "" {
		pod.ObjectMeta.Namespace = "test"
	}

	tracker, err := newPodTracker(log, clientset, pod, 0*time.Second)
	if err != nil {
		return nil, err
	}

	// pre-populate the podInformer cache
	err = tracker.podInformer.Informer().GetIndexer().Add(pod)
	if err != nil {
		return nil, err
	}

	return tracker, err
}
