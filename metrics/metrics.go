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
	ht "net/http"

	"github.com/neuralnorthwest/mu/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

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
	// buckets, and labels. The labels are used to create a unique histogram for
	// each label combination. Buckets can be nil, in which case the default
	// buckets are used.
	NewHistogram(name, help string, buckets []float64, labels ...string) Histogram
	// NewSummary registers a summary with the given name, help text, and
	// objectives, and labels. The labels are used to create a unique summary
	// for each label combination. Objectives can be nil, in which case the
	// default objectives are used.
	NewSummary(name, help string, objectives map[float64]float64, labels ...string) Summary
	// Handler returns a ht.Handler that serves the metrics.
	Handler(logger logging.Logger) ht.Handler
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

// Option is a type that can be used to configure a Metrics.
type Option interface {
	apply(*metrics) error
}

// withRegistryOption is an Option that configures the Metrics to use the given
// prometheus registry.
type withRegistryOption struct {
	registry *prometheus.Registry
}

func (opt *withRegistryOption) apply(m *metrics) error {
	m.registry = opt.registry
	return nil
}

// funcOption is an Option that is implemented by a function.
type funcOption struct {
	f func(*metrics) error
}

func (opt *funcOption) apply(m *metrics) error {
	return opt.f(m)
}

// WithRegistry returns an Option that configures the Metrics to use the given
// prometheus registry. If specified multiple times, only the last one is used.
func WithRegistry(registry *prometheus.Registry) Option {
	return &withRegistryOption{registry: registry}
}

// WithCollector returns an Option that configures the Metrics to register the
// given collector. This can be specified multiple times to register multiple
// collectors.
func WithCollector(collector prometheus.Collector) Option {
	return &funcOption{f: func(m *metrics) error {
		return m.registry.Register(collector)
	}}
}

// WithGoCollector returns an Option that configures the Metrics to register
// the Go collector. This should only be specified once.
func WithGoCollector() Option {
	return WithCollector(collectors.NewGoCollector())
}

// New returns a new Metrics.
func New(opts ...Option) (Metrics, error) {
	met := &metrics{
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
	// Apply any options that configure the registry first, then apply the rest.
	for _, applyRegistry := range []bool{true, false} {
		for _, opt := range opts {
			if _, ok := opt.(*withRegistryOption); ok == applyRegistry {
				if err := opt.apply(met); err != nil {
					return nil, err
				}
			}
		}
	}
	return met, nil
}
