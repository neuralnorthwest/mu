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

import "crypto/tls"

// DefaultTLSConfig is the default TLS configuration if one is not specified
// when calling WithTLS. This is initialized with a config that supports only
// TLS 1.3. By changing the fields of DefaultTLSConfig prior to creating
// servers, you can change the default TLS configuration for all servers. You
// can also override it on a per-server basis by passing a non-nil *tls.Config
// to WithTLS.
var DefaultTLSConfig = &tls.Config{
	MinVersion: tls.VersionTLS13,
	MaxVersion: tls.VersionTLS13,
}

// tlsConfig is the TLS configuration for the server.
type tlsConfig struct {
	tlsConfig *tls.Config
	// certFile is the certificate file.
	certFile string
	// keyFile is the keyFile file.
	keyFile string
}

// WithTLS returns an option that sets the TLS configuration for the server.
//   - certFile is the certificate file in PEM format. It must contain the
//     complete certificate chain, including the CA and any intermediates.
//   - keyFile is the key file in PEM format. It must contain the private key
//     for the certificate.
//
// If tlsConf is nil, then the server will use DefaultTLSConfig.
//
// If certFile or keyFile are empty, you must set the Certificates field of
// tlsConf.
func WithTLS(tlsConf *tls.Config, certFile, keyFile string) ServerOption {
	return func(s *Server) error {
		if tlsConf == nil {
			tlsConf = DefaultTLSConfig
		}
		s.tlsConfig = &tlsConfig{
			tlsConfig: tlsConf,
			certFile:  certFile,
			keyFile:   keyFile,
		}
		return nil
	}
}
