// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"testing"
)

func TestDistributionService_DistributionInspect(t *testing.T) {
	d := &DistributionService{}

	inspect, err := d.DistributionInspect(context.Background(), "test-image", "test-tag")
	if err != nil {
		t.Errorf("DistributionInspect() returned error: %v", err)
	}

	// Should return empty DistributionInspect struct
	if inspect.Descriptor.Digest != "" {
		t.Errorf("DistributionInspect() Descriptor.Digest = %v, want empty string", inspect.Descriptor.Digest)
	}

	if inspect.Descriptor.MediaType != "" {
		t.Errorf("DistributionInspect() Descriptor.MediaType = %v, want empty string", inspect.Descriptor.MediaType)
	}

	if inspect.Descriptor.Size != 0 {
		t.Errorf("DistributionInspect() Descriptor.Size = %v, want 0", inspect.Descriptor.Size)
	}

	// Test with different image names
	inspect2, err := d.DistributionInspect(context.Background(), "alpine", "latest")
	if err != nil {
		t.Errorf("DistributionInspect() with alpine:latest returned error: %v", err)
	}

	if inspect2.Descriptor.Digest != "" {
		t.Errorf("DistributionInspect() alpine:latest Descriptor.Digest = %v, want empty string", inspect2.Descriptor.Digest)
	}

	// Test with empty parameters
	inspect3, err := d.DistributionInspect(context.Background(), "", "")
	if err != nil {
		t.Errorf("DistributionInspect() with empty params returned error: %v", err)
	}

	if inspect3.Descriptor.Digest != "" {
		t.Errorf("DistributionInspect() empty params Descriptor.Digest = %v, want empty string", inspect3.Descriptor.Digest)
	}
}
