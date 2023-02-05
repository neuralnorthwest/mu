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

package main

import (
	ht "net/http"
	"os"

	"github.com/neuralnorthwest/mu/http"
	"github.com/neuralnorthwest/mu/metrics"
	"github.com/neuralnorthwest/mu/service"
	"github.com/neuralnorthwest/mu/worker"
)

// Demonstrates metrics collection and reporting.
//
//  - Define MetricsApp, which embeds an instance of service.Service
//  - Define newMetricsApp, which initializes the application:
//    - Create a new service.Service
//    - Register the setup hook `setup`
//  - Define setup, which sets up the application:
//    - Create a new HTTP server
//    - Register a handler for the `/hello` endpoint
//    - Register a handler for the `/metrics` endpoint
//  - Define main, which runs the application:
//    - Create a new basic application using newMetricsApp
//    - Call Main on the application
//    - If an error occurs, print it to stderr and exit with a non-zero status

// MetricsApp is a basic application.
type MetricsApp struct {
	*service.Service
	// helloCounter is a counter for the number of times the `/hello` endpoint
	// has been called. Note that this is for example purposes only. For
	// general HTTP metrics in production, you should use the WithMetrics
	// option when creating the HTTP server.
	helloCounter metrics.Counter
}

// NewBasicApp returns a new metrics application.
func newMetricsApp() (*MetricsApp, error) {
	svc, err := service.New("metrics")
	if err != nil {
		return nil, err
	}
	app := &MetricsApp{
		Service: svc,
	}
	app.Setup(app.setup)
	return app, nil
}

// setup sets up the application.
func (a *MetricsApp) setup(workerGroup worker.Group) error {
	met := metrics.New()
	a.helloCounter = met.NewCounter("hello_counter", "Number of times the /hello endpoint has been called")
	httpServer, err := http.NewServer()
	if err != nil {
		return err
	}
	httpServer.HandleFunc("/hello", func(w ht.ResponseWriter, r *ht.Request) {
		a.helloCounter.Inc()
		_, _ = w.Write([]byte(a.Config().String("MESSAGE")))
	})
	// This example serves the metrics on the same HTTP server as the
	// application. In production, you should serve the metrics on a separate
	// HTTP server accessible only to the monitoring system.
	httpServer.Handle("/metrics", met.Handler(a.Logger()))
	return workerGroup.Add("http_server", httpServer)
}

func main() {
	app, err := newMetricsApp()
	if err != nil {
		panic(err)
	}
	if err := app.Main(); err != nil {
		os.Exit(1)
	}
}
