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
