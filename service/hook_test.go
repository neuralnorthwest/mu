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
	"testing"

	"github.com/neuralnorthwest/mu/config"
	"github.com/neuralnorthwest/mu/http"
	"github.com/neuralnorthwest/mu/status"
	"github.com/neuralnorthwest/mu/worker"
)

// Test_hooks_NoHooks tests that no hooks are invoked when no hooks are registered.
func Test_hooks_NoHooks(t *testing.T) {
	t.Parallel()
	h := &hookstruct{}
	if err := h.invokeConfigSetup(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := h.invokeSetupWorkers(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := h.invokeCleanup(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// Test_hooks_InvokeCleanup tests that the cleanup hook is invoked.
func Test_hooks_InvokeCleanup(t *testing.T) {
	t.Parallel()
	wasInvoked := false
	h := &hookstruct{}
	h.Cleanup(func() error {
		wasInvoked = true
		return nil
	})
	if err := h.invokeCleanup(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !wasInvoked {
		t.Error("expected cleanup hook to be invoked")
	}
}

// Test_hooks_InvokeConfigSetup tests that the config setup hook is invoked.
func Test_hooks_InvokeConfigSetup(t *testing.T) {
	t.Parallel()
	wasInvoked := false
	h := &hookstruct{}
	h.ConfigSetup(func(c config.Config) error {
		wasInvoked = true
		return nil
	})
	if err := h.invokeConfigSetup(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !wasInvoked {
		t.Error("expected config setup hook to be invoked")
	}
}

// Test_hooks_InvokePreRun tests that the pre-run hook is invoked.
func Test_hooks_InvokePreRun(t *testing.T) {
	t.Parallel()
	wasInvoked := false
	h := &hookstruct{}
	h.PreRun(func() error {
		wasInvoked = true
		return nil
	})
	if err := h.invokePreRun(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !wasInvoked {
		t.Error("expected pre-run hook to be invoked")
	}
}

// Test_hooks_InvokeSetupWorkers tests that the setupWorkers hook is invoked.
func Test_hooks_InvokeSetupWorkers(t *testing.T) {
	t.Parallel()
	wasInvoked := false
	h := &hookstruct{}
	h.SetupWorkers(func(workerGroup worker.Group) error {
		wasInvoked = true
		return nil
	})
	if err := h.invokeSetupWorkers(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !wasInvoked {
		t.Error("expected setupWorkers hook to be invoked")
	}
}

// Test_hooks_InvokeSetupHTTP tests that the setup HTTP hook is invoked.
func Test_hooks_InvokeSetupHTTP(t *testing.T) {
	t.Parallel()
	wasInvoked := false
	h := &hookstruct{}
	h.SetupHTTP(func(server *http.Server) error {
		wasInvoked = true
		return nil
	})
	var httpServer *http.Server
	var err error
	if httpServer, err = h.invokeSetupHTTP(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !wasInvoked {
		t.Error("expected setup HTTP hook to be invoked")
	}
	if httpServer == nil {
		t.Error("expected http server to be returned")
	}
}

// Test_hooks_InvokeSetupHTTP_NewServerError tests that the setup HTTP hook is invoked
// and checks the logic when we fail to create a new server.
func Test_hooks_InvokeSetupHTTP_NewServerError(t *testing.T) {
	t.Parallel()
	h := &hookstruct{}
	h.httpNewServer = func() (*http.Server, error) {
		return nil, status.ErrInvalidArgument
	}
	h.SetupHTTP(func(server *http.Server) error {
		return nil
	})
	if _, err := h.invokeSetupHTTP(); err == nil {
		t.Error("expected error")
	}
}

// Test_hooks_InvokeSetupHTTP_HookError tests that the setup HTTP hook is invoked
// and checks the logic when the user's hook returns an error.
func Test_hooks_InvokeSetupHTTP_HookError(t *testing.T) {
	t.Parallel()
	h := &hookstruct{}
	h.SetupHTTP(func(server *http.Server) error {
		return status.ErrInvalidArgument
	})
	if _, err := h.invokeSetupHTTP(); err == nil {
		t.Error("expected error")
	}
}
