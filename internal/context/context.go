// SPDX-License-Identifier: Apache-2.0

package context

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Status string
var TimeFinished int64

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
		logger.Info("34 time is set")
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
				fmt.Println("There's content in stdout:", stdout.Text(), stderr.Text())
			} else if Status == "success" && time.Now().UTC().Unix() >= TimeFinished {
				cancel()
			}
		case <-ctx.Done():
			logger.Tracef("%s finished, stopping timeout timer", name)
			timer.Stop()

			return
		}
	}()

	//logger.Infof("CURRENT STATUS at 50 is %s", Status)
	return ctx, cancel
}

//
//for {
//// Sleep for 10 seconds
//time.Sleep(10 * time.Second)
//stdout, stderr := bufio.NewScanner(os.Stdout), bufio.NewScanner(os.Stderr)
//
//if stdout.Scan() || stderr.Scan() {
//fmt.Println("There's content in stdout:", stdout.Text(), stderr.Text())
//continue
//} else if c.build.GetStatus() == constants.StatusSuccess && time.Now().UTC().Unix() >= c.build.GetFinished() {
//c.Logger.Info("Build succeed and already finished")
//cancelStreaming()
//fmt.Println("There's no content in stdout.")
//break
//}
//}
