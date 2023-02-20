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

package worker

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/neuralnorthwest/mu/logging"
	mock_logging "github.com/neuralnorthwest/mu/logging/mock"
	"github.com/stretchr/testify/assert"
)

// adapterTestServer is a test server for the adapter.
type adapterTestServer struct {
	stopchan         chan struct{}
	runErr           error       // if set, this error will be immediately returned from run
	stoppingErr      error       // if set, this error will be returned from run when it is stopped
	stopErr          error       // if set, this error will be returned from stop
	workerDoesntStop bool        // if true, the worker will not stop when the context is canceled
	stopped          atomic.Bool // if true, the worker has been stopped
}

// newAdapterTestServer creates a new test server for the adapter.
func newAdapterTestServer() *adapterTestServer {
	return &adapterTestServer{
		stopchan: make(chan struct{}, 1),
	}
}

// run is the function that runs the resource. It should be blocking,
// and return an error if it fails.
func (s *adapterTestServer) run(logger logging.Logger) error {
	if s.runErr != nil {
		return s.runErr
	}
	<-s.stopchan
	if s.stoppingErr != nil {
		return s.stoppingErr
	}
	return nil
}

// stop is the function that stops the resource.
func (s *adapterTestServer) stop(logger logging.Logger) error {
	if !s.workerDoesntStop {
		s.stopchan <- struct{}{}
		s.stopped.Store(true)
	}
	return s.stopErr
}

// Test_Adapter_Case is a test case for the adapter.
type Test_Adapter_Case struct {
	// name is the name of the test case.
	name string
	// runErr is the error that will be returned from run.
	runErr error
	// stoppingErr is the error that will be returned from run when it is stopped.
	stoppingErr error
	// stopErr is the error that will be returned from stop.
	stopErr error
	// workerDoesntStop is true if the worker should not stop when the context is
	// canceled.
	workerDoesntStop bool
	// expectedErr is the expected error.
	expectedErr error
	// noLogger is true if the logger should be nil.
	noLogger bool
}

// Test_Adapter_Cases is a list of test cases for the adapter.
var Test_Adapter_Cases = []Test_Adapter_Case{
	{
		name: "success",
	},
	{
		name:             "run error",
		runErr:           errors.New("run error"),
		expectedErr:      errors.New("run error"),
		workerDoesntStop: true,
	},
	{
		name:        "stop error",
		stopErr:     errors.New("stop error"),
		expectedErr: errors.New("stop error"),
	},
	{
		name:        "stop error, no logger",
		stopErr:     errors.New("stop error"),
		expectedErr: errors.New("stop error"),
		noLogger:    true,
	},
	{
		name:        "stopping error",
		stoppingErr: errors.New("stopping error"),
		expectedErr: errors.New("stopping error"),
	},
	{
		name:             "worker doesn't stop",
		stopErr:          errors.New("stop error"),
		expectedErr:      errors.New("stop error"),
		workerDoesntStop: true,
	},
}

// Test_Adapter tests the adapter.
func Test_Adapter(t *testing.T) {
	t.Parallel()
	for _, tc := range Test_Adapter_Cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			logger := mock_logging.NewMockLogger(mockCtrl)
			if tc.stopErr != nil && !tc.noLogger {
				logger.EXPECT().Errorw("error stopping worker", "err", tc.stopErr)
			}
			server := newAdapterTestServer()
			server.runErr = tc.runErr
			server.stoppingErr = tc.stoppingErr
			server.stopErr = tc.stopErr
			server.workerDoesntStop = tc.workerDoesntStop
			worker := Adapter(server.run, server.stop)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			errchan := make(chan error, 1)
			go func() {
				if tc.noLogger {
					errchan <- worker.Run(ctx, nil)
					return
				}
				errchan <- worker.Run(ctx, logger)
			}()
			if tc.runErr == nil {
				cancel()
			}
			if tc.runErr == nil && tc.workerDoesntStop {
				// If the worker doesn't stop, we need to stop it manually.
				server.stopchan <- struct{}{}
			}
			err := <-errchan
			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedErr.Error())
			}
			assert.NotEqual(t, server.stopped.Load(), tc.workerDoesntStop)
		})
	}
}
