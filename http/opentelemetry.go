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

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// OpenTelemetryTracingMiddleware is an HTTP middleware that adds OpenTelemetry
// tracing to the server.
func OpenTelemetryTracingMiddleware(opts ...otelhttp.Option) Middleware {
	return func(pattern string, next ht.Handler) ht.Handler {
		return otelhttp.NewHandler(otelhttp.WithRouteTag(pattern, next), "handle", opts...)
	}
}

// WithOpenTelemetryTracing returns an HTTP middleware that adds OpenTelemetry
// tracing to the server.
func WithOpenTelemetryTracing(opts ...otelhttp.Option) ServerOption {
	return WithMiddleware(OpenTelemetryTracingMiddleware(opts...))
}
