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
