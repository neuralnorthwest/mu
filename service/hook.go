package service

import "github.com/neuralnorthwest/mu/config"

// SetupConfigFunc is a function that sets up a service configuration.
type SetupConfigFunc func(c config.Config) error

// SetupFunc is a function that sets up a service.
type SetupFunc func() error

// CleanupFunc is a function that cleans up a service.
type CleanupFunc func() error

type hooks interface {
	RegisterSetupConfig(f SetupConfigFunc)
	RegisterSetup(f SetupFunc)
	RegisterCleanup(f CleanupFunc)

	invokeSetupConfig(c config.Config) error
	invokeSetup() error
	invokeCleanup() error
}

// hookstruct holds the hookstruct for a service.
type hookstruct struct {
	// setupConfig is the setup configuration hook.
	setupConfig SetupConfigFunc
	// setup is the setup hook.
	setup SetupFunc
	// cleanup is the cleanup hook.
	cleanup CleanupFunc
}

var _ hooks = (*hookstruct)(nil)

// RegisterSetupConfig registers a setup configuration hook.
func (h *hookstruct) RegisterSetupConfig(f SetupConfigFunc) {
	h.setupConfig = f
}

// RegisterSetup registers a setup hook.
func (h *hookstruct) RegisterSetup(f SetupFunc) {
	h.setup = f
}

// RegisterCleanup registers a cleanup hook.
func (h *hookstruct) RegisterCleanup(f CleanupFunc) {
	h.cleanup = f
}

// invokeSetupConfig invokes the setup configuration hook.
func (h *hookstruct) invokeSetupConfig(c config.Config) error {
	if h.setupConfig != nil {
		return h.setupConfig(c)
	}
	return nil
}

// invokeSetup invokes the setup hook.
func (h *hookstruct) invokeSetup() error {
	if h.setup != nil {
		return h.setup()
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
