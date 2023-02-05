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

import "github.com/prometheus/client_golang/prometheus"

// Metrics is a type that allows you to register and report metrics.
type Metrics interface {
	// NewCounter registers a counter with the given name, help text, and
	// labels. The labels are used to create a unique counter for each label
	// combination.
	NewCounter(name, help string, labels ...string) Counter
	// NewGauge registers a gauge with the given name, help text, and
	// labels. The labels are used to create a unique gauge for each label
	// combination.
	NewGauge(name, help string, labels ...string) Gauge
	// NewHistogram registers a histogram with the given name, help text,
	// and labels. The labels are used to create a unique histogram for each
	// label combination.
	NewHistogram(name, help string, labels ...string) Histogram
	// NewSummary registers a summary with the given name, help text, and
	// labels. The labels are used to create a unique summary for each label
	// combination.
	NewSummary(name, help string, labels ...string) Summary
}

// metrics is the default implementation of Metrics.
type metrics struct {
	// registry is the prometheus registry.
	registry *prometheus.Registry
	// counters is a map of counter names to counters.
	counters map[string]prometheus.Counter
	// counterVecs is a map of counter names to counter vecs.
	counterVecs map[string]*prometheus.CounterVec
	// gauges is a map of gauge names to gauges.
	gauges map[string]prometheus.Gauge
	// gaugeVecs is a map of gauge names to gauge vecs.
	gaugeVecs map[string]*prometheus.GaugeVec
	// histograms is a map of histogram names to histograms.
	histograms map[string]prometheus.Histogram
	// histogramVecs is a map of histogram names to histogram vecs.
	histogramVecs map[string]*prometheus.HistogramVec
	// summaries is a map of summary names to summaries.
	summaries map[string]prometheus.Summary
	// summaryVecs is a map of summary names to summary vecs.
	summaryVecs map[string]*prometheus.SummaryVec
}

// New returns a new Metrics.
func New() Metrics {
	return &metrics{
		registry:      prometheus.NewRegistry(),
		counters:      make(map[string]prometheus.Counter),
		counterVecs:   make(map[string]*prometheus.CounterVec),
		gauges:        make(map[string]prometheus.Gauge),
		gaugeVecs:     make(map[string]*prometheus.GaugeVec),
		histograms:    make(map[string]prometheus.Histogram),
		histogramVecs: make(map[string]*prometheus.HistogramVec),
		summaries:     make(map[string]prometheus.Summary),
		summaryVecs:   make(map[string]*prometheus.SummaryVec),
	}
}
