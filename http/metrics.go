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

package http

import (
	ht "net/http"

	"github.com/neuralnorthwest/mu/logging"
	"github.com/neuralnorthwest/mu/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsOptions specifies options for HTTP metrics
type MetricsOptions struct {
	// Server is the server to use for the metrics endpoint. If nil, the
	// server that is being constructed will be used.
	Server *Server
	// Path is the path to the metrics endpoint. If empty, the metrics
	// are served at /metrics.
	Path string
	// RequestDurationBuckets are the buckets for the request duration
	// histogram. If empty, the default buckets are used.
	RequestDurationBuckets []float64
	// RequestSizeBuckets are the buckets for the request size histogram.
	// If empty, the default buckets are used.
	RequestSizeBuckets []float64
	// ResponseSizeBuckets are the buckets for the response size
	// histogram. If empty, the default buckets are used.
	ResponseSizeBuckets []float64
	// TimeToWriteHeaderBuckets are the buckets for the time to write the
	// response header histogram. If empty, the default buckets are used.
	TimeToWriteHeaderBuckets []float64
}

// MetricsMiddleware returns an HTTP middleware that adds HTTP request metrics
// to the server.
func MetricsMiddleware(met metrics.Metrics, opts MetricsOptions) Middleware {
	requestsTotal := met.NewCounter("http_requests_total", "The total number of HTTP requests.", "method", "code", "path")
	requestsInProgress := met.NewGauge("http_requests_in_progress", "The number of HTTP requests currently in progress.", "path")
	requestDurations := met.NewHistogram("http_request_duration_seconds", "The HTTP request latencies in seconds.", opts.RequestDurationBuckets, "method", "code", "path")
	requestSizes := met.NewHistogram("http_request_size_bytes", "The HTTP request sizes in bytes.", opts.RequestSizeBuckets, "method", "code", "path")
	responseSizes := met.NewHistogram("http_response_size_bytes", "The HTTP response sizes in bytes.", opts.ResponseSizeBuckets, "method", "code", "path")
	timeToWriteHeader := met.NewHistogram("http_time_to_write_header_seconds", "The time to write the HTTP response header in seconds.", opts.TimeToWriteHeaderBuckets, "method", "code", "path")
	return func(pattern string, next ht.Handler) ht.Handler {
		pathLabel := prometheus.Labels{"path": pattern}
		var h ht.Handler
		h = promhttp.InstrumentHandlerCounter(metrics.PrometheusCounterVec(requestsTotal).MustCurryWith(pathLabel), next)
		h = promhttp.InstrumentHandlerInFlight(metrics.PrometheusGaugeVec(requestsInProgress).With(pathLabel), h)
		h = promhttp.InstrumentHandlerDuration(metrics.PrometheusHistogramVec(requestDurations).MustCurryWith(pathLabel), h)
		h = promhttp.InstrumentHandlerRequestSize(metrics.PrometheusHistogramVec(requestSizes).MustCurryWith(pathLabel), h)
		h = promhttp.InstrumentHandlerResponseSize(metrics.PrometheusHistogramVec(responseSizes).MustCurryWith(pathLabel), h)
		h = promhttp.InstrumentHandlerTimeToWriteHeader(metrics.PrometheusHistogramVec(timeToWriteHeader).MustCurryWith(pathLabel), h)
		return h
	}
}

// WithMetrics returns a ServerOption that adds HTTP request metrics to
// the server.
func WithMetrics(met metrics.Metrics, logger logging.Logger, opts MetricsOptions) ServerOption {
	middleware := MetricsMiddleware(met, opts)
	return func(server *Server) error {
		err := WithMiddleware(middleware)(server)
		if err != nil {
			return err
		}
		path := opts.Path
		if path == "" {
			path = "/metrics"
		}
		if opts.Server != nil {
			server = opts.Server
		}
		server.Handle(path, promhttp.InstrumentMetricHandler(metrics.PrometheusRegistry(met), met.Handler(logger)))
		return nil
	}
}
