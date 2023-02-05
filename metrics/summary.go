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

// Summary is a type that allows you to update a summary.
type Summary interface {
	// Observe records a value.
	Observe(value float64, labelValues ...string)
}

// summary is the default implementation of Summary.
type summary struct {
	// name is the name of the summary.
	name string
	// summary is the prometheus summary.
	summary prometheus.Summary
}

// Observe records a value.
func (h *summary) Observe(value float64, labelValues ...string) {
	if len(labelValues) > 0 {
		bug.Bugf("summary %s has no labels", h.name)
		return
	}
	h.summary.Observe(value)
}

// summaryVec is the default implementation of Summary for summarys with
// labels.
type summaryVec struct {
	// name is the name of the summary.
	name string
	// summaryVec is the prometheus summary vec.
	summaryVec *prometheus.SummaryVec
}

// Observe records a value.
func (h *summaryVec) Observe(value float64, labelValues ...string) {
	withLabels, err := h.summaryVec.GetMetricWithLabelValues(labelValues...)
	if err != nil {
		bug.Bugf("summary %s: %s", h.name, err)
		return
	}
	withLabels.Observe(value)
}

// dummySummary is a dummy implementation of Summary.
type dummySummary struct{}

// Observe records a value.
func (h *dummySummary) Observe(value float64, labelValues ...string) {}

// NewSummary returns a new Summary.
func (m *metrics) NewSummary(name, help string, labels ...string) Summary {
	if len(labels) == 0 {
		return m.newSummary(name, help)
	}
	return m.newSummaryVec(name, help, labels...)
}

// newSummary returns a new Summary.
func (m *metrics) newSummary(name, help string) Summary {
	if _, ok := m.summaries[name]; ok {
		bug.Bugf("summary %s already registered", name)
		return &dummySummary{}
	}
	if _, ok := m.summaryVecs[name]; ok {
		bug.Bugf("summary %s already registered as summary vec", name)
		return &dummySummary{}
	}
	h := prometheus.NewSummary(prometheus.SummaryOpts{
		Name: name,
		Help: help,
	})
	if err := m.registry.Register(h); err != nil {
		bug.Bugf("summary %s: %s", name, err)
		return &dummySummary{}
	}
	m.summaries[name] = h
	return &summary{
		name:    name,
		summary: h,
	}
}

// newSummaryVec returns a new summaryVec.
func (m *metrics) newSummaryVec(name, help string, labels ...string) Summary {
	if _, ok := m.summaryVecs[name]; ok {
		bug.Bugf("summary vector %s already registered", name)
		return &dummySummary{}
	}
	if _, ok := m.summaries[name]; ok {
		bug.Bugf("summary vector %s already registered as summary", name)
		return &dummySummary{}
	}
	h := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: name,
		Help: help,
	}, labels)
	if err := m.registry.Register(h); err != nil {
		bug.Bugf("summary vector %s: %s", name, err)
		return &dummySummary{}
	}
	m.summaryVecs[name] = h
	return &summaryVec{
		name:       name,
		summaryVec: h,
	}
}
