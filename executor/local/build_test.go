// SPDX-License-Identifier: Apache-2.0

package local

import (
	"context"
	"testing"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/compiler/native"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/worker/internal/message"
	"github.com/go-vela/worker/runtime/docker"
)

func TestLocal_CreateBuild(t *testing.T) {
	// setup types
	cmd := new(cli.Command)
	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "clone-image",
			Value: "target/vela-git:latest",
			Usage: "doc",
		},
	}
	compiler, err := native.FromCLICommand(context.Background(), cmd)
	if err != nil {
		t.Errorf("unable to create compiler engine: %v", err)
	}

	_build := testBuild()

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	tests := []struct {
		name     string
		failure  bool
		pipeline string
	}{
		{
			name:     "basic services pipeline",
			failure:  false,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "basic steps pipeline",
			failure:  false,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "basic stages pipeline",
			failure:  false,
			pipeline: "testdata/build/stages/basic.yml",
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_build.GetRepo()).
				WithLocal(true).
				WithUser(_build.GetRepo().GetOwner()).
				Compile(context.Background(), test.pipeline)
			if err != nil {
				t.Errorf("unable to compile pipeline %s: %v", test.pipeline, err)
			}

			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRuntime(_runtime),
				WithOutputCtn(testOutputsCtn()),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			err = _engine.CreateBuild(context.Background())

			if test.failure {
				if err == nil {
					t.Errorf("CreateBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("CreateBuild returned err: %v", err)
			}
		})
	}
}

func TestLocal_PlanBuild(t *testing.T) {
	// setup types
	cmd := new(cli.Command)
	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "clone-image",
			Value: "target/vela-git:latest",
			Usage: "doc",
		},
	}
	compiler, err := native.FromCLICommand(context.Background(), cmd)
	if err != nil {
		t.Errorf("unable to create compiler engine: %v", err)
	}

	_build := testBuild()

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	tests := []struct {
		name     string
		failure  bool
		pipeline string
	}{
		{
			name:     "basic services pipeline",
			failure:  false,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "basic steps pipeline",
			failure:  false,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "basic stages pipeline",
			failure:  false,
			pipeline: "testdata/build/stages/basic.yml",
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_build.GetRepo()).
				WithLocal(true).
				WithUser(_build.GetRepo().GetOwner()).
				Compile(context.Background(), test.pipeline)
			if err != nil {
				t.Errorf("unable to compile pipeline %s: %v", test.pipeline, err)
			}

			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRuntime(_runtime),
				WithOutputCtn(testOutputsCtn()),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			// run create to init steps to be created properly
			err = _engine.CreateBuild(context.Background())
			if err != nil {
				t.Errorf("unable to create build: %v", err)
			}

			err = _engine.PlanBuild(context.Background())

			if test.failure {
				if err == nil {
					t.Errorf("PlanBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("PlanBuild returned err: %v", err)
			}
		})
	}
}

func TestLocal_AssembleBuild(t *testing.T) {
	// setup types
	cmd := new(cli.Command)
	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "clone-image",
			Value: "target/vela-git:latest",
			Usage: "doc",
		},
	}
	compiler, err := native.FromCLICommand(context.Background(), cmd)
	if err != nil {
		t.Errorf("unable to create compiler engine: %v", err)
	}

	_build := testBuild()

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	streamRequests, done := message.MockStreamRequestsWithCancel(context.Background())
	defer done()

	tests := []struct {
		name     string
		failure  bool
		pipeline string
	}{
		{
			name:     "basic services pipeline",
			failure:  false,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "services pipeline with image not found",
			failure:  true,
			pipeline: "testdata/build/services/img_notfound.yml",
		},
		{
			name:     "services pipeline with ignoring image not found",
			failure:  true,
			pipeline: "testdata/build/services/img_ignorenotfound.yml",
		},
		{
			name:     "basic steps pipeline",
			failure:  false,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "steps pipeline with image not found",
			failure:  true,
			pipeline: "testdata/build/steps/img_notfound.yml",
		},
		{
			name:     "steps pipeline with ignoring image not found",
			failure:  true,
			pipeline: "testdata/build/steps/img_ignorenotfound.yml",
		},
		{
			name:     "basic stages pipeline",
			failure:  false,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "stages pipeline with image not found",
			failure:  true,
			pipeline: "testdata/build/stages/img_notfound.yml",
		},
		{
			name:     "stages pipeline with ignoring image not found",
			failure:  true,
			pipeline: "testdata/build/stages/img_ignorenotfound.yml",
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_build.GetRepo()).
				WithLocal(true).
				WithUser(_build.GetRepo().GetOwner()).
				Compile(context.Background(), test.pipeline)
			if err != nil {
				t.Errorf("unable to compile pipeline %s: %v", test.pipeline, err)
			}

			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRuntime(_runtime),
				WithOutputCtn(testOutputsCtn()),
				withStreamRequests(streamRequests),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			// run create to init steps to be created properly
			err = _engine.CreateBuild(context.Background())
			if err != nil {
				t.Errorf("unable to create build: %v", err)
			}

			err = _engine.AssembleBuild(context.Background())

			if test.failure {
				if err == nil {
					t.Errorf("AssembleBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("AssembleBuild returned err: %v", err)
			}
		})
	}
}

