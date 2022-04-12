// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"sync"
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

	logBuffer "github.com/go-vela/worker/internal/log"
)

// ErrTruncatedLogs is an error that allows the log streaming to indicate
// that it had to stop streaming logs due to maxLogSize.
var ErrTruncatedLogs = errors.New("TruncatedLogs")

// containerTracker contains useful signals that are managed by the podTracker.
type containerTracker struct {
	// Name is the name of the container
	Name string
	// terminatedOnce ensures that the Terminated channel only gets closed once.
	terminatedOnce sync.Once
	// Terminated will be closed once the container reaches a terminal state.
	Terminated chan struct{}
	// Logs contains all logs streamed for this container (they may be truncated if too large).
	Logs logBuffer.Buffer
	// LogsError holds an error if streaming the logs had an unrecoverable failure
	LogsError error
}

// podTracker contains Informers used to watch and synchronize local k8s caches.
// This is similar to a typical Kubernetes controller (eg like k8s.io/sample-controller.Controller).
type podTracker struct {
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	Logger *logrus.Entry
	// Namespace is the namespace the tracked pod is in
	Namespace string
	// Name is the name of the pod
	Name string
	// TrackedPod is the Namespace/Name of the tracked pod
	TrackedPod string

	// informerFactory is used to create Informers and Listers
	informerFactory kubeinformers.SharedInformerFactory
	// podInformer watches the given pod, caches the results, and makes them available in podLister
	podInformer informers.PodInformer

	// https://pkg.go.dev/k8s.io/client-go/kubernetes#Interface
	Kubernetes kubernetes.Interface
	// PodLister helps list Pods. All objects returned here must be treated as read-only.
	PodLister listers.PodLister
	// PodSynced is a function that can be used to determine if an informer has synced.
	// This is useful for determining if caches have synced.
	PodSynced cache.InformerSynced

	// Containers maps the container name to a containerTracker
	Containers map[string]*containerTracker
}

// HandlePodAdd is an AddFunc for cache.ResourceEventHandlerFuncs for Pods.
func (p podTracker) HandlePodAdd(newObj interface{}) {
	newPod := p.getTrackedPod(newObj)
	if newPod == nil {
		// not valid or not our tracked pod
		return
	}

	p.inspectContainerStatuses(newPod)
}

// HandlePodUpdate is an UpdateFunc for cache.ResourceEventHandlerFuncs for Pods.
func (p podTracker) HandlePodUpdate(oldObj, newObj interface{}) {
	oldPod := p.getTrackedPod(oldObj)
	newPod := p.getTrackedPod(newObj)

	if oldPod == nil || newPod == nil {
		// not valid or not our tracked pod
		return
	}
	// if we need to optimize and avoid the resync update events, we can do this:
	//if newPod.ResourceVersion == oldPod.ResourceVersion {
	//	// Periodic resync will send update events for all known Pods
	//	// If ResourceVersion is the same we have to look harder for Status changes
	//	if newPod.Status.Phase == oldPod.Status.Phase && newPod.Status.Size() == oldPod.Status.Size() {
	//		return
	//	}
	//}

	p.inspectContainerStatuses(newPod)
}

// HandlePodDelete is an DeleteFunc for cache.ResourceEventHandlerFuncs for Pods.
func (p podTracker) HandlePodDelete(oldObj interface{}) {
	oldPod := p.getTrackedPod(oldObj)
	if oldPod == nil {
		// not valid or not our tracked pod
		return
	}

	p.inspectContainerStatuses(oldPod)
}

// getTrackedPod tries to convert the obj into a Pod and makes sure it is the tracked Pod.
// This should only be used by the funcs of cache.ResourceEventHandlerFuncs.
func (p podTracker) getTrackedPod(obj interface{}) *v1.Pod {
	var (
		pod *v1.Pod
		ok  bool
	)

	if pod, ok = obj.(*v1.Pod); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			p.Logger.Errorf("error decoding pod, invalid type")
			return nil
		}

		pod, ok = tombstone.Obj.(*v1.Pod)
		if !ok {
			p.Logger.Errorf("error decoding pod tombstone, invalid type")
			return nil
		}
	}

	trackedPod := pod.GetNamespace() + "/" + pod.GetName()
	if trackedPod != p.TrackedPod {
		p.Logger.Errorf("error got unexpected pod: %s", trackedPod)
		return nil
	}

	return pod
}

// Start kicks off the API calls to start populating the cache.
// There is no need to run this in a separate goroutine (ie go podTracker.Start(stopCh)).
func (p podTracker) Start(ctx context.Context, maxLogSize uint) {
	p.Logger.Tracef("starting PodTracker for pod %s", p.TrackedPod)

	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	p.informerFactory.Start(ctx.Done())

	// Begin log streaming before any of the step containers have started
	// to ensure to capture logs for short-lived steps.
	for _, ctnTracker := range p.Containers {
		go p.streamContainerLogs(ctx, ctnTracker, maxLogSize)
	}
}

// TrackContainers creates a containerTracker for each container.
func (p *podTracker) TrackContainers(containers []v1.Container) {
	p.Logger.Tracef("adding %d containerTrackers for pod %s", len(containers), p.TrackedPod)
	if p.Containers == nil {
		p.Containers = map[string]*containerTracker{}
	}

	for _, ctn := range containers {
		p.Containers[ctn.Name] = &containerTracker{
			Name:       ctn.Name,
			Terminated: make(chan struct{}),
			Logs:       logBuffer.NewBuffer(),
		}
	}
}

// newPodTracker initializes a podTracker with a given clientset for a given pod.
func newPodTracker(log *logrus.Entry, clientset kubernetes.Interface, pod *v1.Pod, defaultResync time.Duration) (*podTracker, error) {
	if pod == nil {
		return nil, fmt.Errorf("newPodTracker expected a pod, got nil")
	}

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
		Namespace:       pod.ObjectMeta.Namespace,
		Name:            pod.ObjectMeta.Name,
		TrackedPod:      trackedPod,
		Kubernetes:      clientset,
		informerFactory: informerFactory,
		podInformer:     podInformer,
		PodLister:       podInformer.Lister(),
		PodSynced:       podInformer.Informer().HasSynced,
	}

	// register event handler funcs in podInformer
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    tracker.HandlePodAdd,
		UpdateFunc: tracker.HandlePodUpdate,
		DeleteFunc: tracker.HandlePodDelete,
	})

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

	// init containerTrackers as well
	tracker.TrackContainers(pod.Spec.Containers)

	// pre-populate the podInformer cache
	err = tracker.podInformer.Informer().GetIndexer().Add(pod)
	if err != nil {
		return nil, err
	}

	return tracker, err
}
