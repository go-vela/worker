// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

func TestNewPodTracker(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())
	clientset := fake.NewSimpleClientset()

	tests := []struct {
		name    string
		pod     *v1.Pod
		wantErr bool
	}{
		{
			name:    "pass-with-pod",
			pod:     _pod,
			wantErr: false,
		},
		{
			name:    "error-with-nil-pod",
			pod:     nil,
			wantErr: true,
		},
		{
			name:    "error-with-empty-pod",
			pod:     &v1.Pod{},
			wantErr: true,
		},
		{
			name: "error-with-pod-without-namespace",
			pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
			},
			wantErr: true,
		},
		{
			name: "fail-with-pod",
			pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "github-octocat-1-for-some-odd-reason-this-name-is-way-too-long-and-will-cause-an-error",
					Namespace: _pod.ObjectMeta.Namespace,
					Labels:    _pod.ObjectMeta.Labels,
				},
				TypeMeta: _pod.TypeMeta,
				Spec:     _pod.Spec,
				Status:   _pod.Status,
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := newPodTracker(logger, clientset, test.pod, 0*time.Second)
			if (err != nil) != test.wantErr {
				t.Errorf("newPodTracker() error = %v, wantErr %v", err, test.wantErr)
				return
			}
		})
	}
}

func Test_podTracker_getTrackedPod(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	tests := []struct {
		name       string
		trackedPod string // namespace/podName
		obj        interface{}
		want       *v1.Pod
	}{
		{
			name:       "got-tracked-pod",
			trackedPod: "test/github-octocat-1",
			obj:        _pod,
			want:       _pod,
		},
		{
			name:       "wrong-pod",
			trackedPod: "test/github-octocat-2",
			obj:        _pod,
			want:       nil,
		},
		{
			name:       "invalid-type",
			trackedPod: "test/github-octocat-1",
			obj:        new(v1.PodTemplate),
			want:       nil,
		},
		{
			name:       "nil",
			trackedPod: "test/nil",
			obj:        nil,
			want:       nil,
		},
		{
			name:       "tombstone-pod",
			trackedPod: "test/github-octocat-1",
			obj: cache.DeletedFinalStateUnknown{
				Key: "test/github-octocat-1",
				Obj: _pod,
			},
			want: _pod,
		},
		{
			name:       "tombstone-nil",
			trackedPod: "test/github-octocat-1",
			obj: cache.DeletedFinalStateUnknown{
				Key: "test/github-octocat-1",
				Obj: nil,
			},
			want: nil,
		},
		{
			name:       "tombstone-invalid-type",
			trackedPod: "test/github-octocat-1",
			obj: cache.DeletedFinalStateUnknown{
				Key: "test/github-octocat-1",
				Obj: new(v1.PodTemplate),
			},
			want: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := podTracker{
				Logger:     logger,
				TrackedPod: test.trackedPod,
				// other fields not used by getTrackedPod
				// if they're needed, use newPodTracker
			}
			if got := p.getTrackedPod(test.obj); !reflect.DeepEqual(got, test.want) {
				t.Errorf("getTrackedPod() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_podTracker_HandlePodAdd(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	tests := []struct {
		name       string
		trackedPod string // namespace/podName
		obj        interface{}
	}{
		{
			name:       "got-tracked-pod",
			trackedPod: "test/github-octocat-1",
			obj:        _pod,
		},
		{
			name:       "wrong-pod",
			trackedPod: "test/github-octocat-2",
			obj:        _pod,
		},
		{
			name:       "invalid-type",
			trackedPod: "test/github-octocat-1",
			obj:        new(v1.PodTemplate),
		},
		{
			name:       "nil",
			trackedPod: "test/nil",
			obj:        nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &podTracker{
				Logger:     logger,
				TrackedPod: test.trackedPod,
				// other fields not used by getTrackedPod
				// if they're needed, use newPodTracker
			}

			// just make sure this doesn't panic
			p.HandlePodAdd(test.obj)
		})
	}
}

func Test_podTracker_HandlePodUpdate(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	tests := []struct {
		name       string
		trackedPod string // namespace/podName
		oldObj     interface{}
		newObj     interface{}
	}{
		{
			name:       "re-sync event without change",
			trackedPod: "test/github-octocat-1",
			oldObj:     _pod,
			newObj:     _pod,
		},
		{
			name:       "wrong-pod",
			trackedPod: "test/github-octocat-2",
			oldObj:     _pod,
			newObj:     _pod,
		},
		{
			name:       "invalid-type-old",
			trackedPod: "test/github-octocat-1",
			oldObj:     new(v1.PodTemplate),
			newObj:     _pod,
		},
		{
			name:       "nil-old",
			trackedPod: "test/github-octocat-1",
			oldObj:     nil,
			newObj:     _pod,
		},
		{
			name:       "invalid-type-new",
			trackedPod: "test/github-octocat-1",
			oldObj:     _pod,
			newObj:     new(v1.PodTemplate),
		},
		{
			name:       "nil-new",
			trackedPod: "test/github-octocat-1",
			oldObj:     _pod,
			newObj:     nil,
		},
		{
			name:       "invalid-type-both",
			trackedPod: "test/github-octocat-1",
			oldObj:     new(v1.PodTemplate),
			newObj:     new(v1.PodTemplate),
		},
		{
			name:       "nil-both",
			trackedPod: "test/nil",
			oldObj:     nil,
			newObj:     nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &podTracker{
				Logger:     logger,
				TrackedPod: test.trackedPod,
				// other fields not used by getTrackedPod
				// if they're needed, use newPodTracker
			}

			// just make sure this doesn't panic
			p.HandlePodUpdate(test.oldObj, test.newObj)
		})
	}
}

func Test_podTracker_HandlePodDelete(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	tests := []struct {
		name       string
		trackedPod string // namespace/podName
		obj        interface{}
	}{
		{
			name:       "got-tracked-pod",
			trackedPod: "test/github-octocat-1",
			obj:        _pod,
		},
		{
			name:       "wrong-pod",
			trackedPod: "test/github-octocat-2",
			obj:        _pod,
		},
		{
			name:       "invalid-type",
			trackedPod: "test/github-octocat-1",
			obj:        new(v1.PodTemplate),
		},
		{
			name:       "nil",
			trackedPod: "test/nil",
			obj:        nil,
		},
		{
			name:       "tombstone-pod",
			trackedPod: "test/github-octocat-1",
			obj: cache.DeletedFinalStateUnknown{
				Key: "test/github-octocat-1",
				Obj: _pod,
			},
		},
		{
			name:       "tombstone-nil",
			trackedPod: "test/github-octocat-1",
			obj: cache.DeletedFinalStateUnknown{
				Key: "test/github-octocat-1",
				Obj: nil,
			},
		},
		{
			name:       "tombstone-invalid-type",
			trackedPod: "test/github-octocat-1",
			obj: cache.DeletedFinalStateUnknown{
				Key: "test/github-octocat-1",
				Obj: new(v1.PodTemplate),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &podTracker{
				Logger:     logger,
				TrackedPod: test.trackedPod,
				// other fields not used by getTrackedPod
				// if they're needed, use newPodTracker
			}

			// just make sure this doesn't panic
			p.HandlePodDelete(test.obj)
		})
	}
}

func Test_podTracker_Stop(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())
	clientset := fake.NewSimpleClientset()

	tests := []struct {
		name    string
		pod     *v1.Pod
		started bool
	}{
		{
			name:    "started",
			pod:     _pod,
			started: true,
		},
		{
			name:    "not started",
			pod:     _pod,
			started: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tracker, err := newPodTracker(logger, clientset, test.pod, 0*time.Second)
			if err != nil {
				t.Errorf("newPodTracker() error = %v", err)
				return
			}

			if test.started {
				tracker.Start(context.Background())
			}
			tracker.Stop()
		})
	}
}

