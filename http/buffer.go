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
	ht "net/http"
)

// BufferedResponseWriter is a wrapper around http.ResponseWriter that buffers
// the response body.
type BufferedResponseWriter struct {
	writer     ht.ResponseWriter
	body       bytes.Buffer
	statusCode int
}

var _ ht.ResponseWriter = &BufferedResponseWriter{}

// NewBufferedResponseWriter returns a new BufferedResponseWriter.
func NewBufferedResponseWriter(w ht.ResponseWriter) *BufferedResponseWriter {
	return &BufferedResponseWriter{
		writer: w,
	}
}

// Header returns the header map that will be sent by WriteHeader.
func (w *BufferedResponseWriter) Header() ht.Header {
	return w.writer.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
func (w *BufferedResponseWriter) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

// WriteHeader sends an HTTP response header with the provided status code.
func (w *BufferedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// Flush writes the buffered response to the underlying http.ResponseWriter.
func (w *BufferedResponseWriter) Flush() (int, error) {
	if w.statusCode != 0 {
		w.writer.WriteHeader(w.statusCode)
	}
	n, err := w.body.WriteTo(w.writer)
	return int(n), err
}
