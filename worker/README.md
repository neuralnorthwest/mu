# mu/worker

The `worker` package provides two types: `Worker` and `Group`.

## Worker

`Worker` is an interface that standardizes the running of worker processes.
It has a single method, `Run`, which runs the worker and returns when the worker
is done, the context is canceled, or an error occurs.

```go
type Worker interface {
    // Run runs the worker. The worker will stop when the context is canceled
    // or an error occurs.
    Run(ctx context.Context, logger logging.Logger) error
}
```

It is the responsibility of the worker to monitor `ctx` and stop when it
is canceled.

`logger` should be used for all logging. The logger will automatically include
the name of the worker in all log messages.

## Group


`Group` is a collection of workers. It provides an interface for registering
workers and starting them all at once.

```go
type Group interface {
    // Add adds a worker to the worker group. The worker will be started when the
    // worker group is started. If the group has already been started, the worker
    // will be started immediately.
    Add(name string, worker Worker) error
    // Run runs the worker group. This will start all the workers in the
    // worker group. This will block until the context is canceled or a worker
    // returns an error.
    Run(ctx context.Context, logger logging.Logger) error
    // Start starts the worker group. This will start all the workers in the
    // worker group. This will not block. To wait for the workers to stop, call
    // Wait and cancel the context.
    Start(ctx context.Context, logger logging.Logger) error
    // Wait waits for the worker group to stop.
    Wait() error
}
```

## Usage

To create a new `Group`, use `NewGroup`.

```go
g := worker.NewGroup()
```

To add a worker to the group, use `Add`.

```go
g.Add("worker1", func(ctx context.Context, logger logging.Logger) error {
    // Do work
    return nil
})
```

To start the group, use `Start`.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
g.Start(ctx, logger)
```

To wait for the group to stop, use `Wait`.

```go
g.Wait()
```

For convenience, `Run` is provided. This will start the group and wait for it
to stop.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
g.Run(ctx, logger)
```
