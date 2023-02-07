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
)

// errorLoggingInterceptor is a ht.ResponseWriter that intercepts the status
// code and body of the response.
type errorLoggingInterceptor struct {
	ht.ResponseWriter
	statusCode int
	body       []byte
}

// WriteHeader intercepts the status code.
func (i *errorLoggingInterceptor) WriteHeader(statusCode int) {
	i.statusCode = statusCode
	i.ResponseWriter.WriteHeader(statusCode)
}

// Write intercepts the body.
func (i *errorLoggingInterceptor) Write(body []byte) (int, error) {
	// If the status code is >= 500, then the body is an error message.
	// Capture the first 512 bytes of the body.
	if i.statusCode >= 500 {
		if len(body) <= 512 {
			i.body = body
		} else {

			//i.body = append(body[:512], []byte("...")...)
			i.body = make([]byte, 515)
			copy(i.body, body[:512])
			copy(i.body[512:], []byte("..."))
		}
	}
	return i.ResponseWriter.Write(body)
}

// ErrorLoggingMiddleware is an HTTP middleware that logs errors.
func ErrorLoggingMiddleware(logger logging.Logger) Middleware {
	return func(pattern string, next ht.Handler) ht.Handler {
		return ht.HandlerFunc(func(w ht.ResponseWriter, r *ht.Request) {
			i := &errorLoggingInterceptor{
				ResponseWriter: w,
			}
			next.ServeHTTP(i, r)
			if i.statusCode >= 500 {
				logger.Errorw("HTTP error", "status_code", i.statusCode, "path", pattern, "body", string(i.body))
			}
		})
	}
}

// WithErrorLogging returns an option that adds the error logging middleware to
// the server.
func WithErrorLogging(logger logging.Logger) ServerOption {
	return WithMiddleware(ErrorLoggingMiddleware(logger))
}
