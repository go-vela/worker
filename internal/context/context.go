// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package context

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func WithDelayedCancelPropagation(parent context.Context, timeout time.Duration, name string, logger *logrus.Entry) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		var timer *time.Timer

		// start the timer once the parent context is canceled
		select {
		case <-parent.Done():
			logger.Tracef("%s timer for %s started now that parent context is done", name, timeout)
			timer = time.NewTimer(timeout)
		case <-ctx.Done():
			logger.Tracef("%s finished before the parent context", name)
			return
		}

		// wait for the timer to elapse or the context to naturally finish.
		select {
		case <-timer.C:
			logger.Tracef("%s timed out, canceling %s", name, name)
			cancel()
			return
		case <-ctx.Done():
			logger.Tracef("%s finished, stopping timeout timer", name)
			timer.Stop()
			return
		}
	}()

	return ctx, cancel
}
