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

// Test_NewGauge tests that NewGauge returns a Gauge
// without labels.
func Test_NewGauge(t *testing.T) {
	t.Parallel()
	m := New()
	g := m.NewGauge("test", "test")
	if g == nil {
		t.Fatal("NewGauge returned nil")
	}
	if _, ok := g.(*gauge); !ok {
		t.Fatal("NewGauge did not return a gauge")
	}
}

// Test_NewGauge_WithLabels tests that NewGauge returns a Gauge
// with labels.
func Test_NewGauge_WithLabels(t *testing.T) {
	t.Parallel()
	m := New()
	g := m.NewGauge("test", "test", "label")
	if g == nil {
		t.Fatal("NewGauge returned nil")
	}
	if _, ok := g.(*gaugeVec); !ok {
		t.Fatal("NewGauge did not return a gauge vector")
	}
}

// Test_Gauge_Set tests that Set sets the gauge.
func Test_Gauge_Set(t *testing.T) {
	t.Parallel()
	m := New()
	g := m.NewGauge("test", "test")
	g.Set(1)
}

// Test_Gauge_Inc tests that Inc increments the gauge.
func Test_Gauge_Inc(t *testing.T) {
	t.Parallel()
	m := New()
	g := m.NewGauge("test", "test")
	g.Inc()
}

// Test_Gauge_Dec tests that Dec increments the gauge.
func Test_Gauge_Dec(t *testing.T) {
	t.Parallel()
	m := New()
	g := m.NewGauge("test", "test")
	g.Dec()
}

// Test_Gauge_Add tests that Add increments the gauge.
func Test_Gauge_Add(t *testing.T) {
	t.Parallel()
	m := New()
	g := m.NewGauge("test", "test")
	g.Add(1)
}

// Test_Gauge_Sub tests that Sub increments the gauge.
func Test_Gauge_Sub(t *testing.T) {
	t.Parallel()
	m := New()
	g := m.NewGauge("test", "test")
	g.Sub(1)
}

// Test_Gauge_Set_WithLabels tests that Set sets the gauge with labels.
func Test_Gauge_Set_WithLabels(t *testing.T) {
	t.Parallel()
	m := New()
	g := m.NewGauge("test", "test", "label")
	g.Set(1, "value")
}

// Test_Gauge_Inc_WithLabels tests that Inc increments the gauge with
// labels.
func Test_Gauge_Inc_WithLabels(t *testing.T) {
	t.Parallel()
	m := New()
	g := m.NewGauge("test", "test", "label")
	g.Inc("value")
}

// Test_Gauge_Dec_WithLabels tests that Dec increments the gauge with
// labels.
func Test_Gauge_Dec_WithLabels(t *testing.T) {
	t.Parallel()
	m := New()
	g := m.NewGauge("test", "test", "label")
	g.Dec("value")
}

// Test_Gauge_Add_WithLabels tests that Add increments the gauge with
// labels.
func Test_Gauge_Add_WithLabels(t *testing.T) {
	t.Parallel()
	m := New()
	g := m.NewGauge("test", "test", "label")
	g.Add(1, "value")
}

// Test_Gauge_Sub_WithLabels tests that Sub increments the gauge with
// labels.
func Test_Gauge_Sub_WithLabels(t *testing.T) {
	t.Parallel()
	m := New()
	g := m.NewGauge("test", "test", "label")
	g.Sub(1, "value")
}

// Test_Gauge_Set_Failed_WithLabels tests that the bug handler is called if
// a gauge without labels is set with labels.
func Test_Gauge_Set_Failed_WithLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	g := m.NewGauge("test", "test")
	g.Set(1, "value")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Inc_Failed_WithLabels tests that the bug handler is called if
// a gauge without labels is incremented with labels.
func Test_Gauge_Inc_Failed_WithLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	g := m.NewGauge("test", "test")
	g.Inc("value")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Dec_Failed_WithLabels tests that the bug handler is called if
