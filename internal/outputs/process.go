// SPDX-License-Identifier: Apache-2.0

package outputs

import (
	"maps"
	"strings"

	"github.com/go-vela/server/compiler/types/pipeline"
)

// Process adds the outputs to the container environment after sanitizing them
// and replacing $REFERENCES already in the environment. It also adds masked outputs
// to container secrets for masking in logs.
func Process(c *pipeline.Container, outputs, maskedOutputs map[string]string) {
	outputs = sanitize(c, outputs)
	maskedOutputs = sanitize(c, maskedOutputs)

	// combine sanitized maps for reference replacement
	allOutputs := make(map[string]string, len(outputs)+len(maskedOutputs))

	// copy outputs into single map
	maps.Copy(allOutputs, outputs)
	maps.Copy(allOutputs, maskedOutputs)

	// output reference replacement - only check values starting with $
	for k, v := range c.Environment {
		if strings.HasPrefix(v, "$") {
			if newV, ok := allOutputs[v[1:]]; ok {
				c.Environment[k] = newV
			}
		}
	}

	// add all outputs to container environment
	maps.Copy(c.Environment, allOutputs)

	// add masked outputs to container secrets for masking in logs
	if len(maskedOutputs) > 0 {
		outputSecrets := make([]*pipeline.StepSecret, 0, len(maskedOutputs))
		for key := range maskedOutputs {
			outputSecrets = append(outputSecrets, &pipeline.StepSecret{
				Target: key,
			})
		}

		c.Secrets = append(c.Secrets, outputSecrets...)
	}
}

func sanitize(c *pipeline.Container, outputs map[string]string) map[string]string {
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
