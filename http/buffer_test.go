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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_BufferedResponseWriter tests that the BufferedResponseWriter works as
// expected. This means:
//   - The response body is buffered. No writes to the underlying
//     http.ResponseWriter should occur until Flush is called.
//   - The status code is buffered. No writes to the underlying
//     http.ResponseWriter should occur until Flush is called.
//   - The buffered response is flushed to the underlying http.ResponseWriter
//     when Flush is called.
//   - The number of bytes written to the underlying http.ResponseWriter is
//     returned by Flush.
func Test_BufferedResponseWriter(t *testing.T) {
	// Use a ResponseRecorder to capture the response.
	rr := httptest.NewRecorder()
	rr.Code = 0
	bw := NewBufferedResponseWriter(rr)
	// Write some data to the BufferedResponseWriter.
	n, err := bw.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	// There should be no writes to the underlying http.ResponseWriter yet.
	assert.Equal(t, 0, rr.Code)
	assert.Len(t, rr.Header(), 0)
	assert.Equal(t, "", rr.Body.String())
	// Write the status code to the BufferedResponseWriter.
	bw.WriteHeader(http.StatusTeapot)
	// Set a header on the BufferedResponseWriter.
	bw.Header().Set("X-Test", "test")
	// There should be no writes to the underlying http.ResponseWriter yet.
	// The headers should be up to date.
	assert.Equal(t, 0, rr.Code)
	assert.Len(t, rr.Header(), 1)
	assert.Equal(t, "test", rr.Header().Get("X-Test"))
	assert.Equal(t, "", rr.Body.String())
	// Flush the BufferedResponseWriter.
	n, err = bw.Flush()
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	// The response should now be written to the underlying
	// http.ResponseWriter.
	assert.Equal(t, http.StatusTeapot, rr.Code)
	assert.Equal(t, "test", rr.Header().Get("X-Test"))
	assert.Equal(t, "hello", rr.Body.String())
}
