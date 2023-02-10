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

// PreRunFunc is a function that runs before the workers are started.
type PreRunFunc func() error

// SetupConfigFunc is a function that sets up a service configuration.
type SetupConfigFunc func(c config.Config) error

// SetupWorkersFunc is a function that sets up workers for the service.
type SetupWorkersFunc func(workerGroup worker.Group) error

// SetupHTTPFunc is a function that sets up a service HTTP server.
type SetupHTTPFunc func(server *http.Server) error

type Hooks interface {
	PreRun(f PreRunFunc)
	SetupConfig(f SetupConfigFunc)
	SetupWorkers(f SetupWorkersFunc)
	SetupHTTP(f SetupHTTPFunc)

	invokePreRun() error
	invokeSetupConfig(c config.Config) error
	invokeSetupWorkers(workerGroup worker.Group) error
	invokeSetupHTTP() (*http.Server, error)
}

// hookstruct holds the hooks for a service.
type hookstruct struct {
	// setupConfig is the setup configuration hook.
	setupConfig SetupConfigFunc
	// prerun is the prerun hook.
	prerun PreRunFunc
	// setupWorkers is the setupWorkers hook.
	setupWorkers SetupWorkersFunc
	// setupHTTP is the setup HTTP hook.
	setupHTTP SetupHTTPFunc
	// httpNewServer is a function that creates a new HTTP server. This is used
	// for testing.
	httpNewServer func() (*http.Server, error)
}

var _ Hooks = (*hookstruct)(nil)

// SetupConfig registers a configuration setup hook.
func (h *hookstruct) SetupConfig(f SetupConfigFunc) {
	h.setupConfig = f
}

// PreRun registers a prerun hook.
func (h *hookstruct) PreRun(f PreRunFunc) {
	h.prerun = f
}

// SetupWorkers registers a worker setup hook.
func (h *hookstruct) SetupWorkers(f SetupWorkersFunc) {
	h.setupWorkers = f
}

// SetupHTTP registers a setup HTTP hook.
func (h *hookstruct) SetupHTTP(f SetupHTTPFunc) {
	h.setupHTTP = f
}

// invokeSetupConfig invokes the setup configuration hook.
func (h *hookstruct) invokeSetupConfig(c config.Config) error {
	if h.setupConfig != nil {
		return h.setupConfig(c)
	}
	return nil
}

// invokePreRun invokes the prerun hook.
func (h *hookstruct) invokePreRun() error {
	if h.prerun != nil {
		return h.prerun()
	}
	return nil
}

// invokeSetupWorkers invokes the setup hook.
func (h *hookstruct) invokeSetupWorkers(workerGroup worker.Group) error {
	if h.setupWorkers != nil {
		return h.setupWorkers(workerGroup)
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
