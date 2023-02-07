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
	"net"
	ht "net/http"
	"net/http/httptest"
	"testing"
)

// Test_Middleware_Case is a test case for Test_Middleware.
type Test_Middleware_Case struct {
	name            string
	middlewareNames []string
}

// middleTrack tracks middleware invocation.
type middleTrack struct {
	sequence []string
}

// middlewareTest is a middleware that tracks invocation.
func middlewareTest(name string, track *middleTrack) Middleware {
	return func(pattern string, next ht.Handler) ht.Handler {
		return ht.HandlerFunc(func(w ht.ResponseWriter, r *ht.Request) {
			track.sequence = append(track.sequence, name)
			next.ServeHTTP(w, r)
		})
	}
}

// Test_Middleware tests that middleware is invoked in the correct order.
func Test_Middleware(t *testing.T) {
	t.Parallel()
	testCases := []Test_Middleware_Case{
		{
			name: "single middleware",
			middlewareNames: []string{
				"first",
			},
		},
		{
			name: "multiple middleware",
			middlewareNames: []string{
				"first",
				"second",
				"third",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			track := &middleTrack{}
			listener, err := net.Listen("tcp", ":0")
			if err != nil {
				t.Fatalf("net.Listen() error = %v", err)
			}
			defer listener.Close()
			opts := []ServerOption{}
			for _, name := range tc.middlewareNames {
				opts = append(opts, WithMiddleware(middlewareTest(name, track)))
			}
			opts = append(opts, WithListener(listener))
			srv, err := NewServer(opts...)
			if err != nil {
				t.Fatalf("NewServer() error = %v", err)
			}
			// check the middleware on the server
			if len(srv.middleware) != len(tc.middlewareNames) {
				t.Errorf("len(srv.middleware) = %v, want %v", len(srv.middleware), len(tc.middlewareNames))
			}
			// wrap a handler
			handler := srv.wrapHandler("/", ht.HandlerFunc(func(w ht.ResponseWriter, r *ht.Request) { w.WriteHeader(ht.StatusOK) }))
			// invoke the handler
			req := httptest.NewRequest(ht.MethodGet, "http://localhost/", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			// check the order in the tracker
			if len(track.sequence) != len(tc.middlewareNames) {
				t.Errorf("len(track.sequence) = %v, want %v", len(track.sequence), len(tc.middlewareNames))
			}
			for i, name := range tc.middlewareNames {
				if track.sequence[i] != name {
					t.Errorf("track.sequence[%v] = %v, want %v", i, track.sequence[i], name)
				}
			}
		})
	}
}
