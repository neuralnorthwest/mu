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

// Counter is a type that allows you to increment a counter.
type Counter interface {
	// Inc increments the counter by 1.
	Inc(labelValues ...string)
	// Add increments the counter by the given value.
	Add(value float64, labelValues ...string)
}

// counter is the default implementation of Counter.
type counter struct {
	// name is the name of the counter.
	name string
	// counter is the prometheus counter.
	counter prometheus.Counter
}

// Inc increments the counter by 1.
func (c *counter) Inc(labelValues ...string) {
	if len(labelValues) > 0 {
		bug.Bugf("counter %s has no labels", c.name)
		return
	}
	c.counter.Inc()
}

// Add increments the counter by the given value.
func (c *counter) Add(value float64, labelValues ...string) {
	if len(labelValues) > 0 {
		bug.Bugf("counter %s has no labels", c.name)
		return
	}
	c.counter.Add(value)
}

// counterVec is the default implementation of Counter for counters with labels.
type counterVec struct {
	// name is the name of the counter.
	name string
	// counterVec is the prometheus counter vec.
	counterVec *prometheus.CounterVec
}

// Inc increments the counter by 1.
func (c *counterVec) Inc(labelValues ...string) {
	withLabels, err := c.counterVec.GetMetricWithLabelValues(labelValues...)
	if err != nil {
		bug.Bugf("counter %s: %s", c.name, err)
		return
	}
	withLabels.Inc()
}

// Add increments the counter by the given value.
func (c *counterVec) Add(value float64, labelValues ...string) {
	withLabels, err := c.counterVec.GetMetricWithLabelValues(labelValues...)
	if err != nil {
		bug.Bugf("counter %s: %s", c.name, err)
		return
	}
	withLabels.Add(value)
}

// dummyCounter is a dummy implementation of Counter.
type dummyCounter struct{}

// Inc increments the counter by 1.
func (c *dummyCounter) Inc(labelValues ...string) {}

// Add increments the counter by the given value.
func (c *dummyCounter) Add(value float64, labelValues ...string) {}

// NewCounter returns a new Counter. Is uses Register, not MustRegister, so that
// it can call bug.Bugf if there is an error.
func (m *metrics) NewCounter(name, help string, labels ...string) Counter {
	if len(labels) == 0 {
		return m.newCounter(name, help)
	}
	return m.newCounterVec(name, help, labels...)
}

// newCounter returns a new Counter.
func (m *metrics) newCounter(name, help string) Counter {
	if _, ok := m.counters[name]; ok {
		bug.Bugf("counter %s already registered", name)
		return &dummyCounter{}
	}
	if _, ok := m.counterVecs[name]; ok {
		bug.Bugf("counter %s already registered as counter vec", name)
		return &dummyCounter{}
	}
	c := prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})
	if err := m.registry.Register(c); err != nil {
		bug.Bugf("counter %s: %s", name, err)
		return &dummyCounter{}
	}
	m.counters[name] = c
	return &counter{
		name:    name,
		counter: c,
	}
}

// newCounterVec returns a new counterVec.
func (m *metrics) newCounterVec(name, help string, labels ...string) Counter {
	if _, ok := m.counterVecs[name]; ok {
		bug.Bugf("counter vector %s already registered", name)
		return &dummyCounter{}
	}
	if _, ok := m.counters[name]; ok {
		bug.Bugf("counter vector %s already registered as counter", name)
		return &dummyCounter{}
	}
	c := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: help,
	}, labels)
	if err := m.registry.Register(c); err != nil {
		bug.Bugf("counter vector %s: %s", name, err)
		return &dummyCounter{}
	}
	m.counterVecs[name] = c
	return &counterVec{
		name:       name,
		counterVec: c,
	}
}
