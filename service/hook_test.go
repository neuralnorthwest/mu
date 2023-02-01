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
