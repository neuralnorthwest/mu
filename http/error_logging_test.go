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
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mock_logging "github.com/neuralnorthwest/mu/logging/mock"
)

// Test_ErrorLoggingMiddleware_Case is a test case for Test_ErrorLoggingMiddleware.
type Test_ErrorLoggingMiddleware_Case struct {
	name        string
	status      int
	body        string
	wantLogBody string
}

// longString is a long string used to test that the error logging middleware
// truncates the body.
var longString = func() string {
	var s string
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	for i := 0; i < 513; i++ {
		s += string(alphabet[i%26])
	}
	return s
}()

// Test_ErrorLoggingMiddleware tests that the error logging middleware logs
// errors.
func Test_ErrorLoggingMiddleware(t *testing.T) {
	t.Parallel()
	testCases := []Test_ErrorLoggingMiddleware_Case{
		{
			name:        "success",
			status:      200,
			body:        "success",
			wantLogBody: "",
		},
		{
			name:        "success long body",
			status:      200,
			body:        longString,
			wantLogBody: "",
		},
		{
			name:        "client error",
			status:      400,
			body:        "client error",
			wantLogBody: "",
		},
		{
			name:        "client error long body",
			status:      400,
			body:        longString,
			wantLogBody: "",
		},
		{
			name:        "server error",
			status:      500,
			body:        "server error",
			wantLogBody: "server error",
		},
		{
			name:        "server error long body",
			status:      500,
			body:        longString,
			wantLogBody: longString[:512] + "...",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			logger := mock_logging.NewMockLogger(mockCtrl)
			if tc.wantLogBody != "" {
				logger.EXPECT().Errorw("HTTP error", "status_code", tc.status, "path", "/test", "body", tc.wantLogBody)
			}
			middleware := ErrorLoggingMiddleware(logger)
			handler := middleware("/test", ht.HandlerFunc(func(w ht.ResponseWriter, r *ht.Request) {
				w.WriteHeader(tc.status)
				_, err := w.Write([]byte(tc.body))
				if err != nil {
					t.Fatal(err)
				}
			}))
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, httptest.NewRequest("GET", "/test", nil))
			// Check that the response is correct.
			if got, want := rr.Code, tc.status; got != want {
				t.Errorf("got %d, want %d", got, want)
			}
			if got, want := rr.Body.String(), tc.body; got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}
