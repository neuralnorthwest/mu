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

	"github.com/neuralnorthwest/mu/bug"
)

// overrideBugHandler sets the bug handler to a function that records the bug
// message in the given string pointer. It returns a function that restores the
// original bug handler.
func overrideBugHandler(t *testing.T, bugMessage *string) func() {
	t.Helper()
	originalHandler := bug.Handler()
	bug.SetHandler(func(message string) {
		*bugMessage = message
	})
	return func() {
		bug.SetHandler(originalHandler)
	}
}

// Test_NewCounter tests that NewCounter returns a Counter
// without labels.
func Test_NewCounter(t *testing.T) {
	t.Parallel()
	m := newMetrics(t)
	c := m.NewCounter("test", "test")
	if c == nil {
		t.Fatal("NewCounter returned nil")
	}
	if _, ok := c.(*counter); !ok {
		t.Fatal("NewCounter did not return a counter")
	}
}

// Test_NewCounter_WithLabels tests that NewCounter returns a Counter
// with labels.
func Test_NewCounter_WithLabels(t *testing.T) {
	t.Parallel()
	m := newMetrics(t)
	c := m.NewCounter("test", "test", "label")
	if c == nil {
		t.Fatal("NewCounter returned nil")
	}
	if _, ok := c.(*counterVec); !ok {
		t.Fatal("NewCounter did not return a counter vector")
	}
}

// Test_Counter_Inc tests that Inc increments the counter.
func Test_Counter_Inc(t *testing.T) {
	t.Parallel()
	m := newMetrics(t)
	c := m.NewCounter("test", "test")
	c.Inc()
}

// Test_Counter_Add tests that Add increments the counter.
func Test_Counter_Add(t *testing.T) {
	t.Parallel()
	m := newMetrics(t)
	c := m.NewCounter("test", "test")
	c.Add(1)
}

// Test_Counter_Inc_WithLabels tests that Inc increments the counter with
// labels.
func Test_Counter_Inc_WithLabels(t *testing.T) {
	t.Parallel()
	m := newMetrics(t)
	c := m.NewCounter("test", "test", "label")
	c.Inc("value")
}

// Test_Counter_Add_WithLabels tests that Add increments the counter with
// labels.
func Test_Counter_Add_WithLabels(t *testing.T) {
	t.Parallel()
	m := newMetrics(t)
	c := m.NewCounter("test", "test", "label")
	c.Add(1, "value")
}

// Test_Counter_Inc_Failed_WithLabels tests that the bug handler is called if
// a counter without labels is incremented with labels.
func Test_Counter_Inc_Failed_WithLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := newMetrics(t)
	c := m.NewCounter("test", "test")
	c.Inc("value")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Counter_Add_Failed_WithLabels tests that the bug handler is called if
// a counter without labels is added with labels.
func Test_Counter_Add_Failed_WithLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := newMetrics(t)
	c := m.NewCounter("test", "test")
	c.Add(1, "value")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Counter_Inc_Failed_WithoutLabels tests that the bug handler is called
// if a counter with labels is incremented without labels.
func Test_Counter_Inc_Failed_WithoutLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := newMetrics(t)
	c := m.NewCounter("test", "test", "label")
	c.Inc()
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Counter_Add_Failed_WithoutLabels tests that the bug handler is called
// if a counter with labels is added without labels.
func Test_Counter_Add_Failed_WithoutLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := newMetrics(t)
	c := m.NewCounter("test", "test", "label")
	c.Add(1)
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Counter_Register_WithoutLabels_Twice tests that the bug handler is
// called if a counter without labels is registered twice.
func Test_Counter_Register_WithoutLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := newMetrics(t)
	m.NewCounter("test", "test")
	m.NewCounter("test", "test")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Counter_Register_WithLabels_Twice tests that the bug handler is
// called if a counter with labels is registered twice.
func Test_Counter_Register_WithLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := newMetrics(t)
	m.NewCounter("test", "test", "label")
	m.NewCounter("test", "test", "label")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Counter_Register_WithAndWithoutLabels_Twice tests that the bug handler
// is called if a counter without labels is registered and then a counter
// with labels is registered with the same name.
func Test_Counter_Register_WithAndWithoutLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := newMetrics(t)
	m.NewCounter("test", "test")
	m.NewCounter("test", "test", "label")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Counter_Register_WithoutAndWithLabels_Twice tests that the bug handler
// is called if a counter with labels is registered and then a counter
// without labels is registered with the same name.
func Test_Counter_Register_WithoutAndWithLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := newMetrics(t)
	m.NewCounter("test", "test", "label")
	m.NewCounter("test", "test")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Counter_Register_WithDifferentLabels_Twice tests that the bug handler
// is called if a counter with labels is registered and then a counter
// with different labels is registered with the same name.
func Test_Counter_Register_WithDifferentLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := newMetrics(t)
	m.NewCounter("test", "test", "label")
	m.NewCounter("test", "test", "label2")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Counter_Register_PrometheusRegistration_Failed tests that the bug
// handler is called if the Prometheus registration fails.
func Test_Counter_Register_WithoutLabels_PrometheusRegistration_Failed(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := newMetrics(t)
	m.NewCounter("", "")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Counter_Register_WithLabels_PrometheusRegistration_Failed tests that the
// bug handler is called if the Prometheus registration fails.
func Test_Counter_Register_WithLabels_PrometheusRegistration_Failed(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := newMetrics(t)
	m.NewCounter("", "", "")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Counter_Dummy tests the dummy counter that is returned if an error
// occurs.
func Test_Counter_Dummy(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := newMetrics(t)
	c := m.NewCounter("", "", "")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
	c.Inc()
	c.Add(1)
}
