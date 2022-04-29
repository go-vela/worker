// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package message

import (
	"context"

	"github.com/go-vela/types/pipeline"
)

// StreamFunc is either StreamService or StreamStep in executor.Engine.
type StreamFunc = func(context.Context, *pipeline.Container) error

// StreamRequest is the message used to begin streaming for a container
// (requests goes from ExecService / ExecStep to StreamBuild in executor).
type StreamRequest struct {
	// Key is either "service" or "step".
	Key string
	// Stream is either Engine.StreamService or Engine.StreamStep.
	Stream StreamFunc
	// Container is the container for the service or step to stream logs for.
	Container *pipeline.Container
}

// MockStreamRequestsWithCancel discards all requests until you call the cancel function.
func MockStreamRequestsWithCancel(ctx context.Context) (chan StreamRequest, context.CancelFunc) {
	cancelCtx, done := context.WithCancel(ctx)
	streamRequests := make(chan StreamRequest)

	// discard all stream requests
	go func() {
		for {
			select {
			case <-streamRequests:
			case <-cancelCtx.Done():
				return
			}
		}
	}()

	return streamRequests, done
}
