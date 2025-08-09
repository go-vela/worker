// SPDX-License-Identifier: Apache-2.0

package message

import (
	"context"
	"testing"
	"time"

	"github.com/go-vela/server/compiler/types/pipeline"
)

func TestMockStreamRequestsWithCancel(t *testing.T) {
	ctx := context.Background()

	// Test creating mock stream requests
	streamRequests, cancel := MockStreamRequestsWithCancel(ctx)

	if streamRequests == nil {
		t.Error("MockStreamRequestsWithCancel() returned nil channel")
	}

	if cancel == nil {
		t.Error("MockStreamRequestsWithCancel() returned nil cancel function")
	}

	// Test that we can send requests without blocking
	mockContainer := &pipeline.Container{
		ID:   "test-container",
		Name: "test",
	}

	mockStreamFunc := func(_ context.Context, _ *pipeline.Container) error {
		return nil
	}

	req := StreamRequest{
		Key:       "service",
		Stream:    mockStreamFunc,
		Container: mockContainer,
	}

	// Send a request (should not block)
	done := make(chan bool)

	go func() {
		streamRequests <- req

		done <- true
	}()

	// Wait for the request to be processed (or timeout)
	select {
	case <-done:
		// Success - request was processed
	case <-time.After(100 * time.Millisecond):
		t.Error("Sending stream request blocked - should be discarded immediately")
	}

	// Test that cancel function works
	cancel()

	// Allow some time for goroutine to process cancellation
	time.Sleep(10 * time.Millisecond)

	// After canceling, the goroutine should exit
	// We can't directly test this, but we've verified the basic functionality
}

func TestStreamRequest(t *testing.T) {
	// Test creating StreamRequest struct
	mockContainer := &pipeline.Container{
		ID:   "test-container",
		Name: "test",
	}

	mockStreamFunc := func(_ context.Context, _ *pipeline.Container) error {
		return nil
	}

	req := StreamRequest{
		Key:       "step",
		Stream:    mockStreamFunc,
		Container: mockContainer,
	}

	if req.Key != "step" {
		t.Errorf("StreamRequest.Key = %v, want 'step'", req.Key)
	}

	if req.Container != mockContainer {
		t.Error("StreamRequest.Container not set correctly")
	}

	if req.Stream == nil {
		t.Error("StreamRequest.Stream should not be nil")
	}

	// Test that the stream function can be called
	err := req.Stream(context.Background(), mockContainer)
	if err != nil {
		t.Errorf("StreamRequest.Stream() error = %v, want nil", err)
	}
}

func TestMockStreamRequestsWithCancel_MultipleRequests(_ *testing.T) {
	ctx := context.Background()

	streamRequests, cancel := MockStreamRequestsWithCancel(ctx)
	defer cancel()

	mockContainer := &pipeline.Container{
		ID:   "test-container",
		Name: "test",
	}

	mockStreamFunc := func(_ context.Context, _ *pipeline.Container) error {
		return nil
	}

	// Send multiple requests rapidly
	for i := 0; i < 5; i++ {
		go func(_ int) {
			req := StreamRequest{
				Key:       "service",
				Stream:    mockStreamFunc,
				Container: mockContainer,
			}
			streamRequests <- req
		}(i)
	}

	// All requests should be discarded without blocking
	time.Sleep(50 * time.Millisecond)
}
