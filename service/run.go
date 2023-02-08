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
	"os/signal"
	"syscall"

	"github.com/neuralnorthwest/mu/logging"
	"github.com/neuralnorthwest/mu/worker"
)

// Run runs the service.
func (s *Service) Run() (status error) {
	if s.MockMode() {
		s.logger.Info("running in mock mode")
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	defer s.cancel()
	if err := s.invokeSetupConfig(s.config); err != nil {
		return err
	}
	workerGroup := worker.NewGroup()
	if err := s.invokeSetupWorkers(workerGroup); err != nil {
		return err
	}
	defer func() {
		err := s.invokeCleanup()
		if status == nil {
			status = err
		}
	}()
	if httpServer, err := s.invokeSetupHTTP(); err != nil {
		return err
	} else if httpServer != nil {
		if err := workerGroup.Add("http_server", httpServer); err != nil {
			return err
		}
	}
	s.startInterruptListener(s.ctx, s.logger, s.cancel)
	if err := s.invokePreRun(); err != nil {
		return err
	}
	werr := workerGroup.Run(s.ctx, s.logger)
	if werr != nil {
		return werr
	}
	return
}

// startInterruptListener starts the interrupt listener. This registers a
// listener for SIGINT and SIGTERM signals, and starts a goroutine that
// cancels the context when a signal is received.
func (s *Service) startInterruptListener(ctx context.Context, logger logging.Logger, cancel context.CancelFunc) {
	signal.Notify(s.sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-s.sigChan
		logger.Infow("received interrupt signal", "signal", sig)
		cancel()
	}()
}
