// SPDX-License-Identifier: Apache-2.0

package service

import (
	"fmt"
	"sync"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
)

// Load attempts to capture the library service
// representing the container from the map.
func Load(c *pipeline.Container, m *sync.Map) (*api.Service, error) {
	// check if the container provided is empty
	if c == nil {
		return nil, fmt.Errorf("empty container provided")
	}

	// load the container ID as the service key from the map
	result, ok := m.Load(c.ID)
	if !ok {
		return nil, fmt.Errorf("unable to load service %s", c.ID)
	}

	// cast the value from the service key to the expected type
	s, ok := result.(*api.Service)
	if !ok {
		return nil, fmt.Errorf("unable to cast value for service %s", c.ID)
	}

	return s, nil
}

// LoadLogs attempts to capture the library service logs
// representing the container from the map.
func LoadLogs(c *pipeline.Container, m *sync.Map) (*api.Log, error) {
	// check if the container provided is empty
	if c == nil {
		return nil, fmt.Errorf("empty container provided")
	}

	// load the container ID as the service log key from the map
	result, ok := m.Load(c.ID)
	if !ok {
		return nil, fmt.Errorf("unable to load logs for service %s", c.ID)
	}

	// cast the value from the service log key to the expected type
	l, ok := result.(*api.Log)
	if !ok {
		return nil, fmt.Errorf("unable to cast value to logs for service %s", c.ID)
	}

	return l, nil
}
