// SPDX-License-Identifier: Apache-2.0

package outputs

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/compiler/types/pipeline"
)

func TestOutputs_Process(t *testing.T) {
	tests := []struct {
		name          string
		ctn           *pipeline.Container
		outputs       map[string]string
		maskedOutputs map[string]string
		wantCtn       *pipeline.Container
	}{
		{
			name: "no outputs",
			ctn: &pipeline.Container{
				Environment: map[string]string{"FOO": "bar"},
			},
			outputs:       nil,
			maskedOutputs: nil,
			wantCtn: &pipeline.Container{
				Environment: map[string]string{"FOO": "bar"},
			},
		},
		{
			name: "outputs no references",
			ctn: &pipeline.Container{
				Environment: map[string]string{"FOO": "bar"},
			},
			outputs:       map[string]string{"OUTPUT1": "value1", "OUTPUT2": "value2"},
			maskedOutputs: nil,
			wantCtn: &pipeline.Container{
				Environment: map[string]string{"FOO": "bar", "OUTPUT1": "value1", "OUTPUT2": "value2"},
			},
		},
		{
			name: "outputs and masked outputs",
			ctn: &pipeline.Container{
				Environment: map[string]string{"FOO": "bar"},
			},
			outputs:       map[string]string{"OUTPUT1": "value1"},
			maskedOutputs: map[string]string{"MASKED_OUTPUT1": "masked_value1"},
			wantCtn: &pipeline.Container{
				Environment: map[string]string{"FOO": "bar", "OUTPUT1": "value1", "MASKED_OUTPUT1": "masked_value1"},
				Secrets: []*pipeline.StepSecret{
					{Target: "MASKED_OUTPUT1"},
				},
			},
		},
		{
			name: "outputs and masked outputs with references",
			ctn: &pipeline.Container{
				Environment: map[string]string{"FOO": "bar", "REF_OUTPUT": "$OUTPUT1", "REF_MASKED_OUTPUT": "$MASKED_OUTPUT1"},
			},
			outputs:       map[string]string{"OUTPUT1": "value1"},
			maskedOutputs: map[string]string{"MASKED_OUTPUT1": "masked_value1"},
			wantCtn: &pipeline.Container{
				Environment: map[string]string{"FOO": "bar", "REF_OUTPUT": "value1", "REF_MASKED_OUTPUT": "masked_value1", "OUTPUT1": "value1", "MASKED_OUTPUT1": "masked_value1"},
				Secrets: []*pipeline.StepSecret{
					{Target: "MASKED_OUTPUT1"},
				},
			},
		},
		{
			name: "outputs with attempted overwrite of existing env",
			ctn: &pipeline.Container{
				Environment: map[string]string{"FOO": "bar"},
			},
			outputs:       map[string]string{"FOO": "new_value"},
			maskedOutputs: nil,
			wantCtn: &pipeline.Container{
				Environment: map[string]string{"FOO": "bar"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(_ *testing.T) {
			Process(test.ctn, test.outputs, test.maskedOutputs)

			if diff := cmp.Diff(test.wantCtn, test.ctn); diff != "" {
				t.Errorf("(Process: -want +got):\n%s", diff)
			}
		})
	}
}
