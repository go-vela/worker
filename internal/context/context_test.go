// SPDX-License-Identifier: Apache-2.0

package context

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	logrusTest "github.com/sirupsen/logrus/hooks/test"
)

// shortDuration copied from
// https://github.com/golang/go/blob/go1.19.1/src/context/context_test.go#L45
const shortDuration = 1 * time.Millisecond // a reasonable duration to block in a test

// quiescent returns an arbitrary duration by which the program should have
// completed any remaining work and reached a steady (idle) state.
//
// copied from https://github.com/golang/go/blob/go1.19.1/src/context/context_test.go#L49-L59
func quiescent(t *testing.T) time.Duration {
	deadline, ok := t.Deadline()
	if !ok {
		return 5 * time.Second
	}

	const arbitraryCleanupMargin = 1 * time.Second

	return time.Until(deadline) - arbitraryCleanupMargin
}

// testCancelPropagated is a helper that tests deadline/timeouts.
//
// based on testDeadline from
// https://github.com/golang/go/blob/go1.19.1/src/context/context_test.go#L272-L285
func testCancelPropagated(c context.Context, name string, t *testing.T) {
	t.Helper()
	d := quiescent(t)

	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
		t.Fatalf("%s: context not timed out after %v", name, d)
	case <-c.Done():
	}
	//if e := c.Err(); e != context.DeadlineExceeded { // original line
	if e := c.Err(); !errors.Is(e, context.Canceled) { // delayedCancelPropagation triggers this instead
		t.Errorf("%s: c.Err() == %v; want %v", name, e, context.Canceled)
	}
}

func TestWithDelayedCancelPropagation(t *testing.T) {
	parentLogger := logrus.New()
	parentLogger.SetLevel(logrus.TraceLevel)

	loggerHook := logrusTest.NewLocal(parentLogger)

	nameArg := "streaming"

	tests := []struct {
		name           string
		cancelParent   string // before, after, never
		timeout        time.Duration
		cancelCtxAfter time.Duration
		lastLogMessage string
	}{
		{
			name:           "cancel parent before call-child finishes before timeout",
			cancelParent:   "before",
			timeout:        shortDuration * 5,
			cancelCtxAfter: shortDuration,
			lastLogMessage: nameArg + " finished, stopping timeout timer",
		},
		{
			name:           "cancel parent before call-child exceeds timeout",
			cancelParent:   "before",
			timeout:        shortDuration,
			cancelCtxAfter: shortDuration * 5,
			lastLogMessage: nameArg + " timed out, propagating cancel to " + nameArg + " context",
		},
		{
			name:           "child finished before timeout and before parent cancel",
			cancelParent:   "never",
			timeout:        shortDuration * 5,
			cancelCtxAfter: shortDuration,
			lastLogMessage: nameArg + " finished before the parent context",
		},
		{
			name:           "cancel parent after call-child finishes before timeout",
			cancelParent:   "after",
			timeout:        shortDuration * 5,
			cancelCtxAfter: shortDuration,
			lastLogMessage: nameArg + " finished, stopping timeout timer",
		},
		{
			name:           "cancel parent after call-child exceeds timeout",
			cancelParent:   "after",
			timeout:        shortDuration,
			cancelCtxAfter: shortDuration * 5,
			lastLogMessage: nameArg + " timed out, propagating cancel to " + nameArg + " context",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer loggerHook.Reset()

			logger := parentLogger.WithFields(logrus.Fields{"test": test.name})

			parentCtx, parentCtxCancel := context.WithCancel(context.Background())
			defer parentCtxCancel() // handles test.CancelParent == "never" case

			if test.cancelParent == "before" {
				parentCtxCancel()
			}

			ctx, cancel := WithDelayedCancelPropagation(parentCtx, test.timeout, nameArg, logger)
			defer cancel()

			// test based on test for context.WithCancel
			if got, want := fmt.Sprint(ctx), "context.Background.WithCancel"; got != want {
				t.Errorf("ctx.String() = %q want %q", got, want)
			}

			if d := ctx.Done(); d == nil {
				t.Errorf("ctx.Done() == %v want non-nil", d)
			}

			if e := ctx.Err(); e != nil {
				t.Errorf("ctx.Err() == %v want nil", e)
			}

			if test.cancelParent == "after" {
				parentCtxCancel()
			}

			go func() {
				time.Sleep(test.cancelCtxAfter)
				cancel()
			}()

			testCancelPropagated(ctx, "WithDelayedCancelPropagation", t)

			time.Sleep(shortDuration)

			lastLogEntry := loggerHook.LastEntry()
			if lastLogEntry.Message != test.lastLogMessage {
				t.Errorf("unexpected last log entry: want = %s ; got = %s", test.lastLogMessage, lastLogEntry.Message)
			}
		})
	}
}
