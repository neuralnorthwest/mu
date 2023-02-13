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

package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/neuralnorthwest/mu/logging"
	gr "google.golang.org/grpc"
)

// Server is a gRPC server.
type Server struct {
	// server is the underlying gRPC server.
	server *gr.Server
	// addr is the address to listen on.
	addr string
	// logger is the logger for the server.
	logger logging.Logger
	// grOpts are the gRPC server options.
	grOpts []gr.ServerOption
	// listener is the listener for the gRPC server.
	listener net.Listener
}

// ServerOption is a gRPC server option.
type ServerOption func(*Server) error

// WithAddress returns an option that sets the address to listen on.
func WithAddress(addr string) ServerOption {
	return func(s *Server) error {
		s.addr = addr
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

// WithOption returns an option that sets a gRPC server option.
func WithOption(opt gr.ServerOption) ServerOption {
	return func(s *Server) error {
		s.grOpts = append(s.grOpts, opt)
		return nil
	}
}

// NewServer returns a new Server.
func NewServer(opts ...ServerOption) (*Server, error) {
	s := &Server{
		addr: ":8081",
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	if s.listener == nil {
		var err error
		s.listener, err = net.Listen("tcp", s.addr)
		if err != nil {
			return nil, fmt.Errorf("failed to listen: %w", err)
		}
	} else {
		s.addr = s.listener.Addr().String()
	}
	s.server = gr.NewServer(s.grOpts...)
	return s, nil
}

// Logger returns the logger for the server.
func (s *Server) Logger() logging.Logger {
	return s.logger
}

// Run runs the gRPC server. It implements worker.Worker.
func (s *Server) Run(ctx context.Context, logger logging.Logger) error {
	s.logger = logger
	shutdown := make(chan error, 1)
	go func() {
		<-ctx.Done()
		logger.Debugw("shutting down gRPC server")
		s.server.GracefulStop()
		shutdown <- nil
	}()
	logger.Debugw("starting gRPC server", "addr", s.listener.Addr())
	err := s.server.Serve(s.listener)
	if err != nil {
		logger.Errorw("gRPC server error", "err", err)
		shutdown <- err
	}
	err = <-shutdown
	logger.Debugw("gRPC server stopped")
	return err
}

// Address returns the address the server is listening on.
func Address(s *Server) string {
	return s.listener.Addr().String()
}
