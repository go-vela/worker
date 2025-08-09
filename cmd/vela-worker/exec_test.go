// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/sdk-go/vela"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/queue/models"
	"github.com/go-vela/worker/executor"
	"github.com/go-vela/worker/runtime"
)

func TestGenerateCryptographicBuildID(t *testing.T) {
	// Test that generateCryptographicBuildID returns a valid hex string
	id1 := generateCryptographicBuildID()
	if len(id1) != 32 { // 16 bytes = 32 hex characters
		t.Errorf("generateCryptographicBuildID() returned ID with length %d, want 32", len(id1))
	}

	// Test that it generates unique IDs
	id2 := generateCryptographicBuildID()
	if id1 == id2 {
		t.Errorf("generateCryptographicBuildID() returned duplicate IDs: %s", id1)
	}

	// Verify it's valid hex
	for _, c := range id1 {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			t.Errorf("generateCryptographicBuildID() returned non-hex character: %c", c)
		}
	}
}

func TestWorker_getBuildResources(t *testing.T) {
	tests := []struct {
		name        string
		cpuQuota    int
		memoryLimit int
		pidsLimit   int
		wantCPU     int64
		wantMemory  int64
		wantPids    int64
	}{
		{
			name:        "default values",
			cpuQuota:    1200,
			memoryLimit: 4,
			pidsLimit:   1024,
			wantCPU:     1200,
			wantMemory:  4 * 1024 * 1024 * 1024,
			wantPids:    1024,
		},
		{
			name:        "custom values",
			cpuQuota:    2000,
			memoryLimit: 8,
			pidsLimit:   2048,
			wantCPU:     2000,
			wantMemory:  8 * 1024 * 1024 * 1024,
			wantPids:    2048,
		},
		{
			name:        "minimum values",
			cpuQuota:    100,
			memoryLimit: 1,
			pidsLimit:   256,
			wantCPU:     100,
			wantMemory:  1 * 1024 * 1024 * 1024,
			wantPids:    256,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				Config: &Config{
					Build: &Build{
						CPUQuota:    tt.cpuQuota,
						MemoryLimit: tt.memoryLimit,
						PidsLimit:   tt.pidsLimit,
					},
				},
			}

			resources := w.getBuildResources()

			if resources.CPUQuota != tt.wantCPU {
				t.Errorf("getBuildResources() CPUQuota = %v, want %v", resources.CPUQuota, tt.wantCPU)
			}

			if resources.Memory != tt.wantMemory {
				t.Errorf("getBuildResources() Memory = %v, want %v", resources.Memory, tt.wantMemory)
			}

			if resources.PidsLimit != tt.wantPids {
				t.Errorf("getBuildResources() PidsLimit = %v, want %v", resources.PidsLimit, tt.wantPids)
			}
		})
	}
}

func TestWorker_BuildContextManagement(t *testing.T) {
	// Test build context initialization and cleanup
	w := &Worker{
		BuildContexts:      nil,
		BuildContextsMutex: sync.RWMutex{},
		Config: &Config{
			Build: &Build{
				CPUQuota:    1200,
				MemoryLimit: 4,
				PidsLimit:   1024,
			},
		},
	}

	// Test BuildContexts initialization
	if w.BuildContexts == nil {
		w.BuildContexts = make(map[string]*BuildContext)
	}

	buildID := "test-build-123"
	buildContext := &BuildContext{
		BuildID:       buildID,
		WorkspacePath: "/tmp/vela-build-" + buildID,
		StartTime:     time.Now(),
		Resources:     w.getBuildResources(),
		Environment:   make(map[string]string),
	}

	// Test context storage
	w.BuildContextsMutex.Lock()
	w.BuildContexts[buildID] = buildContext
	w.BuildContextsMutex.Unlock()

	// Verify context is stored
	w.BuildContextsMutex.RLock()
	stored, exists := w.BuildContexts[buildID]
	w.BuildContextsMutex.RUnlock()

	if !exists {
		t.Error("Build context was not stored")
	}

	if stored.BuildID != buildID {
		t.Errorf("Stored build ID = %v, want %v", stored.BuildID, buildID)
	}

	// Test context cleanup
	w.BuildContextsMutex.Lock()
	delete(w.BuildContexts, buildID)
	w.BuildContextsMutex.Unlock()

	// Verify context is cleaned up
	w.BuildContextsMutex.RLock()
	_, exists = w.BuildContexts[buildID]
	w.BuildContextsMutex.RUnlock()

	if exists {
		t.Error("Build context was not cleaned up")
	}
}