// a gauge without labels is decremented with labels.
func Test_Gauge_Dec_Failed_WithLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	g := m.NewGauge("test", "test")
	g.Dec("value")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Add_Failed_WithLabels tests that the bug handler is called if
// a gauge without labels is added with labels.
func Test_Gauge_Add_Failed_WithLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	g := m.NewGauge("test", "test")
	g.Add(1, "value")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Sub_Failed_WithLabels tests that the bug handler is called if
// a gauge without labels is subtracted with labels.
func Test_Gauge_Sub_Failed_WithLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	g := m.NewGauge("test", "test")
	g.Sub(1, "value")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Set_Failed_WithoutLabels tests that the bug handler is called
// if a gauge with labels is set without labels.
func Test_Gauge_Set_Failed_WithoutLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	g := m.NewGauge("test", "test", "label")
	g.Set(1)
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Inc_Failed_WithoutLabels tests that the bug handler is called
// if a gauge with labels is incremented without labels.
func Test_Gauge_Inc_Failed_WithoutLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	g := m.NewGauge("test", "test", "label")
	g.Inc()
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Dec_Failed_WithoutLabels tests that the bug handler is called
// if a gauge with labels is decremented without labels.
func Test_Gauge_Dec_Failed_WithoutLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	g := m.NewGauge("test", "test", "label")
	g.Dec()
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Add_Failed_WithoutLabels tests that the bug handler is called
// if a gauge with labels is added without labels.
func Test_Gauge_Add_Failed_WithoutLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	g := m.NewGauge("test", "test", "label")
	g.Add(1)
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Sub_Failed_WithoutLabels tests that the bug handler is called
// if a gauge with labels is subtracted without labels.
func Test_Gauge_Sub_Failed_WithoutLabels(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	g := m.NewGauge("test", "test", "label")
	g.Sub(1)
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Register_WithoutLabels_Twice tests that the bug handler is
// called if a gauge without labels is registered twice.
func Test_Gauge_Register_WithoutLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewGauge("test", "test")
	m.NewGauge("test", "test")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Register_WithLabels_Twice tests that the bug handler is
// called if a gauge with labels is registered twice.
func Test_Gauge_Register_WithLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewGauge("test", "test", "label")
	m.NewGauge("test", "test", "label")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Register_WithAndWithoutLabels_Twice tests that the bug handler
// is called if a gauge without labels is registered and then a gauge
// with labels is registered with the same name.
func Test_Gauge_Register_WithAndWithoutLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewGauge("test", "test")
	m.NewGauge("test", "test", "label")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Register_WithoutAndWithLabels_Twice tests that the bug handler
// is called if a gauge with labels is registered and then a gauge
// without labels is registered with the same name.
func Test_Gauge_Register_WithoutAndWithLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewGauge("test", "test", "label")
	m.NewGauge("test", "test")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Register_WithDifferentLabels_Twice tests that the bug handler
// is called if a gauge with labels is registered and then a gauge
// with different labels is registered with the same name.
func Test_Gauge_Register_WithDifferentLabels_Twice(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewGauge("test", "test", "label")
	m.NewGauge("test", "test", "label2")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Register_PrometheusRegistration_Failed tests that the bug
// handler is called if the Prometheus registration fails.
func Test_Gauge_Register_WithoutLabels_PrometheusRegistration_Failed(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewGauge("", "")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}

// Test_Gauge_Register_WithLabels_PrometheusRegistration_Failed tests that the
// bug handler is called if the Prometheus registration fails.
func Test_Gauge_Register_WithLabels_PrometheusRegistration_Failed(t *testing.T) {
	var bugMsg string
	defer overrideBugHandler(t, &bugMsg)()
	m := New()
	m.NewGauge("", "", "")
	if bugMsg == "" {
		t.Fatal("bug handler was not called")
	}
}
