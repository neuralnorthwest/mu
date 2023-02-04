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
	"github.com/neuralnorthwest/mu/config"
	"github.com/neuralnorthwest/mu/http"
	"github.com/neuralnorthwest/mu/worker"
)

// CleanupFunc is a function that cleans up a service.
type CleanupFunc func() error

// ConfigSetupFunc is a function that sets up a service configuration.
type ConfigSetupFunc func(c config.Config) error

// SetupFunc is a function that sets up a service.
type SetupFunc func(workerGroup worker.Group) error

// SetupHTTPFunc is a function that sets up a service HTTP server.
type SetupHTTPFunc func(server *http.Server) error

type Hooks interface {
	Cleanup(f CleanupFunc)
	ConfigSetup(f ConfigSetupFunc)
	Setup(f SetupFunc)
	SetupHTTP(f SetupHTTPFunc)

	invokeCleanup() error
	invokeConfigSetup(c config.Config) error
	invokeSetup(workerGroup worker.Group) error
	invokeSetupHTTP() (*http.Server, error)
}

// hookstruct holds the hooks for a service.
type hookstruct struct {
	// cleanup is the cleanup hook.
	cleanup CleanupFunc
	// configSetup is the setup configuration hook.
	configSetup ConfigSetupFunc
	// setup is the setup hook.
	setup SetupFunc
	// setupHTTP is the setup HTTP hook.
	setupHTTP SetupHTTPFunc
	// httpNewServer is a function that creates a new HTTP server. This is used
	// for testing.
	httpNewServer func() (*http.Server, error)
}

var _ Hooks = (*hookstruct)(nil)

// Cleanup registers a cleanup hook.
func (h *hookstruct) Cleanup(f CleanupFunc) {
	h.cleanup = f
}

// SetupConfig registers a setup configuration hook.
func (h *hookstruct) ConfigSetup(f ConfigSetupFunc) {
	h.configSetup = f
}

// Setup registers a setup hook.
func (h *hookstruct) Setup(f SetupFunc) {
	h.setup = f
}

// SetupHTTP registers a setup HTTP hook.
func (h *hookstruct) SetupHTTP(f SetupHTTPFunc) {
	h.setupHTTP = f
}

// invokeCleanup invokes the cleanup hook.
func (h *hookstruct) invokeCleanup() error {
	if h.cleanup != nil {
		return h.cleanup()
	}
	return nil
}

// invokeConfigSetup invokes the setup configuration hook.
func (h *hookstruct) invokeConfigSetup(c config.Config) error {
	if h.configSetup != nil {
		return h.configSetup(c)
	}
	return nil
}

// invokeSetup invokes the setup hook.
func (h *hookstruct) invokeSetup(workerGroup worker.Group) error {
	if h.setup != nil {
		return h.setup(workerGroup)
	}
	return nil
}

// invokeSetupHTTP invokes the setup HTTP hook.
func (h *hookstruct) invokeSetupHTTP() (*http.Server, error) {
	if h.setupHTTP != nil {
		var server *http.Server
		var err error
		if h.httpNewServer == nil {
			server, err = http.NewServer()
		} else {
			server, err = h.httpNewServer()
		}
		if err != nil {
			return nil, err
		}
		err = h.setupHTTP(server)
		if err != nil {
			return nil, err
		}
		return server, nil
	}
	return nil, nil
}
