// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-vela/worker/version"
)

// hostname stores the worker host name reported by the kernel.
var hostname string

// create an init function to set the hostname for the worker.
//
// https://golang.org/doc/effective_go.html#init
func init() {
	// attempt to capture the hostname for the worker
	hostname, _ = os.Hostname()
	// check if a hostname is set
	if len(hostname) == 0 {
		// default the hostname to localhost
		hostname = "localhost"
	}
}

func main() {
	// capture application version information
	v := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	cmd := cli.Command{
		Name:    "vela-worker",
		Version: v.Semantic(),
		Action:  run,
		Usage:   "Vela build daemon designed for executing pipelines",
	}

	// Worker Flags

	cmd.Flags = flags()

	// Worker Start

	if err = cmd.Run(context.Background(), os.Args); err != nil {
		logrus.Fatal(err)
	}
}