func TestLocal_ExecBuild(t *testing.T) {
	// setup types
	cmd := new(cli.Command)
	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "clone-image",
			Value: "target/vela-git:latest",
			Usage: "doc",
		},
	}
	compiler, err := native.FromCLICommand(context.Background(), cmd)
	if err != nil {
		t.Errorf("unable to create compiler engine: %v", err)
	}

	_build := testBuild()

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	streamRequests, done := message.MockStreamRequestsWithCancel(context.Background())
	defer done()

	tests := []struct {
		name     string
		failure  bool
		pipeline string
	}{
		{
			name:     "basic services pipeline",
			failure:  false,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "services pipeline with image not found",
			failure:  true,
			pipeline: "testdata/build/services/img_notfound.yml",
		},
		{
			name:     "basic steps pipeline",
			failure:  false,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "steps pipeline with image not found",
			failure:  true,
			pipeline: "testdata/build/steps/img_notfound.yml",
		},
		{
			name:     "basic stages pipeline",
			failure:  false,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "stages pipeline with image not found",
			failure:  true,
			pipeline: "testdata/build/stages/img_notfound.yml",
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_build.GetRepo()).
				WithLocal(true).
				WithUser(_build.GetRepo().GetOwner()).
				Compile(context.Background(), test.pipeline)
			if err != nil {
				t.Errorf("unable to compile pipeline %s: %v", test.pipeline, err)
			}

			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRuntime(_runtime),
				WithOutputCtn(testOutputsCtn()),
				withStreamRequests(streamRequests),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			// run create to init steps to be created properly
			err = _engine.CreateBuild(context.Background())
			if err != nil {
				t.Errorf("unable to create build: %v", err)
			}

			err = _engine.ExecBuild(context.Background())

			if test.failure {
				if err == nil {
					t.Errorf("ExecBuild for %s should have returned err", test.pipeline)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("ExecBuild for %s returned err: %v", test.pipeline, err)
			}
		})
	}
}

func TestLocal_StreamBuild(t *testing.T) {
	// setup types
	cmd := new(cli.Command)
	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "clone-image",
			Value: "target/vela-git:latest",
			Usage: "doc",
		},
	}
	compiler, err := native.FromCLICommand(context.Background(), cmd)
	if err != nil {
		t.Errorf("unable to create compiler engine: %v", err)
	}

	_build := testBuild()

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	type planFuncType = func(context.Context, *pipeline.Container) error

	// planNothing is a planFuncType that does nothing
	planNothing := func(ctx context.Context, container *pipeline.Container) error {
		return nil
	}

	tests := []struct {
		name       string
		failure    bool
		pipeline   string
		messageKey string
		ctn        *pipeline.Container
		streamFunc func(*client) message.StreamFunc
		planFunc   func(*client) planFuncType
	}{
		{
			name:       "basic services pipeline",
			failure:    false,
			pipeline:   "testdata/build/services/basic.yml",
			messageKey: "service",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamService
			},
			planFunc: func(c *client) planFuncType {
				return c.PlanService
			},
			ctn: &pipeline.Container{
				ID:          "service_github_octocat_1_postgres",
				Detach:      true,
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:latest",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		{
			name:       "basic services pipeline with StreamService failure",
			failure:    false,
			pipeline:   "testdata/build/services/basic.yml",
			messageKey: "service",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamService
			},
			planFunc: func(c *client) planFuncType {
				// simulate failure to call PlanService
				return planNothing
			},
			ctn: &pipeline.Container{
				ID:          "service_github_octocat_1_postgres",
				Detach:      true,
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:latest",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
				Pull:        "not_present",
			},
		},
		{
			name:       "basic steps pipeline",
			failure:    false,
			pipeline:   "testdata/build/steps/basic.yml",
			messageKey: "step",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamStep
			},
			planFunc: func(c *client) planFuncType {
				return c.PlanStep
			},
			ctn: &pipeline.Container{
				ID:          "step_github_octocat_1_test",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "test",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:       "basic steps pipeline with StreamStep failure",
			failure:    false,
			pipeline:   "testdata/build/steps/basic.yml",
			messageKey: "step",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamStep
			},
			planFunc: func(c *client) planFuncType {
				// simulate failure to call PlanStep
				return planNothing
			},
			ctn: &pipeline.Container{
				ID:          "step_github_octocat_1_test",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "test",
				Number:      1,
				Pull:        "not_present",
			},
		},
		{
			name:       "basic stages pipeline",
			failure:    false,
			pipeline:   "testdata/build/stages/basic.yml",
			messageKey: "step",
			streamFunc: func(c *client) message.StreamFunc {
				return c.StreamStep
			},
			planFunc: func(c *client) planFuncType {
				return c.PlanStep
			},
			ctn: &pipeline.Container{
				ID:          "step_github_octocat_1_test_test",
				Directory:   "/vela/src/github.com/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "test",
				Number:      1,
				Pull:        "not_present",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buildCtx, done := context.WithCancel(context.Background())
			defer done()

			streamRequests := make(chan message.StreamRequest)

			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_build.GetRepo()).
				WithLocal(true).
				WithUser(_build.GetRepo().GetOwner()).
				Compile(context.Background(), test.pipeline)
			if err != nil {
				t.Errorf("unable to compile pipeline %s: %v", test.pipeline, err)
			}

			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRuntime(_runtime),
				WithOutputCtn(testOutputsCtn()),
				withStreamRequests(streamRequests),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			// run create to init steps to be created properly
			err = _engine.CreateBuild(buildCtx)
			if err != nil {
				t.Errorf("unable to create build: %v", err)
			}

			// simulate ExecBuild()
			go func() {
				// ExecBuild calls PlanService()/PlanStep() before ExecService()/ExecStep()
				// (ExecStage() calls PlanStep() before ExecStep()).
				_engine.err = test.planFunc(_engine)(buildCtx, test.ctn)

				// ExecService()/ExecStep()/secret.exec() send this message
				streamRequests <- message.StreamRequest{
					Key:       test.messageKey,
					Stream:    test.streamFunc(_engine),
					Container: test.ctn,
				}

				// simulate exec build duration
				time.Sleep(100 * time.Microsecond)

				// signal the end of the build so StreamBuild can terminate
				done()
			}()

			err = _engine.StreamBuild(buildCtx)

			if test.failure {
				if err == nil {
					t.Errorf("StreamBuild for %s should have returned err", test.pipeline)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("StreamBuild for %s returned err: %v", test.pipeline, err)
			}
		})
	}
}

