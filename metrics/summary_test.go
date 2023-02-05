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

// Test_NewSummary tests that NewSummary returns a Summary
// without labels.
func Test_NewSummary(t *testing.T) {
	t.Parallel()
	m := New()
	s := m.NewSummary("test", "test")
	if s == nil {
		t.Fatal("NewSummary returned nil")
	}
	if _, ok := s.(*summary); !ok {
		t.Fatal("NewSummary did not return a summary")
	}
}

// Test_NewSummary_WithLabels tests that NewSummary returns a Summary
// with labels.
func Test_NewSummary_WithLabels(t *testing.T) {
	t.Parallel()
	m := New()
	s := m.NewSummary("test", "test", "label")
	if s == nil {
		t.Fatal("NewSummary returned nil")
	}
	if _, ok := s.(*summaryVec); !ok {
		t.Fatal("NewSummary did not return a summary vector")
	}
}

// Test_Summary_Observe tests that Observe updates the summary.
func Test_Summary_Observe(t *testing.T) {
	t.Parallel()
	m := New()
	s := m.NewSummary("test", "test")
	s.Observe(1)
}

// Test_Summary_Observe_WithLabels tests that Observe updates the summary
// with labels.
func Test_Summary_Observe_WithLabels(t *testing.T) {
	t.Parallel()
	m := New()
	s := m.NewSummary("test", "test", "label")
	s.Observe(1, "value")
}

// Test_Summary_Observe_Failed_WithLabels tests that the bug handler is called
// if a summary without labels is incremented with labels.
func Test_Summary_Observe_Failed_WithLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	s := m.NewSummary("test", "test")
	s.Observe(1, "value")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Summary_Observe_Failed_WithoutLabels tests that the bug handler is
// called if a summary with labels is incremented without labels.
func Test_Summary_Observe_Failed_WithoutLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	s := m.NewSummary("test", "test", "label")
	s.Observe(1)
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Summary_Register_WithoutLabels_Twice tests that the bug handler is
// called if a summary without labels is registered twice.
func Test_Summary_Register_WithoutLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewSummary("test", "test")
	m.NewSummary("test", "test")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Summary_Register_WithLabels_Twice tests that the bug handler is
// called if a summary with labels is registered twice.
func Test_Summary_Register_WithLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewSummary("test", "test", "label")
	m.NewSummary("test", "test", "label")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Summary_Register_WithAndWithoutLabels_Twice tests that the bug handler
// is called if a summary without labels is registered and then a summary
// with labels is registered with the same name.
func Test_Summary_Register_WithAndWithoutLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewSummary("test", "test")
	m.NewSummary("test", "test", "label")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Summary_Register_WithoutAndWithLabels_Twice tests that the bug handler
// is called if a summary with labels is registered and then a summary
// without labels is registered with the same name.
func Test_Summary_Register_WithoutAndWithLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewSummary("test", "test", "label")
	m.NewSummary("test", "test")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Summary_Register_WithDifferentLabels_Twice tests that the bug handler
// is called if a summary with labels is registered and then a summary
// with different labels is registered with the same name.
func Test_Summary_Register_WithDifferentLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewSummary("test", "test", "label")
	m.NewSummary("test", "test", "label2")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Summary_Register_PrometheusRegistration_Failed tests that the bug
// handler is called if the Prometheus registration fails.
func Test_Summary_Register_WithoutLabels_PrometheusRegistration_Failed(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewSummary("", "")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Summary_Register_WithLabels_PrometheusRegistration_Failed tests that the
// bug handler is called if the Prometheus registration fails.
func Test_Summary_Register_WithLabels_PrometheusRegistration_Failed(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewSummary("", "", "")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}
