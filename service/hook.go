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
	"github.com/neuralnorthwest/mu/worker"
)

// ConfigSetupFunc is a function that sets up a service configuration.
type ConfigSetupFunc func(c config.Config) error

// SetupFunc is a function that sets up a service.
type SetupFunc func(workerGroup worker.Group) error

// CleanupFunc is a function that cleans up a service.
type CleanupFunc func() error

type hooks interface {
	RegisterConfigSetup(f ConfigSetupFunc)
	RegisterSetup(f SetupFunc)
	RegisterCleanup(f CleanupFunc)

	invokeConfigSetup(c config.Config) error
	invokeSetup(workerGroup worker.Group) error
	invokeCleanup() error
}

// hookstruct holds the hookstruct for a service.
type hookstruct struct {
	// configSetup is the setup configuration hook.
	configSetup ConfigSetupFunc
	// setup is the setup hook.
	setup SetupFunc
	// cleanup is the cleanup hook.
	cleanup CleanupFunc
}

var _ hooks = (*hookstruct)(nil)

// RegisterSetupConfig registers a setup configuration hook.
func (h *hookstruct) RegisterConfigSetup(f ConfigSetupFunc) {
	h.configSetup = f
}

// RegisterSetup registers a setup hook.
func (h *hookstruct) RegisterSetup(f SetupFunc) {
	h.setup = f
}

// RegisterCleanup registers a cleanup hook.
func (h *hookstruct) RegisterCleanup(f CleanupFunc) {
	h.cleanup = f
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

// invokeCleanup invokes the cleanup hook.
func (h *hookstruct) invokeCleanup() error {
	if h.cleanup != nil {
		return h.cleanup()
	}
	return nil
}
