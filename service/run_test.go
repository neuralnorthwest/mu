// Copyright 2023 Scott M. Long
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"
	"fmt"
	"syscall"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/neuralnorthwest/mu/config"
	"github.com/neuralnorthwest/mu/http"
	"github.com/neuralnorthwest/mu/logging"
	mock_logging "github.com/neuralnorthwest/mu/logging/mock"
	"github.com/neuralnorthwest/mu/status"
	"github.com/neuralnorthwest/mu/worker"
)

// Test_run tests that the service runs. This is a basic test with no
// workers.
func Test_run(t *testing.T) {
	t.Parallel()
	svc, err := New("test-service")
	if err != nil {
		t.Fatalf("New returned an error: %v", err)
	}
	if err := svc.Run(); err != nil {
		t.Errorf("Run returned an error: %v", err)
	}
}

// Test_run_MockMode tests that the service runs in mock mode and logs a
// message indicating that it is running in mock mode.
func Test_run_MockMode(t *testing.T) {
	t.Parallel()
	mc := gomock.NewController(t)
	logger := mock_logging.NewMockLogger(mc)
	logger.EXPECT().Info("running in mock mode")
	svc, err := New("test-service", WithLogger(func() (logging.Logger, error) {
		return logger, nil
	}), WithMockMode())
	if err != nil {
		t.Fatalf("New returned an error: %v", err)
	}
	if err := svc.Run(); err != nil {
		t.Errorf("Run returned an error: %v", err)
	}
}

// Test_run_Hooks_Case is a test case for Test_run_Hooks. It indicates which,
// if any, hooks should return errors.
type Test_run_Hooks_Case struct {
	name                   string
	cleanupErr             error
	configSetupErr         error
	prerunErr              error
	setupWorkersErr        error
	setupHTTPErr           error
	expectedErr            string
	cleanupWasInvoked      bool
	configSetupWasInvoked  bool
	prerunWasInvoked       bool
	setupWorkersWasInvoked bool
	setupHTTPWasInvoked    bool
}

// Test_run_Hooks tests that the service runs with hooks.
func Test_run_Hooks(t *testing.T) {
	t.Parallel()
	testCases := []Test_run_Hooks_Case{
		{
			name:                   "no errors",
			cleanupErr:             nil,
			configSetupErr:         nil,
			prerunErr:              nil,
			setupWorkersErr:        nil,
			setupHTTPErr:           nil,
			cleanupWasInvoked:      true,
			configSetupWasInvoked:  true,
			prerunWasInvoked:       true,
			setupWorkersWasInvoked: true,
			setupHTTPWasInvoked:    true,
		},
		{
			name:                   "cleanup error",
			cleanupErr:             fmt.Errorf("cleanup error"),
			configSetupErr:         nil,
			prerunErr:              nil,
			setupWorkersErr:        nil,
			setupHTTPErr:           nil,
			expectedErr:            "cleanup error",
			cleanupWasInvoked:      true,
			configSetupWasInvoked:  true,
			prerunWasInvoked:       true,
			setupWorkersWasInvoked: true,
			setupHTTPWasInvoked:    true,
		},
		{
			name:                   "config setup error",
			cleanupErr:             nil,
			configSetupErr:         fmt.Errorf("config setup error"),
			prerunErr:              nil,
			setupWorkersErr:        nil,
			setupHTTPErr:           nil,
			expectedErr:            "config setup error",
			cleanupWasInvoked:      false,
			configSetupWasInvoked:  true,
			prerunWasInvoked:       false,
			setupWorkersWasInvoked: false,
			setupHTTPWasInvoked:    false,
		},
		{
			name:                   "prerun error",
			cleanupErr:             nil,
			configSetupErr:         nil,
			prerunErr:              fmt.Errorf("prerun error"),
			setupWorkersErr:        nil,
			setupHTTPErr:           nil,
			expectedErr:            "prerun error",
			cleanupWasInvoked:      true,
			configSetupWasInvoked:  true,
			prerunWasInvoked:       true,
			setupWorkersWasInvoked: true,
			setupHTTPWasInvoked:    true,
		},
		{
			name:                   "setup workers error",
			cleanupErr:             nil,
			configSetupErr:         nil,
			prerunErr:              nil,
			setupWorkersErr:        fmt.Errorf("setup workers error"),
			setupHTTPErr:           nil,
			expectedErr:            "setup workers error",
			cleanupWasInvoked:      false,
			configSetupWasInvoked:  true,
			prerunWasInvoked:       false,
			setupWorkersWasInvoked: true,
			setupHTTPWasInvoked:    false,
		},
		{
			name:                   "setup HTTP error",
			cleanupErr:             nil,
			configSetupErr:         nil,
			prerunErr:              nil,
			setupWorkersErr:        nil,
			setupHTTPErr:           fmt.Errorf("setup HTTP error"),
			expectedErr:            "setup HTTP error",
			cleanupWasInvoked:      true,
			configSetupWasInvoked:  true,
			prerunWasInvoked:       false,
			setupWorkersWasInvoked: true,
			setupHTTPWasInvoked:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := New("test-service")
			if err != nil {
				t.Fatalf("New returned an error: %v", err)
			}
			cleanupWasInvoked := false
			configSetupWasInvoked := false
			prerunWasInvoked := false
			setupWorkersWasInvoked := false
			setupHTTPWasInvoked := false
			svc.Cleanup(func() error {
				cleanupWasInvoked = true
				return tc.cleanupErr
			})
			svc.ConfigSetup(func(config.Config) error {
				configSetupWasInvoked = true
				return tc.configSetupErr
			})
			svc.PreRun(func() error {
				prerunWasInvoked = true
				// Cancel the context so that Run returns.
				svc.Cancel()
				return tc.prerunErr
			})
			svc.SetupWorkers(func(worker.Group) error {
				setupWorkersWasInvoked = true
				return tc.setupWorkersErr
			})
			svc.SetupHTTP(func(*http.Server) error {
				setupHTTPWasInvoked = true
				return tc.setupHTTPErr
			})
			err = svc.Run()
			if tc.expectedErr != "" {
				if err == nil {
					t.Errorf("Run did not return an error")
				} else if err.Error() != tc.expectedErr {
					t.Errorf("Run returned an unexpected error: %v but expected %v", err, tc.expectedErr)
				}
			} else if err != nil {
				t.Errorf("Run returned an error: %v", err)
			}
			if cleanupWasInvoked != tc.cleanupWasInvoked {
				t.Errorf("cleanup hook was invoked: %t", cleanupWasInvoked)
			}
			if configSetupWasInvoked != tc.configSetupWasInvoked {
				t.Errorf("config setup hook was invoked: %t", configSetupWasInvoked)
			}
			if prerunWasInvoked != tc.prerunWasInvoked {
				t.Errorf("prerun hook was invoked: %t", prerunWasInvoked)
			}
			if setupWorkersWasInvoked != tc.setupWorkersWasInvoked {
				t.Errorf("setup hook was invoked: %t", setupWorkersWasInvoked)
			}
			if setupHTTPWasInvoked != tc.setupHTTPWasInvoked {
				t.Errorf("setup HTTP hook was invoked: %t", setupHTTPWasInvoked)
			}
		})
	}
}

