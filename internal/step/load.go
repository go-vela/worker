// SPDX-License-Identifier: Apache-2.0

package step

import (
	"fmt"
	"sync"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
)

// Load attempts to capture the library step
// representing the container from the map.
func Load(c *pipeline.Container, m *sync.Map) (*api.Step, error) {
	// check if the container provided is empty
	if c == nil {
		return nil, fmt.Errorf("empty container provided")
	}

	// load the container ID as the step key from the map
	result, ok := m.Load(c.ID)
	if !ok {
		return nil, fmt.Errorf("unable to load step %s", c.ID)
	}

	// cast the value from the step key to the expected type
	s, ok := result.(*api.Step)
	if !ok {
		return nil, fmt.Errorf("unable to cast value for step %s", c.ID)
	}

	return s, nil
}

// LoadInit attempts to capture the container representing
// the init process from the pipeline.
func LoadInit(p *pipeline.Build) (*pipeline.Container, error) {
	// check if the pipeline provided is empty
	if p == nil {
		return nil, fmt.Errorf("empty pipeline provided")
	}

	// create new container for the init step
	c := new(pipeline.Container)

	// check if there are steps in the pipeline
	if len(p.Steps) > 0 {
		// update the container for the init process
		c = p.Steps[0]
	}

	// check if there are stages in the pipeline
	if len(p.Stages) > 0 {
		// update the container for the init process
		c = p.Stages[0].Steps[0]
	}

	return c, nil
}

// LoadLogs attempts to capture the library step logs
// representing the container from the map.
func LoadLogs(c *pipeline.Container, m *sync.Map) (*api.Log, error) {
	// check if the container provided is empty
	if c == nil {
		return nil, fmt.Errorf("empty container provided")
	}

	// load the container ID as the step log key from the map
	result, ok := m.Load(c.ID)
	if !ok {
		return nil, fmt.Errorf("unable to load logs for step %s", c.ID)
	}

	// cast the value from the step log key to the expected type
	l, ok := result.(*api.Log)
	if !ok {
		return nil, fmt.Errorf("unable to cast value to logs for step %s", c.ID)
	}

	return l, nil
}
