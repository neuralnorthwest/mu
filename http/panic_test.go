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
	"github.com/stretchr/testify/assert"
)

// panicHandler is a handler that panics.
func panicHandler(w ht.ResponseWriter, r *ht.Request) {
	panic("test panic")
}

// writeBeforePanicHandler is a handler that writes to the response before
// panicking.
func writeBeforePanicHandler(w ht.ResponseWriter, r *ht.Request) {
	_, _ = w.Write([]byte("hello"))
	panic("test panic")
}

// Test_PanicMiddleware_NonBuffered tests that an unbuffered PanicMiddleware
// recovers from panics and returns a 500 Internal Server Error so long as
// the handler does not write to the response before panicking.
func Test_PanicMiddleware_NonBuffered(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	logger := mock_logging.NewMockLogger(mockCtrl)
	logger.EXPECT().Debugw("HTTP panic stack trace", "stack", gomock.Any())
	// Use a ResponseRecorder to capture the response.
	rr := httptest.NewRecorder()
	rr.Code = 0
	// Create a panic middleware.
	middleware := PanicMiddleware(logger, false)
	// Create a request.
	req, err := ht.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	// Call the middleware.
	middleware("/panic", ht.HandlerFunc(panicHandler)).ServeHTTP(rr, req)
	// The response should be a 500.
	assert.Equal(t, ht.StatusInternalServerError, rr.Code)
	// The response should contain the panic message.
	assert.Contains(t, rr.Body.String(), "Internal Server Error\n")
}

// Test_PanicMiddleware_NonBuffer_WriteBeforePanic tests that an unbuffered
// PanicMiddleware recovers from panics but does not return a 500 Internal
// Server Error if the handler writes to the response before panicking.
func Test_PanicMiddleware_NonBuffer_WriteBeforePanic(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	logger := mock_logging.NewMockLogger(mockCtrl)
	logger.EXPECT().Debugw("HTTP panic stack trace", "stack", gomock.Any())
	// Use a ResponseRecorder to capture the response.
	rr := httptest.NewRecorder()
	rr.Code = 0
	// Create a panic middleware.
	middleware := PanicMiddleware(logger, false)
	// Create a request.
	req, err := ht.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	// Call the middleware.
	middleware("/panic", ht.HandlerFunc(writeBeforePanicHandler)).ServeHTTP(rr, req)
	// The response should be a 200.
	assert.Equal(t, ht.StatusOK, rr.Code)
	// The response should contain the handler's response.
	assert.Contains(t, rr.Body.String(), "hello")
}

// Test_PanicMiddleware_Buffered tests that a buffered PanicMiddleware
// recovers from panics and returns a 500 Internal Server Error even when
// the handler writes to the response before panicking.
func Test_PanicMiddleware_Buffered(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	logger := mock_logging.NewMockLogger(mockCtrl)
	logger.EXPECT().Debugw("HTTP panic stack trace", "stack", gomock.Any())
	// Use a ResponseRecorder to capture the response.
	rr := httptest.NewRecorder()
	rr.Code = 0
	// Create a panic middleware.
	middleware := PanicMiddleware(logger, true)
	// Create a request.
	req, err := ht.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	// Call the middleware.
	middleware("/panic", ht.HandlerFunc(writeBeforePanicHandler)).ServeHTTP(rr, req)
	// The response should be a 500.
	assert.Equal(t, ht.StatusInternalServerError, rr.Code)
	// The response should contain the panic message.
	assert.Contains(t, rr.Body.String(), "Internal Server Error\n")
}