func Test_podTracker_getTrackedPodEvent(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	_podEvent := mockContainerEvent(
		_pod,
		"step-github-ooctocat-1-echo",
		reasonFailed,
		"foobar",
	)

	_podTemplateEvent := _podEvent.DeepCopy()
	_podTemplateEvent.InvolvedObject.Kind = new(v1.PodTemplate).TypeMeta.Kind

	tests := []struct {
		name       string
		trackedPod string // namespace/podName
		obj        interface{}
		want       *v1.Event
	}{
		{
			name:       "got-tracked-pod-event",
			trackedPod: "test/github-octocat-1",
			obj:        _podEvent,
			want:       _podEvent,
		},
		{
			name:       "wrong-pod",
			trackedPod: "test/github-octocat-2",
			obj:        _podEvent,
			want:       nil,
		},
		{
			name:       "invalid-type",
			trackedPod: "test/github-octocat-1",
			obj:        new(v1.PodTemplate),
			want:       nil,
		},
		{
			name:       "nil",
			trackedPod: "test/nil",
			obj:        nil,
			want:       nil,
		},
		{
			name:       "invalid-involved-object-type",
			trackedPod: "test/github-octocat-1",
			obj:        _podTemplateEvent,
			want:       nil,
		},
		{
			name:       "tombstone-pod-event",
			trackedPod: "test/github-octocat-1",
			obj: cache.DeletedFinalStateUnknown{
				Key: "test/github-octocat-1",
				Obj: _podEvent,
			},
			want: nil,
		},
		{
			name:       "tombstone-nil",
			trackedPod: "test/github-octocat-1",
			obj: cache.DeletedFinalStateUnknown{
				Key: "test/github-octocat-1",
				Obj: nil,
			},
			want: nil,
		},
		{
			name:       "tombstone-invalid-type",
			trackedPod: "test/github-octocat-1",
			obj: cache.DeletedFinalStateUnknown{
				Key: "test/github-octocat-1",
				Obj: new(v1.PodTemplate),
			},
			want: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := podTracker{
				Logger:     logger,
				TrackedPod: test.trackedPod,
				// other fields not used by getTrackedPodEvent
				// if they're needed, use newPodTracker
			}
			if got := p.getTrackedPodEvent(test.obj); !reflect.DeepEqual(got, test.want) {
				t.Errorf("getTrackedPodEvent() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_podTracker_HandleEventAdd(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	_podEvent := mockContainerEvent(
		_pod,
		"step-github-ooctocat-1-echo",
		reasonInspectFailed,
		"foobar",
	)

	_podTemplateEvent := _podEvent.DeepCopy()
	_podTemplateEvent.InvolvedObject.Kind = new(v1.PodTemplate).TypeMeta.Kind

	tests := []struct {
		name       string
		trackedPod string // namespace/podName
		obj        interface{}
	}{
		{
			name:       "got-tracked-pod-event",
			trackedPod: "test/github-octocat-1",
			obj:        _podEvent,
		},
		{
			name:       "wrong-pod",
			trackedPod: "test/github-octocat-2",
			obj:        _podEvent,
		},
		{
			name:       "invalid-type",
			trackedPod: "test/github-octocat-1",
			obj:        new(v1.PodTemplate),
		},
		{
			name:       "nil",
			trackedPod: "test/nil",
			obj:        nil,
		},
		{
			name:       "invalid-involved-object-type",
			trackedPod: "test/github-octocat-1",
			obj:        _podTemplateEvent,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &podTracker{
				Logger:     logger,
				TrackedPod: test.trackedPod,
				// other fields not used by getTrackedPodEvent
				// if they're needed, use newPodTracker
			}

			// just make sure this doesn't panic
			p.HandleEventAdd(test.obj)
		})
	}
}

func Test_podTracker_HandleEventUpdate(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	_podEvent := mockContainerEvent(
		_pod,
		"step-github-ooctocat-1-echo",
		reasonBackOff,
		"foobar",
	)

	_podTemplateEvent := _podEvent.DeepCopy()
	_podTemplateEvent.InvolvedObject.Kind = new(v1.PodTemplate).TypeMeta.Kind

	tests := []struct {
		name       string
		trackedPod string // namespace/podName
		oldObj     interface{}
		newObj     interface{}
	}{
		{
			name:       "re-sync event without change",
			trackedPod: "test/github-octocat-1",
			oldObj:     _podEvent,
			newObj:     _podEvent,
		},
		{
			name:       "wrong-pod",
			trackedPod: "test/github-octocat-2",
			oldObj:     _podEvent,
			newObj:     _podEvent,
		},
		{
			name:       "invalid-type-old",
			trackedPod: "test/github-octocat-1",
			oldObj:     new(v1.PodTemplate),
			newObj:     _podEvent,
		},
		{
			name:       "nil-old",
			trackedPod: "test/github-octocat-1",
			oldObj:     nil,
			newObj:     _podEvent,
		},
		{
			name:       "invalid-involved-object-type-old",
			trackedPod: "test/github-octocat-1",
			oldObj:     _podTemplateEvent,
			newObj:     _podEvent,
		},
		{
			name:       "invalid-type-new",
			trackedPod: "test/github-octocat-1",
			oldObj:     _podEvent,
			newObj:     new(v1.PodTemplate),
		},
		{
			name:       "nil-new",
			trackedPod: "test/github-octocat-1",
			oldObj:     _podEvent,
			newObj:     nil,
		},
		{
			name:       "invalid-involved-object-type-new",
			trackedPod: "test/github-octocat-1",
			oldObj:     _podEvent,
			newObj:     _podTemplateEvent,
		},
		{
			name:       "invalid-type-both",
			trackedPod: "test/github-octocat-1",
			oldObj:     new(v1.PodTemplate),
			newObj:     new(v1.PodTemplate),
		},
		{
			name:       "nil-both",
			trackedPod: "test/nil",
			oldObj:     nil,
			newObj:     nil,
		},
		{
			name:       "invalid-involved-object-type-both",
			trackedPod: "test/github-octocat-1",
			oldObj:     _podTemplateEvent,
			newObj:     _podTemplateEvent,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &podTracker{
				Logger:     logger,
				TrackedPod: test.trackedPod,
				// other fields not used by getTrackedPodEvent
				// if they're needed, use newPodTracker
			}

			// just make sure this doesn't panic
			p.HandleEventUpdate(test.oldObj, test.newObj)
		})
	}
}
