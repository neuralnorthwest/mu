package service

import (
	"github.com/neuralnorthwest/mu/logging"
	"github.com/neuralnorthwest/mu/status"
	"golang.org/x/mod/semver"
)

// Option is an option for a service.
type Option func(*Service) error

// WithVersion returns an option that sets the version of the service.
func WithVersion(version string) Option {
	return func(s *Service) error {
		if !semver.IsValid(version) {
			return status.ErrInvalidVersion
		}
		s.version = version
		return nil
	}
}

// WithMockMode returns an option that sets the mock mode of the service.
func WithMockMode() Option {
	return func(s *Service) error {
		s.mockMode = true
		return nil
	}
}

// WithLogger returns an option that sets the func used to create a new logger.
func WithLogger(newLogger func() (logging.Logger, error)) Option {
	return func(s *Service) error {
		s.newLogger = newLogger
		return nil
	}
}
