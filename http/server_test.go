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
	"fmt"
	"net"
	ht "net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/neuralnorthwest/mu/logging"
	mock_logging "github.com/neuralnorthwest/mu/logging/mock"
)

// createTestServer creates a server with the given options and waits for it
// to be ready.
func createTestServer(t *testing.T, logger logging.Logger, listener net.Listener, srvErr chan<- error, opts ...ServerOption) (*Server, context.CancelFunc) {
	t.Helper()
	if listener != nil {
		opts = append(opts, WithListener(listener))
	}
	srv, err := NewServer(opts...)
	if err != nil {
		t.Fatalf("unable to create server: %s", err)
	}
	srv.HandleFunc("/readyz", func(w ht.ResponseWriter, r *ht.Request) {
		w.WriteHeader(ht.StatusOK)
	})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := srv.Run(ctx, logger)
		srvErr <- err
	}()
	// Make requests to /readyz until it returns 200. Back off exponentially
	// for up to 1 second.
	addr := srv.server.Addr
	if listener != nil {
		addr = listener.Addr().String()
	}
	resp, err := ht.Get("http://" + addr + "/readyz")
	retryInterval := 10 * time.Millisecond
	retryTime := 0 * time.Millisecond
	retryCount := 0
	for err != nil || resp.StatusCode != ht.StatusOK {
		retryCount++
		if retryTime > time.Second {
			t.Fatalf("server did not become ready: %s", err)
		}
		time.Sleep(retryInterval)
		retryInterval *= 2
		retryTime += retryInterval
		resp, err = ht.Get("http://" + addr + "/readyz")
	}
	t.Logf("server became ready after %d retries and %s", retryCount, retryTime)
	return srv, cancel
}

// Test_Server tests that the server can be created, started, and a request
// can be made.
func Test_Server(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	logger := mock_logging.NewMockLogger(mockCtrl)
	logger.EXPECT().Debugw("starting HTTP server", "addr", gomock.Any())
	logger.EXPECT().Debugw("handling request", "path", "/test")
	logger.EXPECT().Debugw("HTTP server closed")
	logger.EXPECT().Debugw("shutting down HTTP server")
	logger.EXPECT().Debugw("HTTP server stopped")
	srvErr := make(chan error, 1)
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("unable to create listener: %s", err)
	}
	srv, cancel := createTestServer(t, logger, listener, srvErr)
	defer cancel()
	// Check the address of the server.
	if Address(srv) != listener.Addr().String() {
		t.Fatalf("expected address %s, got %s", listener.Addr().String(), Address(srv))
	}
	srv.HandleFunc("/test", func(w ht.ResponseWriter, r *ht.Request) {
		srv.Logger().Debugw("handling request", "path", r.URL.Path)
		w.WriteHeader(ht.StatusOK)
	})
	addr := listener.Addr().String()
	t.Logf("making request to http://%s/test", addr)
	resp, err := ht.Get("http://" + addr + "/test")
	if err != nil {
		t.Fatalf("unable to make request: %s", err)
	}
	if resp.StatusCode != ht.StatusOK {
		t.Fatalf("expected status code %d, got %d", ht.StatusOK, resp.StatusCode)
	}
	cancel()
	if err := <-srvErr; err != nil {
		t.Fatalf("server returned error: %s", err)
	}
}

// Test_Server_Error tests that the server returns an error if it fails at
// any point. We cause the server to fail by closing the listener.
func Test_Server_Error(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	logger := mock_logging.NewMockLogger(mockCtrl)
	logger.EXPECT().Debugw("starting HTTP server", "addr", gomock.Any())
	logger.EXPECT().Errorw("HTTP server error", "err", gomock.Any())
	logger.EXPECT().Debugw("HTTP server stopped")
	srvErr := make(chan error)
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("unable to create listener: %s", err)
	}
	_, _ = createTestServer(t, logger, listener, srvErr)
	// Close the listener to cause the server to fail.
	listener.Close()
	//defer cancel()
	addr := listener.Addr().String()
	t.Logf("making request to http://%s/test", addr)
	resp, err := ht.Get("http://" + addr + "/test")
	if err != nil {
		t.Fatalf("unable to make request: %s", err)
	}
	if resp.StatusCode != ht.StatusNotFound {
		t.Fatalf("expected status code %d, got %d", ht.StatusOK, resp.StatusCode)
	}
	if err := <-srvErr; err == nil {
		t.Fatalf("expected server to return error")
	}
}

