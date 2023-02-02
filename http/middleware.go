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
)

// Middleware is an HTTP middleware.
type Middleware func(ht.Handler) ht.Handler

// WithMiddleware returns an option that adds the given middleware to the server.
// The first middleware in the list is the outermost middleware.
func WithMiddleware(m Middleware, more ...Middleware) ServerOption {
	return func(s *Server) error {
		s.middleware = append(s.middleware, m)
		s.middleware = append(s.middleware, more...)
		return nil
	}
}
