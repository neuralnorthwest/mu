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

// Histogram is a type that allows you to update a histogram.
type Histogram interface {
	// Observe records a value.
	Observe(value float64, labelValues ...string)
}

// histogram is the default implementation of Histogram.
type histogram struct {
	// name is the name of the histogram.
	name string
	// histogram is the prometheus histogram.
	histogram prometheus.Histogram
}

// Observe records a value.
func (h *histogram) Observe(value float64, labelValues ...string) {
	if len(labelValues) > 0 {
		bug.Bugf("histogram %s has no labels", h.name)
		return
	}
	h.histogram.Observe(value)
}

// histogramVec is the default implementation of Histogram for histograms with
// labels.
type histogramVec struct {
	// name is the name of the histogram.
	name string
	// histogramVec is the prometheus histogram vec.
	histogramVec *prometheus.HistogramVec
}

// Observe records a value.
func (h *histogramVec) Observe(value float64, labelValues ...string) {
	withLabels, err := h.histogramVec.GetMetricWithLabelValues(labelValues...)
	if err != nil {
		bug.Bugf("histogram %s: %s", h.name, err)
		return
	}
	withLabels.Observe(value)
}

// dummyHistogram is a dummy implementation of Histogram.
type dummyHistogram struct{}

// Observe records a value.
func (h *dummyHistogram) Observe(value float64, labelValues ...string) {}

// NewHistogram returns a new Histogram.
func (m *metrics) NewHistogram(name, help string, buckets []float64, labels ...string) Histogram {
	if len(labels) == 0 {
		return m.newHistogram(name, help, buckets)
	}
	return m.newHistogramVec(name, help, buckets, labels...)
}

// newHistogram returns a new Histogram.
func (m *metrics) newHistogram(name, help string, buckets []float64) Histogram {
	if _, ok := m.histograms[name]; ok {
		bug.Bugf("histogram %s already registered", name)
		return &dummyHistogram{}
	}
	if _, ok := m.histogramVecs[name]; ok {
		bug.Bugf("histogram %s already registered as histogram vec", name)
		return &dummyHistogram{}
	}
	h := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    name,
		Help:    help,
		Buckets: buckets,
	})
	if err := m.registry.Register(h); err != nil {
		bug.Bugf("histogram %s: %s", name, err)
		return &dummyHistogram{}
	}
	m.histograms[name] = h
	return &histogram{
		name:      name,
		histogram: h,
	}
}

// newHistogramVec returns a new histogramVec.
func (m *metrics) newHistogramVec(name, help string, buckets []float64, labels ...string) Histogram {
	if _, ok := m.histogramVecs[name]; ok {
		bug.Bugf("histogram vector %s already registered", name)
		return &dummyHistogram{}
	}
	if _, ok := m.histograms[name]; ok {
		bug.Bugf("histogram vector %s already registered as histogram", name)
		return &dummyHistogram{}
	}
	h := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    name,
		Help:    help,
		Buckets: buckets,
	}, labels)
	if err := m.registry.Register(h); err != nil {
		bug.Bugf("histogram vector %s: %s", name, err)
		return &dummyHistogram{}
	}
	m.histogramVecs[name] = h
	return &histogramVec{
		name:         name,
		histogramVec: h,
	}
}
