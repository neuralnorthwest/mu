package service

import (
	"context"

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
	hooks
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
}

// New returns a new service.
func New(name string, opts ...Option) (*Service, error) {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		name:     name,
		version:  "v0.0.0",
		hooks:    &hookstruct{},
		ctx:      ctx,
		cancel:   cancel,
		config:   config.New(),
		mockMode: false,
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	if s.logger == nil {
		logger, err := logging.New()
		if err != nil {
			return nil, err
		}
		s.logger = logger
	}
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