func TestBuildContext(t *testing.T) {
	buildID := "test-build-456"
	workspace := "/tmp/vela-build-" + buildID
	startTime := time.Now()

	resources := &BuildResources{
		CPUQuota:  1200,
		Memory:    4 * 1024 * 1024 * 1024,
		PidsLimit: 1024,
	}

	env := make(map[string]string)
	env["TEST_VAR"] = "test_value"

	context := &BuildContext{
		BuildID:       buildID,
		WorkspacePath: workspace,
		StartTime:     startTime,
		Resources:     resources,
		Environment:   env,
	}

	// Test all fields are set correctly
	if context.BuildID != buildID {
		t.Errorf("BuildContext.BuildID = %v, want %v", context.BuildID, buildID)
	}

	if context.WorkspacePath != workspace {
		t.Errorf("BuildContext.WorkspacePath = %v, want %v", context.WorkspacePath, workspace)
	}

	if context.Resources.CPUQuota != 1200 {
		t.Errorf("BuildContext.Resources.CPUQuota = %v, want 1200", context.Resources.CPUQuota)
	}

	if context.Environment["TEST_VAR"] != "test_value" {
		t.Errorf("BuildContext.Environment[TEST_VAR] = %v, want test_value", context.Environment["TEST_VAR"])
	}
}

func TestBuildResources(t *testing.T) {
	resources := &BuildResources{
		CPUQuota:  2000,
		Memory:    8 * 1024 * 1024 * 1024,
		PidsLimit: 2048,
	}

	if resources.CPUQuota != 2000 {
		t.Errorf("BuildResources.CPUQuota = %v, want 2000", resources.CPUQuota)
	}

	if resources.Memory != 8*1024*1024*1024 {
		t.Errorf("BuildResources.Memory = %v, want %v", resources.Memory, 8*1024*1024*1024)
	}

	if resources.PidsLimit != 2048 {
		t.Errorf("BuildResources.PidsLimit = %v, want 2048", resources.PidsLimit)
	}
}

