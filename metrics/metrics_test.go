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

package metrics

import (
	"testing"

	"github.com/neuralnorthwest/mu/status"
	"github.com/prometheus/client_golang/prometheus"
)

// Test_Metrics_WithRegistry tests that the metrics package can be used with a
// custom registry.
func Test_Metrics_WithRegistry(t *testing.T) {
	t.Parallel()
	reg := prometheus.NewRegistry()
	m := newMetrics(t, WithRegistry(reg))
	if m.registry != reg {
		t.Fatal("registry not set")
	}
}

// errorOption is an option that always returns an error.
type errorOption struct{}

func (errorOption) apply(*metrics) error {
	return status.ErrInvalidArgument
}

// withErrorOption is an option that always returns an error.
func withErrorOption() Option {
	return &errorOption{}
}

// Test_Metrics_OptionError tests that the metrics package returns an error if
// an invalid option is passed.
func Test_Metrics_OptionError(t *testing.T) {
	t.Parallel()
	_, err := New(withErrorOption())
	if err == nil {
		t.Fatal("expected error")
	}
}