func TestLocal_DestroyBuild(t *testing.T) {
	// setup types
	cmd := new(cli.Command)
	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "clone-image",
			Value: "target/vela-git:latest",
			Usage: "doc",
		},
	}
	compiler, err := native.FromCLICommand(context.Background(), cmd)
	if err != nil {
		t.Errorf("unable to create compiler engine: %v", err)
	}

	_build := testBuild()

	_runtime, err := docker.NewMock()
	if err != nil {
		t.Errorf("unable to create runtime engine: %v", err)
	}

	tests := []struct {
		name     string
		failure  bool
		pipeline string
	}{
		{
			name:     "basic services pipeline",
			failure:  false,
			pipeline: "testdata/build/services/basic.yml",
		},
		{
			name:     "services pipeline with name not found",
			failure:  false,
			pipeline: "testdata/build/services/name_notfound.yml",
		},
		{
			name:     "basic steps pipeline",
			failure:  false,
			pipeline: "testdata/build/steps/basic.yml",
		},
		{
			name:     "steps pipeline with name not found",
			failure:  false,
			pipeline: "testdata/build/steps/name_notfound.yml",
		},
		{
			name:     "basic stages pipeline",
			failure:  false,
			pipeline: "testdata/build/stages/basic.yml",
		},
		{
			name:     "stages pipeline with name not found",
			failure:  false,
			pipeline: "testdata/build/stages/name_notfound.yml",
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_pipeline, _, err := compiler.
				Duplicate().
				WithBuild(_build).
				WithRepo(_build.GetRepo()).
				WithLocal(true).
				WithUser(_build.GetRepo().GetOwner()).
				Compile(context.Background(), test.pipeline)
			if err != nil {
				t.Errorf("unable to compile pipeline %s: %v", test.pipeline, err)
			}

			_engine, err := New(
				WithBuild(_build),
				WithPipeline(_pipeline),
				WithRuntime(_runtime),
				WithOutputCtn(testOutputsCtn()),
			)
			if err != nil {
				t.Errorf("unable to create executor engine: %v", err)
			}

			// run create to init steps to be created properly
			err = _engine.CreateBuild(context.Background())
			if err != nil {
				t.Errorf("unable to create build: %v", err)
			}

			err = _engine.DestroyBuild(context.Background())

			if test.failure {
				if err == nil {
					t.Errorf("DestroyBuild should have returned err")
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("DestroyBuild returned err: %v", err)
			}
		})
	}
}

func testOutputsCtn() *pipeline.Container {
	return &pipeline.Container{
		ID:          "outputs_test",
		Environment: make(map[string]string),
		Detach:      true,
	}
}
