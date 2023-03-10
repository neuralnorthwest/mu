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
	"testing"

	"github.com/neuralnorthwest/mu/logging"
	"github.com/neuralnorthwest/mu/worker"
)

// testTraceWorker is a worker that records that it was run.
type testTraceWorker struct {
	ran *bool
}

// Run implements the worker interface.
func (w *testTraceWorker) Run(ctx context.Context, logger logging.Logger) error {
	*w.ran = true
	return nil
}

// Test_Main tests that the main command is created and executed.
func Test_Main(t *testing.T) {
	t.Parallel()
	workerRan := false
	svc, err := New("test-service")
	if err != nil {
		t.Fatalf("New returned an error: %v", err)
	}
	svc.SetupWorkers(func(group worker.Group) error {
		return group.Add("worker", &testTraceWorker{ran: &workerRan})
	})
	if err := svc.Main(); err != nil {
		t.Fatalf("Main returned an error: %v", err)
	}
	if !workerRan {
		t.Fatal("worker was not run")
	}
}
