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

// PrometheusRegistry returns the internal prometheus registry. You probably
// don't want to use this. Instead, use the NewCounter, NewGauge, NewHistogram,
// and NewSummary methods. Mu uses this internally for certain purposes.
func PrometheusRegistry(m Metrics) *prometheus.Registry {
	return m.(*metrics).registry
}

// PrometheusCounterVec returns the internal prometheus counter vec. You
// probably don't want to use this. Instead, use the NewCounter method. Mu uses
// this internally for certain purposes.
func PrometheusCounterVec(c Counter) *prometheus.CounterVec {
	return c.(*counterVec).counterVec
}

// PrometheusGaugeVec returns the internal prometheus gauge vec. You probably
// don't want to use this. Instead, use the NewGauge method. Mu uses this
// internally for certain purposes.
func PrometheusGaugeVec(g Gauge) *prometheus.GaugeVec {
	return g.(*gaugeVec).gaugeVec
}

// PrometheusHistogramVec returns the internal prometheus histogram vec. You
// probably don't want to use this. Instead, use the NewHistogram method. Mu
// uses this internally for certain purposes.
func PrometheusHistogramVec(h Histogram) *prometheus.HistogramVec {
	return h.(*histogramVec).histogramVec
}
