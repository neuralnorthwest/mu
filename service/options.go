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

// WithLogger returns an option that sets the logger of the service.
func WithLogger(logger logging.Logger) Option {
	return func(s *Service) error {
		s.logger = logger
		return nil
	}
}
