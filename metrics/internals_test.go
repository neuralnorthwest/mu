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

import "testing"

// Test_PrometheusRegistry tests that PrometheusRegistry returns the internal
// prometheus registry.
func Test_PrometheusRegistry(t *testing.T) {
	t.Parallel()
	m := New()
	reg := PrometheusRegistry(m)
	if reg == nil {
		t.Fatal("PrometheusRegistry returned nil")
	}
}

// Test_PrometheusCounterVec tests that PrometheusCounterVec returns the
// internal prometheus counter vector.
func Test_PrometheusCounterVec(t *testing.T) {
	t.Parallel()
	m := New()
	c := m.NewCounter("test", "test", "label")
	c.Inc("value")
	cv := PrometheusCounterVec(c)
	if cv == nil {
		t.Fatal("PrometheusCounterVec returned nil")
	}
}

// Test_PrometheusGaugeVec tests that PrometheusGaugeVec returns the internal
// prometheus gauge vector.
func Test_PrometheusGaugeVec(t *testing.T) {
	t.Parallel()
	m := New()
	g := m.NewGauge("test", "test", "label")
	g.Set(1, "value")
	gv := PrometheusGaugeVec(g)
	if gv == nil {
		t.Fatal("PrometheusGaugeVec returned nil")
	}
}

// Test_PrometheusHistogramVec tests that PrometheusHistogramVec returns the
// internal prometheus histogram vector.
func Test_PrometheusHistogramVec(t *testing.T) {
	t.Parallel()
	m := New()
	h := m.NewHistogram("test", "test", nil, "label")
	h.Observe(1, "value")
	hv := PrometheusHistogramVec(h)
	if hv == nil {
		t.Fatal("PrometheusHistogramVec returned nil")
	}
}
