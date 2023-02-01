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
	svc.RegisterSetup(func(group worker.Group) error {
		group.Add("worker", &testTraceWorker{ran: &workerRan})
		return nil
	})
	if err := svc.Main(); err != nil {
		t.Fatalf("Main returned an error: %v", err)
	}
	if !workerRan {
		t.Fatal("worker was not run")
	}
}
