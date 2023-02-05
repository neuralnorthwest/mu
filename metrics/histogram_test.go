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
)

// Test_NewHistogram tests that NewHistogram returns a Histogram
// without labels.
func Test_NewHistogram(t *testing.T) {
	t.Parallel()
	m := New()
	h := m.NewHistogram("test", "test", nil)
	if h == nil {
		t.Fatal("NewHistogram returned nil")
	}
	if _, ok := h.(*histogram); !ok {
		t.Fatal("NewHistogram did not return a histogram")
	}
}

// Test_NewHistogram_WithLabels tests that NewHistogram returns a Histogram
// with labels.
func Test_NewHistogram_WithLabels(t *testing.T) {
	t.Parallel()
	m := New()
	h := m.NewHistogram("test", "test", nil, "label")
	if h == nil {
		t.Fatal("NewHistogram returned nil")
	}
	if _, ok := h.(*histogramVec); !ok {
		t.Fatal("NewHistogram did not return a histogram vector")
	}
}

// Test_Histogram_Observe tests that Observe updates the histogram.
func Test_Histogram_Observe(t *testing.T) {
	t.Parallel()
	m := New()
	h := m.NewHistogram("test", "test", nil)
	h.Observe(1)
}

// Test_Histogram_Observe_WithLabels tests that Observe updates the histogram
// with labels.
func Test_Histogram_Observe_WithLabels(t *testing.T) {
	t.Parallel()
	m := New()
	h := m.NewHistogram("test", "test", nil, "label")
	h.Observe(1, "value")
}

// Test_Histogram_Observe_Failed_WithLabels tests that the bug handler is called
// if a histogram without labels is incremented with labels.
func Test_Histogram_Observe_Failed_WithLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	h := m.NewHistogram("test", "test", nil)
	h.Observe(1, "value")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Histogram_Observe_Failed_WithoutLabels tests that the bug handler is
// called if a histogram with labels is incremented without labels.
func Test_Histogram_Observe_Failed_WithoutLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	h := m.NewHistogram("test", "test", nil, "label")
	h.Observe(1)
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Histogram_Register_WithoutLabels_Twice tests that the bug handler is
// called if a histogram without labels is registered twice.
func Test_Histogram_Register_WithoutLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewHistogram("test", "test", nil)
	m.NewHistogram("test", "test", nil)
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Histogram_Register_WithLabels_Twice tests that the bug handler is
// called if a histogram with labels is registered twice.
func Test_Histogram_Register_WithLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewHistogram("test", "test", nil, "label")
	m.NewHistogram("test", "test", nil, "label")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Histogram_Register_WithAndWithoutLabels_Twice tests that the bug handler
// is called if a histogram without labels is registered and then a histogram
// with labels is registered with the same name.
func Test_Histogram_Register_WithAndWithoutLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewHistogram("test", "test", nil)
	m.NewHistogram("test", "test", nil, "label")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Histogram_Register_WithoutAndWithLabels_Twice tests that the bug handler
// is called if a histogram with labels is registered and then a histogram
// without labels is registered with the same name.
func Test_Histogram_Register_WithoutAndWithLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewHistogram("test", "test", nil, "label")
	m.NewHistogram("test", "test", nil)
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Histogram_Register_WithDifferentLabels_Twice tests that the bug handler
// is called if a histogram with labels is registered and then a histogram
// with different labels is registered with the same name.
func Test_Histogram_Register_WithDifferentLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewHistogram("test", "test", nil, "label")
	m.NewHistogram("test", "test", nil, "label2")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Histogram_Register_PrometheusRegistration_Failed tests that the bug
// handler is called if the Prometheus registration fails.
func Test_Histogram_Register_WithoutLabels_PrometheusRegistration_Failed(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewHistogram("", "", nil)
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Histogram_Register_WithLabels_PrometheusRegistration_Failed tests that the
// bug handler is called if the Prometheus registration fails.
func Test_Histogram_Register_WithLabels_PrometheusRegistration_Failed(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewHistogram("", "", nil, "")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}