// Test_run_Workers_Case is a test case for Test_run_Workers. It specifies
// which workers to start and what errors they should return.
type Test_run_Workers_Case struct {
	name       string
	workers    []worker.Worker
	workerErrs []error
}

// testWorker is a Worker that immediately returns the given error.
type testWorker struct {
	t   *testing.T
	err error
}

// Run implements the Worker interface.
func (w *testWorker) Run(ctx context.Context, logger logging.Logger) error {
	return w.err
}

// newTestWorker returns a new testWorker.
func newTestWorker(t *testing.T, err error) worker.Worker {
	t.Helper()
	return &testWorker{
		t:   t,
		err: err,
	}
}

// Test_run_Workers tests that the service runs with workers.
func Test_run_Workers(t *testing.T) {
	t.Parallel()
	testCases := []Test_run_Workers_Case{
		{
			name:       "no workers",
			workers:    nil,
			workerErrs: nil,
		},
		{
			name:       "one worker",
			workers:    []worker.Worker{newTestWorker(t, nil)},
			workerErrs: []error{nil},
		},
		{
			name:       "two workers",
			workers:    []worker.Worker{newTestWorker(t, nil), newTestWorker(t, nil)},
			workerErrs: []error{nil, nil},
		},
		{
			name:       "one worker with one error",
			workers:    []worker.Worker{newTestWorker(t, status.ErrInvalidArgument)},
			workerErrs: []error{status.ErrInvalidArgument},
		},
		{
			name:       "two workers with one error",
			workers:    []worker.Worker{newTestWorker(t, nil), newTestWorker(t, status.ErrInvalidArgument)},
			workerErrs: []error{nil, status.ErrInvalidArgument},
		},
		{
			name:       "two workers with two error",
			workers:    []worker.Worker{newTestWorker(t, status.ErrInvalidArgument), newTestWorker(t, status.ErrInvalidArgument)},
			workerErrs: []error{status.ErrInvalidArgument, status.ErrInvalidArgument},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			svc, err := New("test-service")
			if err != nil {
				t.Fatalf("New returned an error: %v", err)
			}
			svc.SetupWorkers(func(group worker.Group) error {
				for i, w := range testCase.workers {
					if err := group.Add(fmt.Sprintf("worker-%d", i), w); err != nil {
						return err
					}
				}
				return nil
			})
			err = svc.Run()
			matchedErr := false
			expectAnError := false
			for _, workerErr := range testCase.workerErrs {
				if err == workerErr {
					matchedErr = true
				}
				if workerErr != nil {
					expectAnError = true
				}
			}
			if expectAnError && !matchedErr {
				t.Errorf("Run did not return an error")
			}
			if !expectAnError && err != nil {
				t.Errorf("Run returned an error: %v", err)
			}
		})
	}
}

// testWaitWorker is a worker that waits for context cancellation.
type testWaitWorker struct {
	t *testing.T
}

// Run implements the Worker interface.
func (w *testWaitWorker) Run(ctx context.Context, logger logging.Logger) error {
	<-ctx.Done()
	return nil
}

// newTestWaitWorker returns a new testWaitWorker.
func newTestWaitWorker(t *testing.T) worker.Worker {
	t.Helper()
	return &testWaitWorker{
		t: t,
	}
}

// Test_run_interrupt tests that the service terminates when interrupted by
// a signal.
func Test_run_interrupt(t *testing.T) {
	t.Parallel()
	mc := gomock.NewController(t)
	logger := mock_logging.NewMockLogger(mc)
	logger.EXPECT().With("worker", "wait-worker").Return(logger)
	logger.EXPECT().Infow("received interrupt signal", "signal", syscall.SIGINT)
	svc, err := New("test-service", WithLogger(func() (logging.Logger, error) { return logger, nil }))
	if err != nil {
		t.Fatalf("New returned an error: %v", err)
	}
	svc.SetupWorkers(func(group worker.Group) error {
		return group.Add("wait-worker", newTestWaitWorker(t))
	})
	svc.sigChan <- syscall.SIGINT
	err = svc.Run()
	if err != nil {
		t.Errorf("Run returned an error: %v", err)
	}
}
