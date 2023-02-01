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
		w.Write([]byte(a.Config().String("MESSAGE")))
	})
	workerGroup.Add("http_server", httpServer)
	return nil
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
