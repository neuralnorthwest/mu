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

	"github.com/neuralnorthwest/mu/config"
	"github.com/neuralnorthwest/mu/http"
	"github.com/neuralnorthwest/mu/service"
	"github.com/neuralnorthwest/mu/worker"
)

// Demonstrates basic usage of the mu framework.
//
//  - Define BasicApp, which embeds an instance of service.Service
//  - Define newBasicApp, which initializes the application:
//    - Create a new service.Service
//    - Register the configuration setup hook `setupConfig`
//    - Register the worker setup hook `setupWorkers`
//    - Register the cleanup hook `cleanup`
//  - Define setupConfig, which sets up the configuration:
//    - Create a new string configuration variable `MESSAGE`
//  - Define setupWorkers, which sets up the application workers:
//    - Create a new HTTP server
//    - Register a handler for the `/hello` endpoint
//    - Add the HTTP server to the worker group
//  - Define cleanup, which cleans up the application:
//    - Nothing to do here
//  - Define main, which runs the application:
//    - Create a new basic application using newBasicApp
//    - Call Main on the application
//    - If an error occurs, print it to stderr and exit with a non-zero status

// BasicApp is a basic application.
type BasicApp struct {
	*service.Service
}

// NewBasicApp returns a new basic application.
func newBasicApp() (*BasicApp, error) {
	svc, err := service.New("basic")
	if err != nil {
		return nil, err
	}
	app := &BasicApp{
		Service: svc,
	}
	app.SetupConfig(app.setupConfig)
	app.SetupWorkers(app.setupWorkers)
	app.Cleanup(app.cleanup)
	return app, nil
}

// setupConfig sets up the configuration.
func (a *BasicApp) setupConfig(c config.Config) error {
	return c.NewString("MESSAGE", "Hello, World!", "The message to print.")
}

// setupWorkers sets up the application.
func (a *BasicApp) setupWorkers(workerGroup worker.Group) error {
	httpServer, err := http.NewServer()
	if err != nil {
		return err
	}
	httpServer.HandleFunc("/hello", func(w ht.ResponseWriter, r *ht.Request) {
		_, _ = w.Write([]byte(a.Config().String("MESSAGE")))
	})
	return workerGroup.Add("http_server", httpServer)
}

// cleanup cleans up the application.
func (a *BasicApp) cleanup() {
	a.Logger().Info("Cleaning up...")
}

func main() {
	app, err := newBasicApp()
	if err != nil {
		panic(err)
	}
	if err := app.Main(); err != nil {
		os.Exit(1)
	}
}
