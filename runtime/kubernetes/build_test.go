// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"context"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-vela/server/compiler/types/pipeline"
	velav1alpha1 "github.com/go-vela/worker/runtime/kubernetes/apis/vela/v1alpha1"
)

func TestKubernetes_InspectBuild(t *testing.T) {
	// setup types
	_engine, err := NewMock(_pod)
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	// setup tests
	tests := []struct {
		name     string
		failure  bool
		pipeline *pipeline.Build
	}{
		{
			name:     "stages",
			failure:  false,
			pipeline: _stages,
		},
		{
			name:     "steps",
			failure:  false,
			pipeline: _steps,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err = _engine.InspectBuild(context.Background(), test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("InspectBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("InspectBuild returned err: %v", err)
			}
		})
	}
}

func TestKubernetes_SetupBuild(t *testing.T) {
	// needed to be able to make a pointers:
	trueBool := true
	twoString := "2"

	// testdata/pipeline-pods-template.yaml
	wantFromTemplateMetadata := velav1alpha1.PipelinePodTemplateMeta{
		Annotations: map[string]string{"annotation/foo": "bar"},
		Labels: map[string]string{
			"foo":      "bar",
			"pipeline": _steps.ID,
		},
	}

	// testdata/pipeline-pods-template-security-context.yaml
	wantFromTemplateSecurityContext := velav1alpha1.PipelinePodSecurityContext{
		RunAsNonRoot: &trueBool,
		Sysctls: []v1.Sysctl{
			{Name: "kernel.shm_rmid_forced", Value: "0"},
			{Name: "net.core.somaxconn", Value: "1024"},
			{Name: "kernel.msgmax", Value: "65536"},
		},
	}

	// testdata/pipeline-pods-template-node-selection.yaml
	wantFromTemplateNodeSelection := velav1alpha1.PipelinePodTemplateSpec{
		NodeSelector: map[string]string{"disktype": "ssd"},
		Affinity: &v1.Affinity{
			NodeAffinity: &v1.NodeAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: &v1.NodeSelector{
					NodeSelectorTerms: []v1.NodeSelectorTerm{
						{MatchExpressions: []v1.NodeSelectorRequirement{
							{Key: "kubernetes.io/os", Operator: v1.NodeSelectorOpIn, Values: []string{"linux"}},
						}},
					},
				},
				PreferredDuringSchedulingIgnoredDuringExecution: []v1.PreferredSchedulingTerm{
					{Weight: 1, Preference: v1.NodeSelectorTerm{
						MatchExpressions: []v1.NodeSelectorRequirement{
							{Key: "another-node-label-key", Operator: v1.NodeSelectorOpIn, Values: []string{"another-node-label-value"}},
						},
					}},
				},
			},
			PodAffinity: &v1.PodAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []v1.PodAffinityTerm{
					{
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{Key: "security", Operator: metav1.LabelSelectorOpIn, Values: []string{"S1"}},
							},
						},
						TopologyKey: "topology.kubernetes.io/zone",
					},
				},
			},
			PodAntiAffinity: &v1.PodAntiAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
					{
						Weight: 100,
						PodAffinityTerm: v1.PodAffinityTerm{
							LabelSelector: &metav1.LabelSelector{
								MatchExpressions: []metav1.LabelSelectorRequirement{
									{Key: "security", Operator: metav1.LabelSelectorOpIn, Values: []string{"S2"}},
								},
							},
							TopologyKey: "topology.kubernetes.io/zone",
						},
					},
				},
			},
		},
		Tolerations: []v1.Toleration{
			{
				Key:      "key1",
				Operator: v1.TolerationOpEqual,
				Value:    "value1",
				Effect:   v1.TaintEffectNoSchedule,
			},
			{
				Key:      "key1",
				Operator: v1.TolerationOpEqual,
				Value:    "value1",
				Effect:   v1.TaintEffectNoExecute,
			},
		},
	}

	// testdata/pipeline-pods-template-dns.yaml
	wantFromTemplateDNS := velav1alpha1.PipelinePodTemplateSpec{
		DNSPolicy: v1.DNSNone,
		DNSConfig: &v1.PodDNSConfig{
			Nameservers: []string{"1.2.3.4"},
			Searches: []string{
				"ns1.svc.cluster-domain.example",
				"my.dns.search.suffix",
			},
			Options: []v1.PodDNSConfigOption{
				{Name: "ndots", Value: &twoString},
				{Name: "edns0"},
			},
		},
	}

	// setup tests
	tests := []struct {
		name             string
		failure          bool
		pipeline         *pipeline.Build
		opts             []ClientOpt
		wantFromTemplate any
	}{
		{
			name:             "stages",
			failure:          false,
			pipeline:         _stages,
			opts:             nil,
			wantFromTemplate: nil,
		},
		{
			name:             "steps",
			failure:          false,
			pipeline:         _steps,
			opts:             nil,
			wantFromTemplate: nil,
		},
		{
			name:             "stages-PipelinePodsTemplate-empty",
			failure:          false,
			pipeline:         _stages,
			opts:             []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-empty.yaml")},
			wantFromTemplate: nil,
		},
		{
			name:             "steps-PipelinePodsTemplate-empty",
			failure:          false,
			pipeline:         _steps,
			opts:             []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-empty.yaml")},
			wantFromTemplate: nil,
		},
		{
			name:             "stages-PipelinePodsTemplate-metadata",
			failure:          false,
			pipeline:         _stages,
			opts:             []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template.yaml")},
			wantFromTemplate: wantFromTemplateMetadata,
		},
		{
			name:             "steps-PipelinePodsTemplate-metadata",
			failure:          false,
			pipeline:         _steps,
			opts:             []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template.yaml")},
			wantFromTemplate: wantFromTemplateMetadata,
		},
		{
			name:             "stages-PipelinePodsTemplate-SecurityContext",
			failure:          false,
			pipeline:         _stages,
			opts:             []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-security-context.yaml")},
			wantFromTemplate: wantFromTemplateSecurityContext,
		},
		{
			name:             "steps-PipelinePodsTemplate-SecurityContext",
			failure:          false,
			pipeline:         _steps,
			opts:             []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-security-context.yaml")},
			wantFromTemplate: wantFromTemplateSecurityContext,
		},
		{
			name:             "stages-PipelinePodsTemplate-NodeSelection",
			failure:          false,
			pipeline:         _stages,
			opts:             []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-node-selection.yaml")},
			wantFromTemplate: wantFromTemplateNodeSelection,
		},
		{
			name:             "steps-PipelinePodsTemplate-NodeSelection",
			failure:          false,
			pipeline:         _steps,
			opts:             []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-node-selection.yaml")},
			wantFromTemplate: wantFromTemplateNodeSelection,
		},
		{
			name:             "stages-PipelinePodsTemplate-dns",
			failure:          false,
			pipeline:         _stages,
			opts:             []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-dns.yaml")},
			wantFromTemplate: wantFromTemplateDNS,
		},
		{
			name:             "steps-PipelinePodsTemplate-dns",
			failure:          false,
			pipeline:         _steps,
			opts:             []ClientOpt{WithPodsTemplate("", "testdata/pipeline-pods-template-dns.yaml")},
			wantFromTemplate: wantFromTemplateDNS,
		},
		{
			name:             "stages-named PipelinePodsTemplate present",
			failure:          false,
			pipeline:         _stages,
			opts:             []ClientOpt{WithPodsTemplate("mock-pipeline-pods-template", "")},
			wantFromTemplate: nil,
		},
		{
			name:             "steps-named PipelinePodsTemplate present",
			failure:          false,
			pipeline:         _steps,
			opts:             []ClientOpt{WithPodsTemplate("mock-pipeline-pods-template", "")},
			wantFromTemplate: nil,
		},
		{
			name:             "stages-named PipelinePodsTemplate missing",
			failure:          true,
			pipeline:         _stages,
			opts:             []ClientOpt{WithPodsTemplate("missing-pipeline-pods-template", "")},
			wantFromTemplate: nil,
		},
		{
			name:             "steps-named PipelinePodsTemplate missing",
			failure:          true,
			pipeline:         _steps,
			opts:             []ClientOpt{WithPodsTemplate("missing-pipeline-pods-template", "")},
			wantFromTemplate: nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup types
			_engine, err := NewMock(&v1.Pod{}, test.opts...)
			if err != nil {
				t.Errorf("unable to create runtime engine: %v", err)
			}

			err = _engine.SetupBuild(context.Background(), test.pipeline)

			// this does not test the resulting pod spec (ie no tests for ObjectMeta, RestartPolicy)

			if test.failure {
				if err == nil {
					t.Errorf("SetupBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("SetupBuild returned err: %v", err)
			}

			// make sure that worker-defined labels are set and cannot be overridden by PipelinePodsTemplate
			if pipelineLabel, ok := _engine.Pod.Labels["pipeline"]; !ok {
				t.Errorf("Pod is missing the pipeline label: %v", _engine.Pod.ObjectMeta)
			} else if pipelineLabel != test.pipeline.ID {
				t.Errorf("Pod's pipeline label is %v, want %v", pipelineLabel, test.pipeline.ID)
			}

			switch test.wantFromTemplate.(type) {
			case velav1alpha1.PipelinePodTemplateMeta:
				want := test.wantFromTemplate.(velav1alpha1.PipelinePodTemplateMeta)

				// PipelinePodsTemplate defined Annotations
				if want.Annotations != nil && !reflect.DeepEqual(_engine.Pod.Annotations, want.Annotations) {
					t.Errorf("Pod.Annotations is %v, want %v", _engine.Pod.Annotations, want.Annotations)
				}

				// PipelinePodsTemplate defined Labels
				if want.Labels != nil && !reflect.DeepEqual(_engine.Pod.Labels, want.Labels) {
					t.Errorf("Pod.Labels is %v, want %v", _engine.Pod.Labels, want.Labels)
				}
			case velav1alpha1.PipelinePodSecurityContext:
				want := test.wantFromTemplate.(velav1alpha1.PipelinePodSecurityContext)

				// PipelinePodsTemplate defined SecurityContext.RunAsNonRoot
				if !reflect.DeepEqual(_engine.Pod.Spec.SecurityContext.RunAsNonRoot, want.RunAsNonRoot) {
					t.Errorf("Pod.SecurityContext.RunAsNonRoot is %v, want %v", _engine.Pod.Spec.SecurityContext.RunAsNonRoot, want.RunAsNonRoot)
				}

				// PipelinePodsTemplate defined SecurityContext.Sysctls
				if want.Sysctls != nil && !reflect.DeepEqual(_engine.Pod.Spec.SecurityContext.Sysctls, want.Sysctls) {
					t.Errorf("Pod.SecurityContext.Sysctls is %v, want %v", _engine.Pod.Spec.SecurityContext.Sysctls, want.Sysctls)
				}
			case velav1alpha1.PipelinePodTemplateSpec:
				want := test.wantFromTemplate.(velav1alpha1.PipelinePodTemplateSpec)

				// PipelinePodsTemplate defined NodeSelector
				if want.NodeSelector != nil && !reflect.DeepEqual(_engine.Pod.Spec.NodeSelector, want.NodeSelector) {
					t.Errorf("Pod.NodeSelector is %v, want %v", _engine.Pod.Spec.NodeSelector, want.NodeSelector)
				}

				// PipelinePodsTemplate defined Affinity
				if want.Affinity != nil && !reflect.DeepEqual(_engine.Pod.Spec.Affinity, want.Affinity) {
					t.Errorf("Pod.Affinity is %v, want %v", _engine.Pod.Spec.Affinity, want.Affinity)
				}

				// PipelinePodsTemplate defined Tolerations
				if want.Tolerations != nil && !reflect.DeepEqual(_engine.Pod.Spec.Tolerations, want.Tolerations) {
					t.Errorf("Pod.Tolerations is %v, want %v", _engine.Pod.Spec.Tolerations, want.Tolerations)
				}

				// PipelinePodsTemplate defined DNSPolicy
				if len(want.DNSPolicy) > 0 && _engine.Pod.Spec.DNSPolicy != want.DNSPolicy {
					t.Errorf("Pod.DNSPolicy is %v, want %v", _engine.Pod.Spec.DNSPolicy, want.DNSPolicy)
				}

				// PipelinePodsTemplate defined DNSConfig
				if want.DNSConfig != nil && !reflect.DeepEqual(_engine.Pod.Spec.DNSConfig, want.DNSConfig) {
					t.Errorf("Pod.DNSConfig is %v, want %v", _engine.Pod.Spec.DNSConfig, want.DNSConfig)
				}
			}
		})
	}
}

