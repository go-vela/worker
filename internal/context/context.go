// SPDX-License-Identifier: Apache-2.0

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
			logger.Tracef("parent context is done, starting %s timer for %s", name, timeout)
			timer = time.NewTimer(timeout)

			break
		case <-ctx.Done():
			logger.Tracef("%s finished before the parent context", name)

			return
		}

		// wait for the timer to elapse or the context to naturally finish.
		select {
		case <-timer.C:
			logger.Tracef("%s timed out, propagating cancel to %s context", name, name)
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