// Test_Server_WithAddress tests that the server can be created, started,
// and a request can be made.
func Test_Server_WithAddress(t *testing.T) {
	t.Parallel()
	port := 31381
	maxPort := 31390
	for p := port; p <= maxPort; p++ {
		t.Run(fmt.Sprintf("port %d", p), func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			logger := mock_logging.NewMockLogger(mockCtrl)
			logger.EXPECT().Debugw("starting HTTP server", "addr", gomock.Any())
			logger.EXPECT().Debugw("HTTP server closed")
			logger.EXPECT().Debugw("shutting down HTTP server")
			logger.EXPECT().Debugw("HTTP server stopped")
			addr := fmt.Sprintf(":%d", p)
			srvErr := make(chan error, 1)
			srv, cancel := createTestServer(t, logger, nil, srvErr, WithAddress(addr))
			defer cancel()
			srv.HandleFunc("/test", func(w ht.ResponseWriter, r *ht.Request) {
				w.WriteHeader(ht.StatusOK)
			})
			t.Logf("making request to http://%s/test", addr)
			resp, err := ht.Get("http://" + addr + "/test")
			if err != nil {
				t.Fatalf("unable to make request: %s", err)
			}
			if resp.StatusCode != ht.StatusOK {
				t.Fatalf("expected status code %d, got %d", ht.StatusOK, resp.StatusCode)
			}
			cancel()
			if err := <-srvErr; err != nil {
				t.Fatalf("server returned error: %s", err)
			}
			maxPort = p
		})
	}
}

// Test_Server_WithShutdownTimeout tests that the server shutdown times out
// if the shutdown timeout is exceeded.
func Test_Server_WithShutdownTimeout(t *testing.T) {
	t.Parallel()
	timescale := time.Duration(1) // increase if test deadlocks
	mockCtrl := gomock.NewController(t)
	logger := mock_logging.NewMockLogger(mockCtrl)
	logger.EXPECT().Debugw("starting HTTP server", "addr", gomock.Any())
	logger.EXPECT().Debugw("HTTP server closed")
	logger.EXPECT().Debugw("shutting down HTTP server")
	logger.EXPECT().Debugw("HTTP server stopped")
	srvErr := make(chan error, 1)
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("unable to create listener: %s", err)
	}
	srv, cancel := createTestServer(t, logger, listener, srvErr, WithShutdownTimeout(100*timescale*time.Millisecond))
	srv.HandleFunc("/test", func(w ht.ResponseWriter, r *ht.Request) {
		// Cancel the context, so that the server will start shutting down.
		cancel()
		// Read from the srvErr, to ensure that the handler does not return
		// until after the server has shut down.
		if err := <-srvErr; err != nil && err.Error() != "context deadline exceeded" {
			t.Fatalf("server returned error: %s", err)
		}
		w.WriteHeader(ht.StatusOK)
	})
	addr := listener.Addr().String()
	t.Logf("making request to http://%s/test", addr)
	resp, err := ht.Get("http://" + addr + "/test")
	if err != nil {
		t.Fatalf("unable to make request: %s", err)
	}
	// The server can't actually force handlers to stop, so we'll get a 200
	// response.
	if resp.StatusCode != ht.StatusOK {
		t.Fatalf("expected status code %d, got %d", ht.StatusOK, resp.StatusCode)
	}
}

// Test_Server_OptionError tests that NewServer returns an error if an option
// returns an error.
func Test_Server_OptionError(t *testing.T) {
	t.Parallel()
	_, err := NewServer(func(*Server) error {
		return errors.New("error")
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}
