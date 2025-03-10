// SPDX-License-Identifier: Apache-2.0

package outputs

import (
	"github.com/go-vela/server/compiler/types/pipeline"
)

// Skip creates the ruledata from the build and repository
// information and returns true if the data does not match
// the ruleset for the given container.
func Sanitize(c *pipeline.Container, outputs map[string]string) map[string]string {
	// check if the container provided is empty
	if c == nil {
		return nil
	}

	if len(outputs) == 0 {
		return nil
	}

	// initialize sanitized outputs max length of outputs
	result := make(map[string]string, len(outputs))

	for k, v := range outputs {
		// do not allow overwrites to compiled starter env for step
		if _, ok := c.Environment[k]; ok {
			continue
		}

		result[k] = v
	}

	return result
}
