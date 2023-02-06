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
	"net"
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
	// listener is the listener for the server.
	listener net.Listener
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

// WithListener returns an option that sets the listener for the server.
// This is useful if you want to bind to ":0" and then get the actual port
// that was bound to from the listener.
func WithListener(listener net.Listener) ServerOption {
	return func(s *Server) error {
		s.listener = listener
		return nil
	}
}

// WithShutdownTimeout returns an option that sets the timeout for graceful
// shutdown. The default is 5 seconds.
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
	s.mux.Handle(pattern, s.wrapHandler(pattern, handler))
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
	err := s.serve(s.listener)
	if err != nil {
		if !errors.Is(err, ht.ErrServerClosed) && err.Error() != "context deadline exceeded" {
			s.logger.Errorw("HTTP server error", "err", err)
			shutdown <- err
		} else {
			s.logger.Debugw("HTTP server closed")
			err = nil
		}
	}
	err = <-shutdown
	s.logger.Debugw("HTTP server stopped")
	return err
}

// serve calls Serve or ListenAndServe on the underlying HTTP server.
func (s *Server) serve(listener net.Listener) error {
	if listener != nil {
		return s.server.Serve(listener)
	}
	return s.server.ListenAndServe()
}
