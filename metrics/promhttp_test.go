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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mock_logging "github.com/neuralnorthwest/mu/logging/mock"
	"github.com/prometheus/common/expfmt"
)

// Test_Prometheus_Metrics_Promhttp_Handler tests that the Prometheus metrics
// are correctly served by the promhttp handler.
func Test_Prometheus_Metrics_Promhttp_Handler(t *testing.T) {
	t.Parallel()
	m := newMetrics(t, WithGoCollector())
	cinc := m.NewCounter("counter_test_inc", "counter_test_inc")
	cadd := m.NewCounter("counter_test_add", "counter_test_add")
	gset := m.NewGauge("gauge_test_set", "gauge_test_set")
	ginc := m.NewGauge("gauge_test_inc", "gauge_test_inc")
	gdec := m.NewGauge("gauge_test_dec", "gauge_test_dec")
	gadd := m.NewGauge("gauge_test_add", "gauge_test_add")
	gsub := m.NewGauge("gauge_test_sub", "gauge_test_sub")
	hobs := m.NewHistogram("histogram_test_observe", "histogram_test_observe", nil)
	sobs := m.NewSummary("summary_test_observe", "summary_test_observe", nil)
	clinc := m.NewCounter("counter_test_inc_label", "counter_test_inc", "label")
	cladd := m.NewCounter("counter_test_add_label", "counter_test_add", "label")
	glset := m.NewGauge("gauge_test_set_label", "gauge_test_set", "label")
	glinc := m.NewGauge("gauge_test_inc_label", "gauge_test_inc", "label")
	gldec := m.NewGauge("gauge_test_dec_label", "gauge_test_dec", "label")
	gladd := m.NewGauge("gauge_test_add_label", "gauge_test_add", "label")
	glsub := m.NewGauge("gauge_test_sub_label", "gauge_test_sub", "label")
	hlobs := m.NewHistogram("histogram_test_observe_label", "histogram_test_observe", nil, "label")
	slobs := m.NewSummary("summary_test_observe_label", "summary_test_observe", nil, "label")

	cinc.Inc()
	cadd.Add(1)
	gset.Set(1)
	ginc.Inc()
	gdec.Dec()
	gadd.Add(1)
	gsub.Sub(1)
	hobs.Observe(1)
	sobs.Observe(1)
	clinc.Inc("label")
	cladd.Add(1, "label")
	glset.Set(1, "label")
	glinc.Inc("label")
	gldec.Dec("label")
	gladd.Add(1, "label")
	glsub.Sub(1, "label")
	hlobs.Observe(1, "label")
	slobs.Observe(1, "label")

	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	logger := mock_logging.NewMockLogger(mockCtrl)
	m.Handler(logger).ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
	resp := w.Result()
	defer resp.Body.Close()
	if resp.Header.Get("Content-Type") != "text/plain; version=0.0.4; charset=utf-8" {
		t.Fatalf("expected Content-Type %s, got %s", "text/plain; version=0.0.4; charset=utf-8", resp.Header.Get("Content-Type"))
	}
	parser := expfmt.TextParser{}
	metrics, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		t.Fatalf("failed to parse metrics: %s", err)
	}
	// 45 total == 18 + 27 (go metrics)
	if len(metrics) != 18+27 {
		t.Fatalf("expected 18 + 27 metrics, got %d", len(metrics))
	}
	if metrics["counter_test_inc"].GetMetric()[0].GetCounter().GetValue() != 1 {
		t.Fatalf("expected counter_test_inc to be 1, got %f", metrics["counter_test_inc"].GetMetric()[0].GetCounter().GetValue())
	}
	if metrics["counter_test_add"].GetMetric()[0].GetCounter().GetValue() != 1 {
		t.Fatalf("expected counter_test_add to be 1, got %f", metrics["counter_test_add"].GetMetric()[0].GetCounter().GetValue())
	}
	if metrics["gauge_test_set"].GetMetric()[0].GetGauge().GetValue() != 1 {
		t.Fatalf("expected gauge_test_set to be 1, got %f", metrics["gauge_test_set"].GetMetric()[0].GetGauge().GetValue())
	}
	if metrics["gauge_test_inc"].GetMetric()[0].GetGauge().GetValue() != 1 {
		t.Fatalf("expected gauge_test_inc to be 1, got %f", metrics["gauge_test_inc"].GetMetric()[0].GetGauge().GetValue())
	}
	if metrics["gauge_test_dec"].GetMetric()[0].GetGauge().GetValue() != -1 {
		t.Fatalf("expected gauge_test_dec to be -1, got %f", metrics["gauge_test_dec"].GetMetric()[0].GetGauge().GetValue())
	}
	if metrics["gauge_test_add"].GetMetric()[0].GetGauge().GetValue() != 1 {
		t.Fatalf("expected gauge_test_add to be 1, got %f", metrics["gauge_test_add"].GetMetric()[0].GetGauge().GetValue())
	}
	if metrics["gauge_test_sub"].GetMetric()[0].GetGauge().GetValue() != -1 {
		t.Fatalf("expected gauge_test_sub to be -1, got %f", metrics["gauge_test_sub"].GetMetric()[0].GetGauge().GetValue())
	}
	if metrics["histogram_test_observe"].GetMetric()[0].GetHistogram().GetSampleCount() != 1 {
		t.Fatalf("expected histogram_test_observe to be 1, got %d", metrics["histogram_test_observe"].GetMetric()[0].GetHistogram().GetSampleCount())
	}
	if metrics["summary_test_observe"].GetMetric()[0].GetSummary().GetSampleCount() != 1 {
		t.Fatalf("expected summary_test_observe to be 1, got %d", metrics["summary_test_observe"].GetMetric()[0].GetSummary().GetSampleCount())
	}
	if metrics["counter_test_inc_label"].GetMetric()[0].GetCounter().GetValue() != 1 {
		t.Fatalf("expected counter_test_inc_label to be 1, got %f", metrics["counter_test_inc_label"].GetMetric()[0].GetCounter().GetValue())
	}
	if metrics["counter_test_add_label"].GetMetric()[0].GetCounter().GetValue() != 1 {
		t.Fatalf("expected counter_test_add_label to be 1, got %f", metrics["counter_test_add_label"].GetMetric()[0].GetCounter().GetValue())
	}
	if metrics["gauge_test_set_label"].GetMetric()[0].GetGauge().GetValue() != 1 {
		t.Fatalf("expected gauge_test_set_label to be 1, got %f", metrics["gauge_test_set_label"].GetMetric()[0].GetGauge().GetValue())
	}
	if metrics["gauge_test_inc_label"].GetMetric()[0].GetGauge().GetValue() != 1 {
		t.Fatalf("expected gauge_test_inc_label to be 1, got %f", metrics["gauge_test_inc_label"].GetMetric()[0].GetGauge().GetValue())
	}
	if metrics["gauge_test_dec_label"].GetMetric()[0].GetGauge().GetValue() != -1 {
		t.Fatalf("expected gauge_test_dec_label to be -1, got %f", metrics["gauge_test_dec_label"].GetMetric()[0].GetGauge().GetValue())
	}
	if metrics["gauge_test_add_label"].GetMetric()[0].GetGauge().GetValue() != 1 {
		t.Fatalf("expected gauge_test_add_label to be 1, got %f", metrics["gauge_test_add_label"].GetMetric()[0].GetGauge().GetValue())
	}
	if metrics["gauge_test_sub_label"].GetMetric()[0].GetGauge().GetValue() != -1 {
		t.Fatalf("expected gauge_test_sub_label to be -1, got %f", metrics["gauge_test_sub_label"].GetMetric()[0].GetGauge().GetValue())
	}
	if metrics["histogram_test_observe_label"].GetMetric()[0].GetHistogram().GetSampleCount() != 1 {
		t.Fatalf("expected histogram_test_observe_label to be 1, got %d", metrics["histogram_test_observe_label"].GetMetric()[0].GetHistogram().GetSampleCount())
	}
	if metrics["summary_test_observe_label"].GetMetric()[0].GetSummary().GetSampleCount() != 1 {
		t.Fatalf("expected summary_test_observe_label to be 1, got %d", metrics["summary_test_observe_label"].GetMetric()[0].GetSummary().GetSampleCount())
	}
}
