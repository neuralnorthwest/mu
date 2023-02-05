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
	"github.com/neuralnorthwest/mu/bug"
	"github.com/prometheus/client_golang/prometheus"
)

// Gauge is a type that allows you to set the value of a gauge.
type Gauge interface {
	// Set sets the value of the gauge.
	Set(value float64, labelValues ...string)
	// Inc increments the value of the gauge by 1.
	Inc(labelValues ...string)
	// Dec decrements the value of the gauge by 1.
	Dec(labelValues ...string)
	// Add increments the value of the gauge by the given value.
	Add(value float64, labelValues ...string)
	// Sub decrements the value of the gauge by the given value.
	Sub(value float64, labelValues ...string)
}

// gauge is the default implementation of Gauge.
type gauge struct {
	// name is the name of the gauge.
	name string
	// gauge is the prometheus gauge.
	gauge prometheus.Gauge
}

// Set sets the value of the gauge.
func (g *gauge) Set(value float64, labelValues ...string) {
	if len(labelValues) > 0 {
		bug.Bugf("gauge %s has no labels", g.name)
		return
	}
	g.gauge.Set(value)
}

// Inc increments the value of the gauge by 1.
func (g *gauge) Inc(labelValues ...string) {
	if len(labelValues) > 0 {
		bug.Bugf("gauge %s has no labels", g.name)
		return
	}
	g.gauge.Inc()
}

// Dec decrements the value of the gauge by 1.
func (g *gauge) Dec(labelValues ...string) {
	if len(labelValues) > 0 {
		bug.Bugf("gauge %s has no labels", g.name)
		return
	}
	g.gauge.Dec()
}

// Add increments the value of the gauge by the given value.
func (g *gauge) Add(value float64, labelValues ...string) {
	if len(labelValues) > 0 {
		bug.Bugf("gauge %s has no labels", g.name)
		return
	}
	g.gauge.Add(value)
}

// Sub decrements the value of the gauge by the given value.
func (g *gauge) Sub(value float64, labelValues ...string) {
	if len(labelValues) > 0 {
		bug.Bugf("gauge %s has no labels", g.name)
		return
	}
	g.gauge.Sub(value)
}

// gaugeVec is the default implementation of Gauge for gauges with labels.
type gaugeVec struct {
	// name is the name of the gauge.
	name string
	// gaugeVec is the prometheus gauge vec.
	gaugeVec *prometheus.GaugeVec
}

// Set sets the value of the gauge.
func (g *gaugeVec) Set(value float64, labelValues ...string) {
	withLabels, err := g.gaugeVec.GetMetricWithLabelValues(labelValues...)
	if err != nil {
		bug.Bugf("gauge %s: %s", g.name, err)
		return
	}
	withLabels.Set(value)
}

// Inc increments the value of the gauge by 1.
func (g *gaugeVec) Inc(labelValues ...string) {
	withLabels, err := g.gaugeVec.GetMetricWithLabelValues(labelValues...)
	if err != nil {
		bug.Bugf("gauge %s: %s", g.name, err)
		return
	}
	withLabels.Inc()
}

// Dec decrements the value of the gauge by 1.
func (g *gaugeVec) Dec(labelValues ...string) {
	withLabels, err := g.gaugeVec.GetMetricWithLabelValues(labelValues...)
	if err != nil {
		bug.Bugf("gauge %s: %s", g.name, err)
		return
	}
	withLabels.Dec()
}

// Add increments the value of the gauge by the given value.
func (g *gaugeVec) Add(value float64, labelValues ...string) {
	withLabels, err := g.gaugeVec.GetMetricWithLabelValues(labelValues...)
	if err != nil {
		bug.Bugf("gauge %s: %s", g.name, err)
		return
	}
	withLabels.Add(value)
}

// Sub decrements the value of the gauge by the given value.
func (g *gaugeVec) Sub(value float64, labelValues ...string) {
	withLabels, err := g.gaugeVec.GetMetricWithLabelValues(labelValues...)
	if err != nil {
		bug.Bugf("gauge %s: %s", g.name, err)
		return
	}
	withLabels.Sub(value)
}

// dummyGauge is a dummy implementation of Gauge.
type dummyGauge struct{}

// Set sets the value of the gauge.
func (g *dummyGauge) Set(value float64, labelValues ...string) {}

// Inc increments the value of the gauge by 1.
func (g *dummyGauge) Inc(labelValues ...string) {}

// Dec decrements the value of the gauge by 1.
func (g *dummyGauge) Dec(labelValues ...string) {}

// Add increments the value of the gauge by the given value.
func (g *dummyGauge) Add(value float64, labelValues ...string) {}

// Sub decrements the value of the gauge by the given value.
func (g *dummyGauge) Sub(value float64, labelValues ...string) {}

// NewGauge creates a new gauge.
func (m *metrics) NewGauge(name, help string, labels ...string) Gauge {
	if len(labels) == 0 {
		return m.newGauge(name, help)
	}
	return m.newGaugeVec(name, help, labels...)
}

// newGauge creates a new gauge.
func (m *metrics) newGauge(name, help string) Gauge {
	if _, ok := m.gauges[name]; ok {
		bug.Bugf("gauge %s already registered", name)
		return &dummyGauge{}
	}
	if _, ok := m.gaugeVecs[name]; ok {
		bug.Bugf("gauge %s already registered as gauge vec", name)
		return &dummyGauge{}
	}
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	})
	if err := m.registry.Register(g); err != nil {
		bug.Bugf("gauge %s: %s", name, err)
		return &dummyGauge{}
	}
	m.gauges[name] = g
	return &gauge{
		name:  name,
		gauge: g,
	}
}

// newGaugeVec returns a new gaugeVec.
func (m *metrics) newGaugeVec(name, help string, labels ...string) Gauge {
	if _, ok := m.gaugeVecs[name]; ok {
		bug.Bugf("gauge vector %s already registered", name)
		return &dummyGauge{}
	}
	if _, ok := m.gauges[name]; ok {
		bug.Bugf("gauge vector %s already registered as gauge", name)
		return &dummyGauge{}
	}
	g := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	}, labels)
	if err := m.registry.Register(g); err != nil {
		bug.Bugf("gauge vector %s: %s", name, err)
		return &dummyGauge{}
	}
	m.gaugeVecs[name] = g
	return &gaugeVec{
		name:     name,
		gaugeVec: g,
	}
}
