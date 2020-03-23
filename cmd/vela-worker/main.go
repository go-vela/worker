// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"os"
	"time"

	"github.com/go-vela/worker/version"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := cli.NewApp()

	// Worker Information

	app.Name = "vela-worker"
	app.HelpName = "vela-executor"
	app.Usage = "Vela executor package for integrating with different executors"
	app.Copyright = "Copyright (c) 2020 Target Brands, Inc. All rights reserved."
	app.Authors = []cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Worker Metadata

	app.Compiled = time.Now()
	app.Action = run
	app.Version = version.Version.String()

	// Worker Flags

	app.Flags = flags()

	// set logrus to log in JSON format
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Worker Start

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