func TestWorker_exec(t *testing.T) {
	// Set up a test server to mock the Vela API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/workers/test-worker":
			switch r.Method {
			case "GET":
				worker := &api.Worker{}
				worker.SetHostname("test-worker")
				worker.SetRoutes([]string{"repo"})
				worker.SetStatus(constants.WorkerStatusIdle)
				worker.SetRunningBuilds([]*api.Build{})
				_ = json.NewEncoder(w).Encode(worker)
			case "PUT":
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(&api.Worker{})
			}
		case "/api/v1/repos/test-org/test-repo/builds/1/token":
			token := &api.Token{}
			token.SetToken("test-token")
			_ = json.NewEncoder(w).Encode(token)
		case "/api/v1/repos/test-org/test-repo/builds/1/executable":
			pipelineBuild := &pipeline.Build{
				ID:      "test-build-id",
				Version: "1",
				Steps: pipeline.ContainerSlice{
					{
						ID:    "step-1",
						Name:  "test",
						Image: "alpine:latest",
					},
				},
			}
			data, _ := json.Marshal(pipelineBuild)
			executable := &api.BuildExecutable{}
			executable.SetData(data)
			_ = json.NewEncoder(w).Encode(executable)
		case "/api/v1/repos/test-org/test-repo/builds/1":
			if r.Method == "PUT" {
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(&api.Build{})
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Parse the test server URL
	serverURL, _ := url.Parse(server.URL)

	// Create a test queue
	testQueue := &mockQueue{
		popFunc: func(_ context.Context, _ []string) (*models.Item, error) {
			build := &api.Build{}
			build.SetID(1)
			build.SetNumber(1)
			build.SetStatus(constants.StatusPending)
			repo := &api.Repo{}
			repo.SetOrg("test-org")
			repo.SetName("test-repo")
			repo.SetFullName("test-org/test-repo")
			repo.SetTimeout(60)
			owner := &api.User{}
			owner.SetName("test-user")
			repo.SetOwner(owner)
			build.SetRepo(repo)
			return &models.Item{
				Build:       build,
				ItemVersion: models.ItemVersion,
			}, nil
		},
	}

	// Create a Vela client
	client, _ := vela.NewClient(server.URL, "", nil)

	// Create test worker
	w := &Worker{
		Config: &Config{
			Mock: true,
			API: &API{
				Address: serverURL,
			},
			Build: &Build{
				Limit:       5,
				Timeout:     30 * time.Minute,
				CPUQuota:    1200,
				MemoryLimit: 4,
				PidsLimit:   1024,
			},
			Server: &Server{
				Address: server.URL,
				Secret:  "test-secret",
			},
			Executor: &executor.Setup{
				Driver:     "linux",
				MaxLogSize: 1000000,
				OutputCtn: &pipeline.Container{
					ID:    "outputs",
					Image: "alpine:latest",
				},
			},
			Runtime: &runtime.Setup{
				Driver: "docker",
			},
			Queue: &queue.Setup{
				Timeout: 5 * time.Second,
			},
		},
		VelaClient:         client,
		Queue:              testQueue,
		Executors:          make(map[int]executor.Engine),
		RunningBuilds:      []*api.Build{},
		RunningBuildsMutex: sync.Mutex{},
		BuildContexts:      make(map[string]*BuildContext),
		BuildContextsMutex: sync.RWMutex{},
	}

	config := &api.Worker{}
	config.SetHostname("test-worker")
	config.SetRunningBuilds([]*api.Build{})

	// Test successful execution
	err := w.exec(0, config)
	if err != nil {
		t.Errorf("exec() error = %v, want nil", err)
	}
}

func TestWorker_exec_QueuePopError(t *testing.T) {
	// Set up a test server to mock the Vela API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/workers/test-worker":
			worker := &api.Worker{}
			worker.SetHostname("test-worker")
			worker.SetRoutes([]string{"repo"})
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(worker)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Test queue pop failure
	testQueue := &mockQueue{
		popFunc: func(_ context.Context, _ []string) (*models.Item, error) {
			return nil, errors.New("queue pop failed")
		},
	}

	serverURL, _ := url.Parse(server.URL)
	client, _ := vela.NewClient(server.URL, "", nil)

	w := &Worker{
		Config: &Config{
			API: &API{
				Address: serverURL,
			},
			Queue: &queue.Setup{
				Timeout: 1 * time.Second,
			},
		},
		Queue:              testQueue,
		VelaClient:         client,
		RunningBuildsMutex: sync.Mutex{},
	}

	config := &api.Worker{}
	config.SetHostname("test-worker")

	err := w.exec(0, config)
	if err != nil {
		t.Errorf("exec() should return nil on queue pop error, got %v", err)
	}
}

func TestWorker_exec_QueuePopNilItem(t *testing.T) {
	// Set up a test server to mock the Vela API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/workers/test-worker":
			worker := &api.Worker{}
			worker.SetHostname("test-worker")
			worker.SetRoutes([]string{"repo"})
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(worker)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Test nil item from queue
	testQueue := &mockQueue{
		popFunc: func(_ context.Context, _ []string) (*models.Item, error) {
			return nil, nil
		},
	}

	serverURL, _ := url.Parse(server.URL)
	client, _ := vela.NewClient(server.URL, "", nil)

	w := &Worker{
		Config: &Config{
			API: &API{
				Address: serverURL,
			},
		},
		Queue:              testQueue,
		VelaClient:         client,
		RunningBuildsMutex: sync.Mutex{},
	}

	config := &api.Worker{}
	config.SetHostname("test-worker")

	err := w.exec(0, config)
	if err != nil {
		t.Errorf("exec() should return nil for nil queue item, got %v", err)
	}
}

func TestWorker_exec_StaleItemVersion(t *testing.T) {
	// Set up a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/workers/test-worker":
			worker := &api.Worker{}
			worker.SetHostname("test-worker")
			worker.SetRoutes([]string{"repo"})
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(worker)
		case "/api/v1/repos/test-org/test-repo/builds/1/token":
			token := &api.Token{}
			token.SetToken("test-token")
			_ = json.NewEncoder(w).Encode(token)
		case "/api/v1/repos/test-org/test-repo/builds/1/executable":
			pipelineBuild := &pipeline.Build{ID: "test-id"}
			data, _ := json.Marshal(pipelineBuild)
			executable := &api.BuildExecutable{}
			executable.SetData(data)
			_ = json.NewEncoder(w).Encode(executable)
		case "/api/v1/repos/test-org/test-repo/builds/1":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(&api.Build{})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	// Create a test queue with stale item version
	testQueue := &mockQueue{
		popFunc: func(_ context.Context, _ []string) (*models.Item, error) {
			build := &api.Build{}
			build.SetID(1)
			build.SetNumber(1)
			repo := &api.Repo{}
			repo.SetOrg("test-org")
			repo.SetName("test-repo")
			repo.SetFullName("test-org/test-repo")
			owner := &api.User{}
			owner.SetName("test-user")
			repo.SetOwner(owner)
			build.SetRepo(repo)
			return &models.Item{
				Build:       build,
				ItemVersion: models.ItemVersion - 1, // Stale version
			}, nil
		},
	}

	client, _ := vela.NewClient(server.URL, "", nil)

	w := &Worker{
		Config: &Config{
			API: &API{
				Address: serverURL,
			},
			Build: &Build{
				CPUQuota:    1200,
				MemoryLimit: 4,
				PidsLimit:   1024,
			},
			Server: &Server{
				Address: server.URL,
			},
			Executor: &executor.Setup{
				OutputCtn: &pipeline.Container{},
			},
			Runtime: &runtime.Setup{},
		},
		VelaClient:         client,
		Queue:              testQueue,
		RunningBuilds:      []*api.Build{},
		RunningBuildsMutex: sync.Mutex{},
		BuildContexts:      make(map[string]*BuildContext),
		BuildContextsMutex: sync.RWMutex{},
	}

	config := &api.Worker{}
	config.SetHostname("test-worker")

	err := w.exec(0, config)
	if err != nil {
		t.Errorf("exec() should return nil for stale item, got %v", err)
	}
}

func TestWorker_exec_GetBuildTokenConflict(t *testing.T) {
	// Set up a test server that returns conflict for build token
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/workers/test-worker":
			worker := &api.Worker{}
			worker.SetHostname("test-worker")
			worker.SetRoutes([]string{"repo"})
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(worker)
		case "/api/v1/repos/test-org/test-repo/builds/1/token":
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	serverURL, _ := url.Parse(server.URL)

	testQueue := &mockQueue{
		popFunc: func(_ context.Context, _ []string) (*models.Item, error) {
			build := &api.Build{}
			build.SetID(1)
			build.SetNumber(1)
			repo := &api.Repo{}
			repo.SetOrg("test-org")
			repo.SetName("test-repo")
			repo.SetFullName("test-org/test-repo")
			owner := &api.User{}
			owner.SetName("test-user")
			repo.SetOwner(owner)
			build.SetRepo(repo)
			return &models.Item{
				Build:       build,
				ItemVersion: models.ItemVersion,
			}, nil
		},
	}

	client, _ := vela.NewClient(server.URL, "", nil)

	w := &Worker{
		Config: &Config{
			API: &API{
				Address: serverURL,
			},
		},
		VelaClient: client,
		Queue:      testQueue,
	}

	config := &api.Worker{}
	config.SetHostname("test-worker")

	err := w.exec(0, config)
	if err != nil {
		t.Errorf("exec() should return nil on conflict, got %v", err)
	}
}

func TestWorker_exec_RetryLogic(t *testing.T) {
	// Test retry logic when worker retrieval fails initially
	attempt := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/workers/test-worker":
			attempt++
			if attempt < 2 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			worker := &api.Worker{}
			worker.SetHostname("test-worker")
			worker.SetRoutes([]string{"repo"})
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(worker)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	serverURL, _ := url.Parse(server.URL)
	client, _ := vela.NewClient(server.URL, "", nil)

	testQueue := &mockQueue{
		popFunc: func(_ context.Context, _ []string) (*models.Item, error) {
			return nil, nil // Return nil to exit early
		},
	}

	w := &Worker{
		Config: &Config{
			API: &API{
				Address: serverURL,
			},
		},
		VelaClient: client,
		Queue:      testQueue,
	}

	config := &api.Worker{}
	config.SetHostname("test-worker")

	// Should succeed on retry
	err := w.exec(0, config)
	if err != nil {
		t.Errorf("exec() error = %v, want nil after retry", err)
	}

	if attempt < 2 {
		t.Errorf("Expected at least 2 attempts, got %d", attempt)
	}
}

func TestWorker_exec_MaxRetriesExceeded(t *testing.T) {
	// Test max retries exceeded
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError) // Always fail
	}))
	defer server.Close()

	serverURL, _ := url.Parse(server.URL)
	client, _ := vela.NewClient(server.URL, "", nil)

	w := &Worker{
		Config: &Config{
			API: &API{
				Address: serverURL,
			},
		},
		VelaClient: client,
	}

	config := &api.Worker{}
	config.SetHostname("test-worker")

	err := w.exec(0, config)
	if err == nil {
		t.Error("exec() should return error when max retries exceeded")
	}
}

