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
	"context"
	"errors"
	ht "net/http"
	"time"

	"github.com/neuralnorthwest/mu/logging"
)

// Server is an HTTP server.
type Server struct {
	// server is the underlying HTTP server.
	server *ht.Server
	// mux is the HTTP request multiplexer.
	mux *ht.ServeMux
	// shutdownTimeout is the timeout for graceful shutdown.
	shutdownTimeout time.Duration
	// logger is the logger for the server.
	logger logging.Logger
	// middleware is the middleware for the server.
	middleware []Middleware
}

// ServerOption is an option for the HTTP server.
type ServerOption func(*Server) error

// WithAddress returns an option that sets the address to listen on.
func WithAddress(addr string) ServerOption {
	return func(s *Server) error {
		s.server.Addr = addr
		return nil
	}
}

// WithShutdownTimeout returns an option that sets the timeout for graceful shutdown.
func WithShutdownTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) error {
		s.shutdownTimeout = timeout
		return nil
	}
}

// NewServer creates a new HTTP server.
func NewServer(opts ...ServerOption) (*Server, error) {
	mux := ht.NewServeMux()
	s := &Server{
		server: &ht.Server{
			Handler: mux,
			Addr:    ":8080",
		},
		mux:             mux,
		shutdownTimeout: 5 * time.Second,
	}
	for _, o := range opts {
		if err := o(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// Handle registers the handler for the given pattern.
func (s *Server) Handle(pattern string, handler ht.Handler) {
	for i := len(s.middleware) - 1; i >= 0; i-- {
		handler = s.middleware[i](pattern, handler)
	}
	s.mux.Handle(pattern, handler)
}

// HandleFunc registers the handler function for the given pattern.
func (s *Server) HandleFunc(pattern string, handler func(ht.ResponseWriter, *ht.Request)) {
	s.Handle(pattern, ht.HandlerFunc(handler))
}

// Logger returns the logger for the server.
func (s *Server) Logger() logging.Logger {
	return s.logger
}

// Run starts the HTTP server. It implements worker.Worker.
func (s *Server) Run(ctx context.Context, logger logging.Logger) error {
	s.logger = logger
	shutdown := make(chan error, 1)
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		s.logger.Debugw("shutting down HTTP server")
		shutdown <- s.server.Shutdown(ctx)
	}()
	s.logger.Debugw("starting HTTP server", "addr", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		if !errors.Is(err, ht.ErrServerClosed) {
			s.logger.Errorw("HTTP server error", "err", err)
			return err
		}
		return nil
	}
	err := <-shutdown
	if err != nil {
		s.logger.Errorw("HTTP server shutdown error", "err", err)
		return err
	}
	s.logger.Debugw("HTTP server stopped")
	return nil
}
