// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/types/constants"

	"github.com/go-vela/worker/queue"
	"github.com/go-vela/worker/queue/redis"

	"github.com/sirupsen/logrus"
)

// helper function to setup the queue from the CLI arguments.
func setupQueue(q *queueSetup) (queue.Service, error) {
	logrus.Debug("Creating queue client from CLI configuration")

	switch q.Driver {
	case constants.DriverKafka:
		return setupKafka(q)
	case constants.DriverRedis:
		return setupRedis(q)
	default:
		return nil, fmt.Errorf("invalid queue driver: %s", q.Driver)
	}
}

// helper function to setup the Kafka queue from the CLI arguments.
func setupKafka(q *queueSetup) (queue.Service, error) {
	logrus.Tracef("Creating %s queue client from CLI configuration", constants.DriverKafka)
	// return kafka.New(c.String("queue-config"), "vela")
	return nil, fmt.Errorf("unsupported queue driver: %s", constants.DriverKafka)
}

// helper function to setup the Redis queue from the CLI arguments.
func setupRedis(q *queueSetup) (queue.Service, error) {
	// setup routes
	routes := append(q.Routes, constants.DefaultRoute)

	if q.Cluster {
		logrus.Tracef("Creating %s queue cluster client from CLI configuration", constants.DriverRedis)
		return redis.NewCluster(q.Config, routes...)
	}

	logrus.Tracef("Creating %s queue client from CLI configuration", constants.DriverRedis)

	return redis.New(q.Config, routes...)
}
