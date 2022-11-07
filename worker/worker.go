// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"net/url"
	"sync"
	"time"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime"
	"github.com/sirupsen/logrus"
)

type (
	// API represents the worker configuration for API information.
	API struct {
		Address *url.URL
	}

	// Build represents the worker configuration for build information.
	Build struct {
		Limit   int
		Timeout time.Duration
	}

	// Logger represents the worker configuration for logger information.
	Logger struct {
		Format string
		Level  string
	}

	// Server represents the worker configuration for server information.
	Server struct {
		Address string
		Secret  string
	}

	// Certificate represents the optional cert and key to enable TLS.
	Certificate struct {
		Cert string
		Key  string
	}

	// Config represents the worker configuration.
	Config struct {
		Mock          bool // Mock should only be true for tests
		API           *API
		Build         *Build
		CheckIn       time.Duration
		Executor      *executor.Setup
		Logger        *Logger
		Queue         *queue.Setup
		Runtime       *runtime.Setup
		Server        *Server
		Certificate   *Certificate
		TLSMinVersion string
	}

	// Worker represents all configuration and
	// system processes for the worker.
	Worker struct {
		Config         *Config
		Executors      map[int]executor.Engine
		Queue          queue.Service
		Runtime        runtime.Engine
		VelaClient     *vela.Client
		PackageChannel chan *types.BuildPackage
		*Activity
	}

	Activity struct {
		ActiveBuilds []*library.Build
		Mutex        sync.Mutex
		Channel      chan message.BuildActivity
	}
)

func NewActivity() *Activity {
	return &Activity{
		ActiveBuilds: []*library.Build{},
		Mutex:        sync.Mutex{},
		Channel:      make(chan message.BuildActivity),
	}
}

func (a *Activity) ToWorkerStatus(w *Worker) string {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	if len(a.ActiveBuilds) < w.Config.Build.Limit {
		return "ready"
	} else if len(a.ActiveBuilds) > 0 {
		return "building"
	} else {
		return "error"
	}
}

func (a *Activity) HandleMessage(msg *message.BuildActivity) {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	// handle message
	switch t := msg.Action.(type) {
	case message.AddBuild:
		a.AddBuild(msg.Build)
	case message.RemoveBuild:
		a.RemoveBuild(msg.Build)
	default:
		logrus.Tracef("received unsupported build activity message %s", t)
	}
}

func (a *Activity) GetBuild(build *library.Build) (*library.Build, int) {
	var _build *library.Build
	idx := -1
	for i, b := range a.ActiveBuilds {
		if b.GetID() == build.GetID() {
			_build = b
			idx = i
		}
	}
	return _build, idx
}

func (a *Activity) AddBuild(build *library.Build) {
	// check activity for incoming build
	_build, idx := a.GetBuild(build)

	// build found
	if _build != nil || idx != -1 {
		return
	}

	// add build
	a.ActiveBuilds = append(a.ActiveBuilds, build)
}

func (a *Activity) RemoveBuild(build *library.Build) {
	// check activity for incoming build
	_build, idx := a.GetBuild(build)

	// build not found
	if _build == nil || idx == -1 {
		return
	}

	// remove build
	a.ActiveBuilds[idx] = a.ActiveBuilds[len(a.ActiveBuilds)-1]
	a.ActiveBuilds = a.ActiveBuilds[:len(a.ActiveBuilds)-1]
}
