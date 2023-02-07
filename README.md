 # Mu

![GitHub](https://img.shields.io/github/license/neuralnorthwest/mu?style=plastic)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/neuralnorthwest/mu?style=plastic)
![GitHub Workflow Status (with branch)](https://img.shields.io/github/actions/workflow/status/neuralnorthwest/mu/cicd.yaml?branch=develop&style=plastic)
![GitHub search hit counter](https://img.shields.io/github/search/neuralnorthwest/mu/goto?style=plastic)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/neuralnorthwest/mu?style=plastic)
![Lines of code](https://img.shields.io/badge/lines%20of%20code-8k-blue?style=plastic)
![Status](https://img.shields.io/badge/status-in%20development-orange?style=plastic)

Mu is a microservice framework written in Go. Built on ideas and learnings from
real-world microservice development, Mu is designed to be simple, fast, and
easy to use.

Mu is highly modular, with functionality divided into small, independent
packages. This allows you to pick and choose the features you need, and
avoid the bloat of unnecessary dependencies.

## Hello world in Mu

The following program is a HTTP microservice that listens on port 8080 and
returns a response when you visit `/hello`.

```go
package main

import (
	ht "net/http"
	"os"

	"github.com/neuralnorthwest/mu/config"
	"github.com/neuralnorthwest/mu/http"
	"github.com/neuralnorthwest/mu/service"
	"github.com/neuralnorthwest/mu/worker"
)

func main() {
	s, _ := service.New("hello-service")
	s.ConfigSetup(func(c config.Config) error {
		return c.NewString("MESSAGE", "Hello world!", "The message to print")
	})
	s.Setup(func(g worker.Group) error {
		httpServer, err := http.NewServer()
		if err != nil {
			return err
		}
		httpServer.HandleFunc("/hello", func(w ht.ResponseWriter, r *ht.Request) {
			s.Logger().Infow("Received request", "method", r.Method, "url", r.URL)
			_, _ = w.Write([]byte(s.Config().String("MESSAGE")))
		})
		return g.Add("http_server", httpServer)
	})
	if err := s.Main(); err != nil {
		os.Exit(1)
	}
}
```

## Features

-   **Hooks** - Mu uses hooks to allow you to easily extend the functionality of
    the framework. Hooks are called at specific points in the lifecycle of a
    microservice. You can use hooks to add custom functionality to your
    microservice. The example above uses two hooks, `ConfigSetup` and `Setup`.
-   **Configuration** - Mu is designed to be configured using environment
    variables. This allows you to easily configure your microservice using
    Docker, Kubernetes, or any other container orchestration system.
-   **Logging** - Mu uses the [zap](https://github.com/uber-go/zap) logging
    library. Zap is fast, structured, and easy to use. Mu also provides a mock
    logger that can be used in tests.
-   **Workers** - Mu is based on the idea of workers. A worker is an independent
    process that performs a specific task. There are pre-built workers for
    HTTP servers, gRPC servers, and more. You can also create your own workers
    to perform any task you need.
-   **Graceful shutdown** - Mu provides a hook for registering cleanup functions
    that are called when the microservice is shutting down. This allows you to
    perform any cleanup tasks you need, such as closing database connections,
    before the microservice exits. Mu automatically handles SIGTERM and SIGINT
    correctly, and will wait for all workers to finish before exiting.

## Developer quick start

If you want work on Mu, you can use the following commands to get started:

```bash
git clone https://github.com/neuralnorthwest/mu.git
cd mu
make setup-dev
```

### Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for more information.

## Code of Conduct

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for more information.
