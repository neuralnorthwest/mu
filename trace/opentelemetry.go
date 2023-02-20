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

package trace

import (
	"context"
	"fmt"
	"time"

	"github.com/neuralnorthwest/mu"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// otelOptions is a group of options for a tracer.
type otelOptions struct {
	baseOptions
	// endpoint is the endpoint to which spans should be sent.
	endpoint string
	// insec is a flag that indicates whether the collector endpoint is insecure.
	insec bool
}

// collectorEndpoint sets the collector endpoint.
func (o *otelOptions) collectorEndpoint(endpoint string) {
	o.endpoint = endpoint
}

// insecure sets the insecure flag.
func (o *otelOptions) insecure() {
	o.insec = true
}

// OpenTelemetryTracing initializes OpenTelemetry tracing. It returns a tracer
// and a cleanup function that should be called when tracing is no longer
// needed.
func OpenTelemetryTracing(ctx context.Context, name, version string, opts ...TracerOption) (Tracer, func(), error) {
	var o otelOptions
	for _, opt := range opts {
		opt.apply(&o)
	}
	otelopts := []otlptracehttp.Option{}
	if o.endpoint != "" {
		otelopts = append(otelopts, otlptracehttp.WithEndpoint(o.endpoint))
	}
	if o.insec {
		otelopts = append(otelopts, otlptracehttp.WithInsecure())
	}
	client := otlptracehttp.NewClient(otelopts...)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, nil, fmt.Errorf("could not initialize OTLP trace exporter: %w", err)
	}
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(name),
			semconv.ServiceVersion(version),
		)),
	)
	otel.SetTracerProvider(tracerProvider)
	tracer := otel.Tracer("github.com/neuralnorthwest/mu",
		trace.WithInstrumentationVersion(mu.Version()),
		trace.WithSchemaURL(semconv.SchemaURL),
	)
	return tracer, func() {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_ = tracerProvider.Shutdown(ctx)
	}, nil
}
