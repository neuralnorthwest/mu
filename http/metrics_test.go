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
	"math/rand"
	ht "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	logging_mock "github.com/neuralnorthwest/mu/logging/mock"
	"github.com/neuralnorthwest/mu/metrics"
	"github.com/prometheus/common/expfmt"
)

// Test_MetricsMiddleware_Case is a test case for Test_MetricsMiddleware.
type Test_MetricsMiddleware_Case struct {
	name string
	opts MetricsOptions
}

// Test_MetricsMiddleware tests that the metrics middleware records metrics.
func Test_MetricsMiddleware(t *testing.T) {
	t.Parallel()
	testCases := []Test_MetricsMiddleware_Case{
		{
			name: "default",
			opts: MetricsOptions{},
		},
		{
			name: "with buckets",
			opts: MetricsOptions{
				RequestDurationBuckets:   []float64{0.1, 0.2, 0.3},
				RequestSizeBuckets:       []float64{100, 200, 300},
				ResponseSizeBuckets:      []float64{100, 200, 300},
				TimeToWriteHeaderBuckets: []float64{0.1, 0.2, 0.3},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			met, err := metrics.New()
			if err != nil {
				t.Fatalf("metrics.New() = %v, want nil", err)
			}
			middleware := MetricsMiddleware(met, tc.opts)
			for _, responseCode := range []int{200, 400, 500} {
				func() {
					handler := middleware("/test", ht.HandlerFunc(func(w ht.ResponseWriter, r *ht.Request) {
						// timeToHeader is a normally distributed random variable with mean 0.2
						// and standard deviation 0.1.
						timeToHeader := 0.2 + 0.1*rand.NormFloat64()
						if timeToHeader < 0.001 {
							timeToHeader = 0.001
						}
						// timeToRespond is a normally distributed random variable with mean 0.2
						// and standard deviation 0.1.
						timeToRespond := 0.2 + 0.1*rand.NormFloat64()
						if timeToRespond < 0.001 {
							timeToRespond = 0.001
						}
						// responseSize is a normally distributed random variable with mean 200
						// and standard deviation 100.
						responseSize := 200 + 100*rand.NormFloat64()
						if responseSize < 0 {
							responseSize = 0
						}
						// wait for timeToHeader
						time.Sleep(time.Duration(timeToHeader * float64(time.Second)))
						// write the header
						w.WriteHeader(responseCode)
						// wait for timeToRespond
						time.Sleep(time.Duration(timeToRespond * float64(time.Second)))
						// write the body
						_, _ = w.Write(make([]byte, int(responseSize)))
					}))
					// requestSize is a normally distributed random variable with mean 200
					// and standard deviation 100.
					requestSize := 200 + 100*rand.NormFloat64()
					if requestSize < 0 {
						requestSize = 0
					}
					// make the request
					req := httptest.NewRequest("GET", "/test", bytes.NewBuffer(make([]byte, int(requestSize))))
					rr := httptest.NewRecorder()
					handler.ServeHTTP(rr, req)
					defer rr.Result().Body.Close()
					// check the response
					if rr.Code != responseCode {
						t.Fatalf("expected response code %d, got %d", responseCode, rr.Code)
					}
					// check the metrics
					mockCtrl := gomock.NewController(t)
					defer mockCtrl.Finish()
					logger := logging_mock.NewMockLogger(mockCtrl)
					handler = met.Handler(logger)
					req = httptest.NewRequest("GET", "/metrics", nil)
					rr = httptest.NewRecorder()
					handler.ServeHTTP(rr, req)
					defer rr.Result().Body.Close()
					// parse the metrics
					parser := expfmt.TextParser{}
					metrics, err := parser.TextToMetricFamilies(rr.Result().Body)
					if err != nil {
						t.Fatalf("failed to parse metrics: %s", err)
					}
					if len(metrics) != 6 {
						t.Fatalf("expected 6 metrics, got %d", len(metrics))
					}
					// check that each metric is present
					for _, name := range []string{
						"http_request_duration_seconds",
						"http_request_size_bytes",
						"http_response_size_bytes",
						"http_requests_total",
						"http_requests_in_progress",
						"http_time_to_write_header_seconds",
					} {
						if _, ok := metrics[name]; !ok {
							t.Fatalf("expected metric %s, got none", name)
						}
					}
				}()
			}
		})
	}
}
