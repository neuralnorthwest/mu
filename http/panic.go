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
	ht "net/http"
	"runtime/debug"

	"github.com/neuralnorthwest/mu/logging"
)

// PanicMiddleware returns an HTTP middleware that recovers from panics and
// returns a 500 Internal Server Error. If buffered is true, the middleware will
// buffer the response and flush it only after the handler has succeeded. This
// is useful if the handler writes to the response before panicking. If you
// know that the handler will *not* write to the response before panicking, then
// you can set buffered to false to avoid the overhead of buffering the
// response.
//
// This middleware will log a debug message containing the stack trace of the
// panic. It will not log an error, as the application can choose to do so or
// not by using ErrorLoggingMiddleware. If using ErrorLoggingMiddleware, it
// must be the outermost (leftmost) middleware.
func PanicMiddleware(logger logging.Logger, buffered bool) Middleware {
	return func(pattern string, next ht.Handler) ht.Handler {
		return ht.HandlerFunc(func(w ht.ResponseWriter, r *ht.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the stack trace as a debug message. Do not log an
					// error, as the application can choose to do so or not by
					// using ErrorLoggingMiddleware.
					logger.Debugw("HTTP panic stack trace", "stack", string(debug.Stack()))
					ht.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			var bw *BufferedResponseWriter
			if buffered {
				bw = NewBufferedResponseWriter(w)
				next.ServeHTTP(bw, r)
				bw.Flush()
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

// WithPanic returns a ServerOption that adds PanicMiddleware to the server.
//
// PanicMiddleware is an HTTP middleware that recovers from panics and returns a
// 500 Internal Server Error. If buffered is true, the middleware will buffer
// the response and flush it only after the handler has succeeded. This is
// useful if the handler writes to the response before panicking. If you know
// that the handler will *not* write to the response before panicking, then you
// can set buffered to false to avoid the overhead of buffering the response.
//
// This middleware will log a debug message containing the stack trace of the
// panic. It will not log an error, as the application can choose to do so or
// not by using ErrorLoggingMiddleware. If using ErrorLoggingMiddleware, it
// must be the outermost (leftmost) middleware.
func WithPanic(logger logging.Logger, buffered bool) ServerOption {
	return WithMiddleware(PanicMiddleware(logger, buffered))
}