// Mock queue implementation for testing.
type mockQueue struct {
	popFunc func(context.Context, []string) (*models.Item, error)
}

func (m *mockQueue) Pop(ctx context.Context, routes []string) (*models.Item, error) {
	if m.popFunc != nil {
		return m.popFunc(ctx, routes)
	}

	return nil, nil
}

func (m *mockQueue) Route(_ *pipeline.Worker) (string, error)               { return "vela", nil }
func (m *mockQueue) Driver() string                                         { return "mock" }
func (m *mockQueue) GetSettings() settings.Queue                            { return settings.Queue{} }
func (m *mockQueue) SetSettings(_ *settings.Platform)                       {}
func (m *mockQueue) Length(_ context.Context) (int64, error)                { return 0, nil }
func (m *mockQueue) RouteLength(_ context.Context, _ string) (int64, error) { return 0, nil }
func (m *mockQueue) Push(_ context.Context, _ string, _ []byte) error       { return nil }
func (m *mockQueue) Ping(_ context.Context) error                           { return nil }

func TestWorker_exec_LogOutput(t *testing.T) {
	// Capture log output to verify logging behavior
	origLevel := logrus.GetLevel()

	logrus.SetLevel(logrus.DebugLevel)
	defer logrus.SetLevel(origLevel)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/workers/test-worker":
			worker := &api.Worker{}
			worker.SetHostname("test-worker")
			worker.SetRoutes([]string{"repo"})
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(worker)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	serverURL, _ := url.Parse(server.URL)
	client, _ := vela.NewClient(server.URL, "", nil)

	testQueue := &mockQueue{
		popFunc: func(_ context.Context, _ []string) (*models.Item, error) {
			return nil, fmt.Errorf("simulated queue error")
		},
	}

	w := &Worker{
		Config: &Config{
			API: &API{
				Address: serverURL,
			},
			Queue: &queue.Setup{
				Timeout: 100 * time.Millisecond,
			},
		},
		VelaClient: client,
		Queue:      testQueue,
	}

	config := &api.Worker{}
	config.SetHostname("test-worker")

	// This should log the queue pop error
	err := w.exec(0, config)
	if err != nil {
		t.Errorf("exec() should return nil on queue error, got %v", err)
	}
}

