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

	"github.com/neuralnorthwest/mu/logging"
)

// adapter struct
type adapter struct {
	// run is the function that runs the resource. It should be blocking,
	// and return an error if it fails.
	run func(logger logging.Logger) error
	// stop is the function that stops the resource.
	stop func(logger logging.Logger) error
}

// Adapter can be used to convert any kind of startable/stoppable resource
// into a worker. Simply provide a function that runs the resource, and a
// function that stops the resource. The worker will stop when the context is
// canceled or an error occurs. If run returns an error, that error will be
// returned. If run does not return an error, but stop returns an error, that
// error will be returned.
//
// For example, this can be used to convert the default HTTP server into a
// worker:
//
//	httpWorker := worker.Adapter(
//		func(logging.Logger) error {
//			return http.ListenAndServe(":8080", nil)
//		},
//		func(logging.Logger) error {
//			return http.Shutdown(context.Background())
//		},
//	)
//
// The httpWorker can then be used in a worker group:
//
//	workerGroup := worker.NewGroup()
//	workerGroup.Add("http", httpWorker)
//	workerGroup.Run(ctx, logger)
func Adapter(run, stop func(logger logging.Logger) error) Worker {
	return &adapter{
		run:  run,
		stop: stop,
	}
}

// Run runs the worker. The worker will stop when the context is canceled
// or an error occurs. If run returns an error, that error will be returned.
// If run does not return an error, but stop returns an error, that error
// will be returned.
func (a *adapter) Run(ctx context.Context, logger logging.Logger) error {
	result := make(chan error, 1)
	go func() {
		result <- a.run(logger)
	}()
	select {
	case <-ctx.Done():
		stoperr := a.stop(logger)
		if stoperr != nil {
			if logger != nil {
				// Log this error, to aid in debugging if the worker does not
				// stop.
				logger.Errorw("error stopping worker", "err", stoperr)
			}
		}
		err := <-result
		if err != nil {
			return err
		}
		return stoperr
	case err := <-result:
		return err
	}
}