func TestKubernetes_StreamBuild(t *testing.T) {
	tests := []struct {
		name     string
		failure  bool
		doCancel bool
		doReady  bool
		pipeline *pipeline.Build
		pod      *v1.Pod
	}{
		{
			name:     "stages canceled",
			failure:  false,
			doCancel: true,
			pipeline: _stages,
			pod:      _stagesPod,
		},
		{
			name:     "steps canceled",
			failure:  false,
			doCancel: true,
			pipeline: _steps,
			pod:      _pod,
		},
		{
			name:     "stages ready",
			failure:  false,
			doReady:  true,
			pipeline: _stages,
			pod:      _stagesPod,
		},
		{
			name:     "steps ready",
			failure:  false,
			doReady:  true,
			pipeline: _steps,
			pod:      _pod,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := NewMock(test.pod)
			if err != nil {
				t.Errorf("unable to create runtime engine: %v", err)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// StreamBuild and AssembleBuild coordinate their work.
			go func() {
				if test.doCancel {
					// simulate canceled build
					cancel()
				} else if test.doReady {
					// simulate AssembleBuild
					_engine.MarkPodTrackerReady()
				}
			}()

			err = _engine.StreamBuild(ctx, test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("StreamBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("StreamBuild returned err: %v", err)
			}
		})
	}
}

func TestKubernetes_AssembleBuild(t *testing.T) {
	// setup tests
	tests := []struct {
		name     string
		failure  bool
		pipeline *pipeline.Build
		// k8sPod is the pod that the mock Kubernetes client will return
		k8sPod *v1.Pod
		// enginePod is the pod under construction in the Runtime Engine
		enginePod *v1.Pod
	}{
		{
			name:      "stages",
			failure:   false,
			pipeline:  _stages,
			k8sPod:    &v1.Pod{},
			enginePod: _stagesPod,
		},
		{
			name:      "steps",
			failure:   false,
			pipeline:  _steps,
			k8sPod:    &v1.Pod{},
			enginePod: _pod,
		},
		{
			name:      "stages-pod already exists",
			failure:   true,
			pipeline:  _stages,
			k8sPod:    _stagesPod,
			enginePod: _stagesPod,
		},
		{
			name:      "steps-pod already exists",
			failure:   true,
			pipeline:  _steps,
			k8sPod:    _pod,
			enginePod: _pod,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := NewMock(test.k8sPod)
			if err != nil {
				t.Errorf("unable to create runtime engine: %v", err)
			}

			_engine.Pod = test.enginePod

			_engine.containersLookup = map[string]int{}
			for i, ctn := range test.enginePod.Spec.Containers {
				_engine.containersLookup[ctn.Name] = i
			}

			// StreamBuild and AssembleBuild coordinate their work, so, emulate
			// executor.StreamBuild which calls runtime.StreamBuild concurrently.
			go func() {
				err := _engine.StreamBuild(context.Background(), test.pipeline)
				if err != nil {
					t.Errorf("unable to start PodTracker via StreamBuild")
				}
			}()

			err = _engine.AssembleBuild(context.Background(), test.pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("AssembleBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("AssembleBuild returned err: %v", err)
			}
		})
	}
}

func TestKubernetes_RemoveBuild(t *testing.T) {
	// setup tests
	tests := []struct {
		name       string
		failure    bool
		createdPod bool
		pipeline   *pipeline.Build
		pod        *v1.Pod
	}{
		{
			name:       "stages-createdPod-pod in k8s",
			failure:    false,
			createdPod: true,
			pipeline:   _stages,
			pod:        _pod,
		},
		{
			name:       "steps-createdPod-pod in k8s",
			failure:    false,
			createdPod: true,
			pipeline:   _steps,
			pod:        _pod,
		},
		{
			name:       "stages-not createdPod-pod not in k8s",
			failure:    false,
			createdPod: false,
			pipeline:   _stages,
			pod:        &v1.Pod{},
		},
		{
			name:       "steps-not createdPod-pod not in k8s",
			failure:    false,
			pipeline:   _steps,
			pod:        &v1.Pod{},
			createdPod: false,
		},
		{
			name:       "stages-createdPod-pod not in k8s",
			failure:    true,
			pipeline:   _stages,
			pod:        &v1.Pod{},
			createdPod: true,
		},
		{
			name:       "steps-createdPod-pod not in k8s",
			failure:    true,
			pipeline:   _steps,
			pod:        &v1.Pod{},
			createdPod: true,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_engine, err := NewMock(test.pod)
			if err != nil {
				t.Errorf("unable to create runtime engine: %v", err)
			}

			_engine.createdPod = test.createdPod

			err = _engine.RemoveBuild(context.Background(), test.pipeline)
			if test.failure {
				if err == nil {
					t.Errorf("RemoveBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("RemoveBuild returned err: %v", err)
			}
		})
	}
}