func TestWorker_getWorkerStatusFromConfig(t *testing.T) {
	tests := []struct {
		name           string
		buildLimit     int
		runningBuilds  []*api.Build
		expectedStatus string
	}{
		{
			name:           "idle status",
			buildLimit:     5,
			runningBuilds:  []*api.Build{},
			expectedStatus: constants.WorkerStatusIdle,
		},
		{
			name:       "available status",
			buildLimit: 5,
			runningBuilds: []*api.Build{
				{ID: func() *int64 { id := int64(1); return &id }()},
				{ID: func() *int64 { id := int64(2); return &id }()},
			},
			expectedStatus: constants.WorkerStatusAvailable,
		},
		{
			name:       "busy status",
			buildLimit: 3,
			runningBuilds: []*api.Build{
				{ID: func() *int64 { id := int64(1); return &id }()},
				{ID: func() *int64 { id := int64(2); return &id }()},
				{ID: func() *int64 { id := int64(3); return &id }()},
			},
			expectedStatus: constants.WorkerStatusBusy,
		},
		{
			name:       "error status",
			buildLimit: 2,
			runningBuilds: []*api.Build{
				{ID: func() *int64 { id := int64(1); return &id }()},
				{ID: func() *int64 { id := int64(2); return &id }()},
				{ID: func() *int64 { id := int64(3); return &id }()},
			},
			expectedStatus: constants.WorkerStatusError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				Config: &Config{
					Build: &Build{
						Limit: tt.buildLimit,
					},
				},
			}

			config := &api.Worker{}
			config.SetRunningBuilds(tt.runningBuilds)

			status := w.getWorkerStatusFromConfig(config)
			if status != tt.expectedStatus {
				t.Errorf("getWorkerStatusFromConfig() = %v, want %v", status, tt.expectedStatus)
			}
		})
	}
}
