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
	"crypto/tls"
	"reflect"
	"testing"
)

// Test_DefaultTLSConfig tests that the default TLS config is correct.
func Test_DefaultTLSConfig(t *testing.T) {
	t.Parallel()
	if !reflect.DeepEqual(DefaultTLSConfig, &tls.Config{
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
	}) {
		t.Errorf("unexpected TLS config: %v", DefaultTLSConfig)
	}
}

// Test_WithTLS tests that the WithTLS option sets the correct values in the
// server.
func Test_WithTLS(t *testing.T) {
	t.Parallel()
	myConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS12,
	}
	srv, err := NewServer(WithTLS(myConfig, "certFile", "keyFile"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(srv.tlsConfig.tlsConfig, myConfig) {
		t.Errorf("unexpected TLS config: %v", srv.tlsConfig)
	}
	if srv.tlsConfig.certFile != "certFile" {
		t.Errorf("unexpected cert file: %v", srv.tlsConfig.certFile)
	}
	if srv.tlsConfig.keyFile != "keyFile" {
		t.Errorf("unexpected key file: %v", srv.tlsConfig.keyFile)
	}
}

// With_OverrideDefaultTLSConfig tests that WithTLS propagates changes to
// DefaultTLSConfig.
func Test_With_OverrideDefaultTLSConfig(t *testing.T) {
	// Not a parallel test, because we are changing the global DefaultTLSConfig.
	myConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS12,
	}
	previousDefault := DefaultTLSConfig
	defer func() {
		DefaultTLSConfig = previousDefault
	}()
	DefaultTLSConfig = myConfig
	srv, err := NewServer(WithTLS(nil, "certFile", "keyFile"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(srv.tlsConfig.tlsConfig, myConfig) {
		t.Errorf("unexpected TLS config: %v", srv.tlsConfig)
	}
	if srv.tlsConfig.certFile != "certFile" {
		t.Errorf("unexpected cert file: %v", srv.tlsConfig.certFile)
	}
	if srv.tlsConfig.keyFile != "keyFile" {
		t.Errorf("unexpected key file: %v", srv.tlsConfig.keyFile)
	}
}
