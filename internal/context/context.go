// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package context

import (
	"context"
	"time"
)

func WithDelayedCancelPropagation(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		var timer *time.Timer

		// start the timer once the parent context is canceled
		select {
		case <-parent.Done():
			timer = time.NewTimer(timeout)
		case <-ctx.Done():
			return
		}

		// wait for the timer to elapse or the context to naturally finish.
		select {
		case <-timer.C:
			cancel()
			return
		case <-ctx.Done():
			timer.Stop()
			return
		}
	}()

	return ctx, cancel
}
