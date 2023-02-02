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
//    - Register the configuration setup hook `configSetup`
//    - Register the setup hook `setup`
//    - Register the cleanup hook `cleanup`
//  - Define configSetup, which sets up the configuration:
//    - Create a new string configuration variable `MESSAGE`
//  - Define setup, which sets up the application:
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
	app.RegisterConfigSetup(app.configSetup)
	app.RegisterSetup(app.setup)
	app.RegisterCleanup(app.cleanup)
	return app, nil
}

// configSetup sets up the configuration.
func (a *BasicApp) configSetup(c config.Config) error {
	return c.NewString("MESSAGE", "Hello, World!", "The message to print.")
}

// setup sets up the application.
func (a *BasicApp) setup(workerGroup worker.Group) error {
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
func (a *BasicApp) cleanup() error {
	a.Logger().Info("Cleaning up...")
	return nil
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
