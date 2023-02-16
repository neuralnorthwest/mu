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
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	grpc_test "github.com/neuralnorthwest/mu/grpc/test/proto/gen/go/grpc_test/v1"
	"github.com/neuralnorthwest/mu/logging"
	mock_logging "github.com/neuralnorthwest/mu/logging/mock"
	"github.com/stretchr/testify/assert"
	gr "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// grpcTestService implements the TestService gRPC service.
type grpcTestService struct {
	server *Server
	grpc_test.UnimplementedGRPCTestServiceServer
	started atomic.Bool
}

// GRPCTest implements the TestService gRPC service.
func (s *grpcTestService) GRPCTest(ctx context.Context, req *grpc_test.GRPCTestRequest) (*grpc_test.GRPCTestResponse, error) {
	if s.started.Load() {
		s.server.Logger().Debugw("received request", "name", req.GetName())
	}
	return &grpc_test.GRPCTestResponse{
		Message: "Hello, " + req.GetName(),
	}, nil
}

// createTestClient creates a GRPCTest client
func createTestClient(t *testing.T, addr string) (grpc_test.GRPCTestServiceClient, error) {
	t.Helper()
	conn, err := gr.Dial(addr, gr.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "unable to dial server")
	return grpc_test.NewGRPCTestServiceClient(conn), nil
}

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
	grpcTestService := &grpcTestService{
		server: srv,
	}
	grpc_test.RegisterGRPCTestServiceServer(srv.server, grpcTestService)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		srvErr <- srv.Run(ctx, logger)
	}()
	// Make requests to GRPCTest until it returns a response. Back off exponentially
	// for up to 1 second.
	addr := ":8081"
	if listener != nil {
		addr = listener.Addr().String()
	}
	grpcTestServiceClient, err := createTestClient(t, addr)
	assert.NoError(t, err, "unable to create client")
	resp, err := grpcTestServiceClient.GRPCTest(ctx, &grpc_test.GRPCTestRequest{Name: "world"})
	retryInterval := 10 * time.Millisecond
	retryTime := 0 * time.Millisecond
	retryCount := 0
	for err != nil || resp.Message != "Hello, world" {
		retryCount++
		if retryTime > time.Second {
			t.Fatalf("server did not become ready: %s", err)
		}
		time.Sleep(retryInterval)
		retryInterval *= 2
		retryTime += retryInterval
		resp, err = grpcTestServiceClient.GRPCTest(ctx, &grpc_test.GRPCTestRequest{Name: "world"})
	}
	t.Logf("server became ready after %d retries and %s", retryCount, retryTime)
	grpcTestService.started.Store(true)
	return srv, cancel
}

// Test_Server tests that the server can be created, started, and a request
// can be made.
func Test_Server(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	logger := mock_logging.NewMockLogger(mockCtrl)
	logger.EXPECT().Debugw("starting gRPC server", "addr", gomock.Any())
	logger.EXPECT().Debugw("received request", "name", "world")
	logger.EXPECT().Debugw("shutting down gRPC server")
	logger.EXPECT().Debugw("gRPC server stopped")
	srvErr := make(chan error, 1)
	listener, err := net.Listen("tcp", ":0")
	assert.NoError(t, err, "unable to create listener")
	srv, cancel := createTestServer(t, logger, listener, srvErr)
	defer cancel()
	// Check the address of the server.
	addr := Address(srv)
	assert.Equal(t, listener.Addr().String(), addr, "expected address %s, got %s", listener.Addr().String(), addr)
	t.Logf("making request to %s", addr)
	// Make a request to the server.
	grpcTestServiceClient, err := createTestClient(t, addr)
	assert.NoError(t, err, "unable to create client")
	resp, err := grpcTestServiceClient.GRPCTest(context.Background(), &grpc_test.GRPCTestRequest{Name: "world"})
	assert.NoError(t, err, "unable to make request")
	assert.Equal(t, "Hello, world", resp.Message, "expected message %s, got %s", "Hello, world", resp.Message)
	// Stop the server.
	cancel()
	if err := <-srvErr; err != nil {
		t.Fatalf("server returned error: %s", err)
	}
}
