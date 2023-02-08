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
	"os"

	"github.com/neuralnorthwest/mu/config"
	"github.com/neuralnorthwest/mu/logging"
)

// Service represents a service.
type Service struct {
	// name is the name of the service.
	name string
	// version is the version of the service.
	version string
	// hooks are the hooks for the service.
	Hooks
	// logger is the logger for the service.
	logger logging.Logger
	// ctx holds the context for the service.
	ctx context.Context
	// cancel cancels the context for the service.
	cancel context.CancelFunc
	// config is the config for the service.
	config config.Config
	// mockMode is true if the service is in mock mode.
	mockMode bool
	// newLogger is a func that returns a new logger.
	newLogger func() (logging.Logger, error)
	// sigChan is the channel for signals.
	sigChan chan os.Signal
}

// New returns a new service.
func New(name string, opts ...Option) (*Service, error) {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		name:     name,
		version:  "v0.0.0",
		Hooks:    &hookstruct{},
		ctx:      ctx,
		cancel:   cancel,
		config:   config.New(),
		mockMode: false,
		newLogger: func() (logging.Logger, error) {
			return logging.New()
		},
		sigChan: make(chan os.Signal, 1),
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	logger, err := s.newLogger()
	if err != nil {
		return nil, err
	}
	s.logger = logger
	return s, nil
}

// Name returns the name of the service.
func (s *Service) Name() string {
	return s.name
}

// Version returns the version of the service.
func (s *Service) Version() string {
	return s.version
}

// Logger returns the logger for the service.
func (s *Service) Logger() logging.Logger {
	return s.logger
}

// Context returns the context for the service.
func (s *Service) Context() context.Context {
	return s.ctx
}

// Cancel cancels the context for the service.
func (s *Service) Cancel() {
	s.cancel()
}

// Config returns the config for the service.
func (s *Service) Config() config.Config {
	return s.config
}

// MockMode returns true if the service is in mock mode.
func (s *Service) MockMode() bool {
	return s.mockMode
}
