// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package kubernetes

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	// The k8s libraries have some quirks around yaml marshaling.
	// They use `json` instead of `yaml` to annotate their struct Tags.
	// So, we need to use "sigs.k8s.io/yaml" instead of "github.com/buildkite/yaml".
	"sigs.k8s.io/yaml"

	velav1alpha1 "github.com/go-vela/worker/runtime/kubernetes/apis/vela/v1alpha1"
)

// ClientOpt represents a configuration option to initialize the runtime client for Kubernetes.
type ClientOpt func(*client) error

// WithConfigFile sets the config file in the runtime client for Kubernetes.
func WithConfigFile(file string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring config file in kubernetes runtime client")

		// set the runtime config file in the kubernetes client
		c.config.File = file

		return nil
	}
}

// WithHostVolumes sets the host volumes in the runtime client for Kubernetes.
func WithHostVolumes(volumes []string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring host volumes in kubernetes runtime client")

		// set the runtime host volumes in the kubernetes client
		c.config.Volumes = volumes

		return nil
	}
}

// WithLogger sets the logger in the runtime client for Kubernetes.
func WithLogger(logger *logrus.Entry) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring logger in kubernetes runtime client")

		// check if the logger provided is empty
		if logger != nil {
			// set the runtime logger in the kubernetes client
			c.Logger = logger
		}

		return nil
	}
}

// WithNamespace sets the namespace in the runtime client for Kubernetes.
func WithNamespace(namespace string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring namespace in kubernetes runtime client")

		// check if the namespace provided is empty
		if len(namespace) == 0 {
			return fmt.Errorf("no Kubernetes namespace provided")
		}

		// set the runtime namespace in the kubernetes client
		c.config.Namespace = namespace

		return nil
	}
}

// WithPodsTemplate sets the PipelinePodsTemplateName or loads the PipelinePodsTemplate
// from file in the runtime client for Kubernetes.
func WithPodsTemplate(name string, path string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring pipeline pods template in kubernetes runtime client")

		// check if a PipelinePodsTemplate was requested
		if len(name) == 0 && len(path) == 0 {
			// no PipelinePodTemplate to load
			return nil
		}

		if len(name) == 0 {
			// load the PodsTemplate from the path (must restart Worker to reload the local file)
			if data, err := os.ReadFile(path); err == nil {
				pipelinePodsTemplate := velav1alpha1.PipelinePodsTemplate{}

				err := yaml.UnmarshalStrict(data, &pipelinePodsTemplate)
				if err != nil {
					return err
				}

				c.PipelinePodTemplate = &pipelinePodsTemplate.Spec.Template
			}

			return nil
		}

		// set the runtime namespace in the kubernetes client for just-in-time retrieval
		c.config.PipelinePodsTemplateName = name

		return nil
	}
}

// WithPrivilegedImages sets the privileged images in the runtime client for Kubernetes.
func WithPrivilegedImages(images []string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring privileged images in kubernetes runtime client")

		// set the runtime privileged images in the kubernetes client
		c.config.Images = images

		return nil
	}
}
