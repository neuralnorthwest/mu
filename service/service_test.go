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
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/neuralnorthwest/mu/logging"
	mock_logging "github.com/neuralnorthwest/mu/logging/mock"
	"github.com/neuralnorthwest/mu/status"
)

// serviceOptionError is an Option that returns an error.
func serviceOptionError(t *testing.T, err error) Option {
	t.Helper()
	return func(s *Service) error {
		return err
	}
}

// Test_Service_New tests the New function.
func Test_Service_New(t *testing.T) {
	t.Parallel()
	svc, err := New("test-service")
	if err != nil {
		t.Fatalf("New returned an error: %v", err)
	}
	if svc.Name() != "test-service" {
		t.Errorf("unexpected service name: %s, expected: %s", svc.Name(), "test-service")
	}
	if svc.Version() != "v0.0.0" {
		t.Errorf("unexpected service version: %s, expected: %s", svc.Version(), "v0.0.0")
	}
	if reflect.ValueOf(svc.Logger()).IsNil() {
		t.Error("unexpected nil logger")
	}
	if reflect.ValueOf(svc.Context()).IsNil() {
		t.Error("unexpected nil context")
	}
	if reflect.ValueOf(svc.Config()).IsNil() {
		t.Error("unexpected nil config")
	}
	if svc.MockMode() {
		t.Error("unexpected mock mode")
	}
}

// Test_Service_Context tests the Context and Cancel functions.
func Test_Service_Context(t *testing.T) {
	t.Parallel()
	svc, err := New("test-service")
	if err != nil {
		t.Fatalf("New returned an error: %v", err)
	}
	// Done should not be closed.
	select {
	case <-svc.Context().Done():
		t.Error("unexpected context done")
	default:
	}
	svc.Cancel()
	// Done should be closed.
	select {
	case <-svc.Context().Done():
	default:
		t.Error("expected context done")
	}
}

// Test_Service_New_WithVersion tests the New function with a version.
func Test_Service_New_WithVersion(t *testing.T) {
	t.Parallel()
	svc, err := New("test-service", WithVersion("v1.0.0"))
	if err != nil {
		t.Fatalf("New returned an error: %v", err)
	}
	if svc.Version() != "v1.0.0" {
		t.Errorf("unexpected service version: %s, expected: %s", svc.Version(), "v1.0.0")
	}
}

// Test_Service_New_WithVersionError tests the New function with a version error.
func Test_Service_New_WithVersionError(t *testing.T) {
	t.Parallel()
	svc, err := New("test-service", WithVersion("1.0.0"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, status.ErrInvalidVersion) {
		t.Errorf("unexpected error: %v, expected: %v", err, status.ErrInvalidVersion)
	}
	if svc != nil {
		t.Error("expected nil service")
	}
}

// Test_Service_New_WithMockMode tests the New function with mock mode.
func Test_Service_New_WithMockMode(t *testing.T) {
	t.Parallel()
	svc, err := New("test-service", WithMockMode())
	if err != nil {
		t.Fatalf("New returned an error: %v", err)
	}
	if !svc.MockMode() {
		t.Error("expected mock mode")
	}
}

// Test_Service_New_WithLogger tests the New function with a logger.
func Test_Service_New_WithLogger(t *testing.T) {
	t.Parallel()
	mc := gomock.NewController(t)
	logger := mock_logging.NewMockLogger(mc)
	logger.EXPECT().Info("test")
	svc, err := New("test-service", WithLogger(func() (logging.Logger, error) {
		return logger, nil
	}))
	if err != nil {
		t.Fatalf("New returned an error: %v", err)
	}
	if svc.Logger() != logger {
		t.Error("unexpected logger")
	}
	svc.Logger().Info("test")
}

// Test_Service_New_WithLoggerError tests the New function with a logger constructor error.
func Test_Service_New_WithLoggerError(t *testing.T) {
	t.Parallel()
	svc, err := New("test-service", WithLogger(func() (logging.Logger, error) {
		return nil, status.ErrInvalidArgument
	}))
	if !errors.Is(err, status.ErrInvalidArgument) {
		t.Errorf("unexpected error: %v, expected: %v", err, status.ErrInvalidArgument)
	}
	if svc != nil {
		t.Error("unexpected service")
	}
}

// Test_Service_New_serviceOptionError tests the New function with a serviceOptionError.
func Test_Service_New_serviceOptionError(t *testing.T) {
	t.Parallel()
	svc, err := New("test-service", serviceOptionError(t, status.ErrInvalidArgument))
	if !errors.Is(err, status.ErrInvalidArgument) {
		t.Errorf("unexpected error: %v, expected: %v", err, status.ErrInvalidArgument)
	}
	if svc != nil {
		t.Error("unexpected service")
	}
}
