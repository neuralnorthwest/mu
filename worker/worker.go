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
	"fmt"
	"sync"

	"github.com/neuralnorthwest/mu/logging"
	"github.com/neuralnorthwest/mu/status"
	"golang.org/x/sync/errgroup"
)

// Worker is an interface for a worker.
type Worker interface {
	// Run runs the worker in the given errgroup. The worker will stop when the
	// context is canceled.
	Run(ctx context.Context, logger logging.Logger) error
}

// Group is the interface for a worker group.
type Group interface {
	// Add adds a worker to the worker group. The worker will be started when the
	// worker group is started. If the group has already been started, the worker
	// will be started immediately.
	Add(name string, worker Worker) error
	// Run runs the worker group. This will start all the workers in the
	// worker group. This will block until the context is canceled or a worker
	// returns an error.
	Run(ctx context.Context, logger logging.Logger) error
	// Start starts the worker group. This will start all the workers in the
	// worker group. This will not block. To wait for the workers to stop, call
	// Wait after canceling the context.
	Start(ctx context.Context, logger logging.Logger) error
	// Wait waits for the worker group to stop.
	Wait() error
}

// group is a group of workers.
type group struct {
	// lock is the lock for the worker group.
	lock sync.Mutex
	// ctx is the context.
	ctx context.Context
	// workers is a map of workers.
	workers map[string]Worker
	// eg is the errgroup.
	eg *errgroup.Group
	// started is true if the worker group has been started.
	started bool
	// logger is the logger for the worker group.
	logger logging.Logger
}

// NewGroup creates a new worker group.
func NewGroup() Group {
	return &group{
		workers: make(map[string]Worker),
	}
}

// Add adds a worker to the worker group. The worker will be started when the
// worker group is started. If the group has already been started, the worker
// will be started immediately.
func (g *group) Add(name string, worker Worker) error {
	g.lock.Lock()
	defer g.lock.Unlock()
	if _, ok := g.workers[name]; ok {
		return fmt.Errorf("%w: %s", status.ErrAlreadyExists, name)
	}
	g.workers[name] = worker
	if g.started {
		return g.startWorker(name, worker)
	}
	return nil
}

// Run runs the worker group. This will start all the workers in the
// worker group. This will block until the context is canceled.
func (g *group) Run(ctx context.Context, logger logging.Logger) error {
	if err := g.Start(ctx, logger); err != nil {
		return err
	}
	return g.Wait()
}

// Start starts the worker group. This will start all the workers in the
// worker group. This will not block. To wait for the workers to stop, call
// Wait after canceling the context.
func (g *group) Start(ctx context.Context, logger logging.Logger) error {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.started {
		return status.ErrAlreadyStarted
	}
	g.ctx = ctx
	g.eg, g.ctx = errgroup.WithContext(ctx)
	g.started = true
	g.logger = logger
	for name, worker := range g.workers {
		if err := g.startWorker(name, worker); err != nil {
			return err
		}
	}
	return nil
}

// Wait waits for the worker group to stop.
func (g *group) Wait() error {
	g.lock.Lock()
	if !g.started {
		g.lock.Unlock()
		return status.ErrNotStarted
	}
	g.lock.Unlock()
	return g.eg.Wait()
}

// startWorker starts a worker.
func (g *group) startWorker(name string, worker Worker) error {
	g.eg.Go(func() error {
		return worker.Run(g.ctx, g.logger.With("worker", name))
	})
	return nil
}
