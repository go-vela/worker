// SPDX-License-Identifier: Apache-2.0

package context

import (
	"bufio"
	"context"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	status       string
	timeFinished int64
	muS          sync.Mutex
	muTF         sync.Mutex
)

func SetBuildStatus(s string) {
	muS.Lock()
	defer muS.Unlock()
	status = s
}

func GetBuildStatus() string {
	muS.Lock()
	defer muS.Unlock()
	return status
}

func SetBuildFinished(tf int64) {
	muTF.Lock()
	defer muTF.Unlock()
	timeFinished = tf
}

func GetBuildFinished() int64 {
	muTF.Lock()
	defer muTF.Unlock()
	return timeFinished
}

func WithDelayedCancelPropagation(parent context.Context, timeout time.Duration, name string, logger *logrus.Entry) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		var timer *time.Timer
		// Create a ticker with a 1-second interval
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
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
		// stop time ticker once finished.
		select {
		case <-timer.C:
			logger.Tracef("%s timed out, propagating cancel to %s context", name, name)
			ticker.Stop()
			cancel()

			return
		case <-ticker.C:
			stdout, stderr := bufio.NewScanner(os.Stdout), bufio.NewScanner(os.Stderr)

			if stdout.Scan() || stderr.Scan() {
				logger.Debug("Logs found, continuing")
			} else if GetBuildStatus() == "success" && time.Now().UTC().Unix() >= GetBuildFinished() {
				logger.Debug("")
				cancel()
			}
		case <-ctx.Done():
			logger.Tracef("%s finished, stopping timeout timer", name)
			timer.Stop()

			return
		}
	}()

	return ctx, cancel
}
