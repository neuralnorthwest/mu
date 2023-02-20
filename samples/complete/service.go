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

package complete

import (
	"github.com/neuralnorthwest/mu/bug"
	"github.com/neuralnorthwest/mu/http"
	"github.com/neuralnorthwest/mu/service"
	"github.com/neuralnorthwest/mu/trace"
)

// complete is the main struct of the complete service.
type complete struct {
	// Embed a Mu service.Service.
	*service.Service
	// metrics holds the service metrics.
	metrics met
	// diagnosticsServer is the diagnostics HTTP server.
	diagnosticsServer *http.Server
	// config is the service configuration wrapper.
	config *Config
	// tracer is the tracer for the service.
	tracer trace.Tracer
}

// New creates a new complete service.
func New(name string) (*complete, error) {
	// First, we initialize the service.Service. This is the core of the
	// service. It handles the service lifecycle, logging, and configuration.
	srv, err := service.New(name)
	if err != nil {
		return nil, err
	}
	// Next, we initialize the complete service. This is the service-specific
	// code. We embed the service.Service so that we can access the service
	// lifecycle, logging, and configuration.
	s := &complete{
		Service: srv,
	}
	// Configure a bug handler that will send bug reports to the service
	// logger. See bug_handler.go.
	bug.SetHandler(s.bugHandler)
	// Next, we initialize OTLP tracing. See trace.go.
	if err := s.setupTracing(); err != nil {
		return nil, err
	}
	// Next, we initialize the metrics for the service. See metrics.go.
	if err := s.setupMetrics(); err != nil {
		return nil, err
	}
	// Add a configuration setup hook. See setup.go.
	s.SetupConfig(s.setupConfig)
	// Add a worker setup hook. We use this to start a secondary HTTP server
	// that serves metrics. See worker.go.
	s.SetupWorkers(s.setupWorkers)
	// Add a HTTP setup hook. We enable error and panic logging, and metrics.
	// See http.go.
	metricsOpts := http.MetricsOptions{
		Server: func() *http.Server { return s.diagnosticsServer },
	}
	s.SetupHTTP(s.setupHTTP,
		http.WithPanicAndErrorLogging(s.Logger(), true),
		http.WithMetrics(s.metrics.metrics, s.Logger(), metricsOpts),
	)
	// Add a pre-run hook. We use this to perform any other initialization
	// that we need to do before the service starts. See prerun.go.
	s.PreRun(s.preRun)
	return s, nil
}
