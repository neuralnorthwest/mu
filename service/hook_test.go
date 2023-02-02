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
	"github.com/neuralnorthwest/mu/worker"
)

// Test_hooks_NoHooks tests that no hooks are invoked when no hooks are registered.
func Test_hooks_NoHooks(t *testing.T) {
	t.Parallel()
	h := &hookstruct{}
	if err := h.invokeConfigSetup(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := h.invokeSetup(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := h.invokeCleanup(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// Test_hooks_InvokeConfigSetup tests that the config setup hook is invoked.
func Test_hooks_InvokeConfigSetup(t *testing.T) {
	t.Parallel()
	wasInvoked := false
	h := &hookstruct{}
	h.RegisterConfigSetup(func(c config.Config) error {
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

// Test_hooks_InvokeSetup tests that the setup hook is invoked.
func Test_hooks_InvokeSetup(t *testing.T) {
	t.Parallel()
	wasInvoked := false
	h := &hookstruct{}
	h.RegisterSetup(func(workerGroup worker.Group) error {
		wasInvoked = true
		return nil
	})
	if err := h.invokeSetup(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !wasInvoked {
		t.Error("expected setup hook to be invoked")
	}
}

// Test_hooks_InvokeCleanup tests that the cleanup hook is invoked.
func Test_hooks_InvokeCleanup(t *testing.T) {
	t.Parallel()
	wasInvoked := false
	h := &hookstruct{}
	h.RegisterCleanup(func() error {
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
