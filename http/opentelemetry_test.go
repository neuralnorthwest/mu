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
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	ht "net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Test_OpenTelemetryTracing tests the OpenTelemetry tracing middleware. We
// use stddouttrace to capture the trace output and then verify that it is
// correct.
func Test_OpenTelemetryTracing(t *testing.T) {
	t.Parallel()

	// Initialize tracing to an in-memory buffer.
	tracebuf := &bytes.Buffer{}
	exporter, err := stdouttrace.New(
		stdouttrace.WithWriter(tracebuf),
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		t.Fatalf("failed to create exporter: %s", err)
	}
	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)
	shutdown := false
	defer func() {
		if !shutdown {
			err := provider.Shutdown(context.Background())
			if err != nil {
				t.Fatalf("failed to shutdown provider: %s", err)
			}
		}
	}()
	otel.SetTracerProvider(provider)

	middleware := OpenTelemetryTracingMiddleware()
	handler := middleware("/test", ht.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()
	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, rr.Result().StatusCode)
	}
	if rr.Result().Body == nil {
		t.Fatalf("expected a response body, got nil")
	}
	body, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatalf("failed to read response body: %s", err)
	}
	if string(body) != "OK" {
		t.Fatalf("expected response body %q, got %q", "OK", string(body))
	}
	err = provider.Shutdown(context.Background())
	if err != nil {
		t.Fatalf("failed to shutdown provider: %s", err)
	}
	shutdown = true
	var tracedata []map[string]interface{}
	//t.Fatalf("trace: %s", tracebuf.String())
	decoder := json.NewDecoder(tracebuf)
	for {
		var data map[string]interface{}
		err := decoder.Decode(&data)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("failed to decode trace data: %s", err)
		}
		tracedata = append(tracedata, data)
	}
	assert.Len(t, tracedata, 1)
	trace := tracedata[0]
	assert.Equal(t, "handle", trace["Name"], "expected trace name %q, got %q", "handle", trace["Name"])
	gotRoute := false
	gotStatusCode := false
	gotWroteBytes := false
	for _, attr := range trace["Attributes"].([]interface{}) {
		attrMap := attr.(map[string]interface{})
		if attrMap["Key"] == "http.route" {
			valueMap := attrMap["Value"].(map[string]interface{})
			assert.Equal(t, "STRING", valueMap["Type"].(string))
			assert.Equal(t, "/test", valueMap["Value"].(string))
			gotRoute = true
		}
		if attrMap["Key"] == "http.status_code" {
			valueMap := attrMap["Value"].(map[string]interface{})
			assert.Equal(t, "INT64", valueMap["Type"].(string))
			assert.Equal(t, float64(http.StatusOK), valueMap["Value"].(float64))
			gotStatusCode = true
		}
		if attrMap["Key"] == "http.wrote_bytes" {
			valueMap := attrMap["Value"].(map[string]interface{})
			assert.Equal(t, "INT64", valueMap["Type"].(string))
			assert.Equal(t, float64(2), valueMap["Value"].(float64))
			gotWroteBytes = true
		}
	}
	assert.True(t, gotRoute)
	assert.True(t, gotStatusCode)
	assert.True(t, gotWroteBytes)
}
