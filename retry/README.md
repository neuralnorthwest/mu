# mu/retry

The `retry` package provides a retry mechanism. This is useful for retrying
operations that may fail due to transient errors.

## Do function

The `Do` function will retry an operation until it succeeds or the strategy
determines that the operation should no longer be retried.

```go
// Do retries the provided function according to the provided strategy, until
// the context is canceled or the function returns nil.
func Do(ctx context.Context, strategy Strategy, fn func() error, opts ...DoOption) error
```

`ctx` is the context to use for the retry operation. If the context is
canceled, the retry operation will stop.

`strategy` is the strategy to use for retrying the operation. The strategy
receives the most recent error returned by the function and returns a
duration to wait before retrying the operation. If the strategy returns a
negative duration, the retry operation will stop, and the most recent error
will be returned.

`fn` is the function to retry. If the function returns nil, the retry
operation will stop. If the function returns an error, the strategy will be
used to determine how long to wait before retrying the operation.

`opts` are options to use for the retry operation.

### Do options

The following options are available:

- `OnRetryAttempt` specifies a function to call after each retry attempt. The
  function receives the most recent error returned by the function and the
  number of attempts that have been made. If the function returns an error,
  the retry operation will stop and the error will be returned. Otherwise, the
  retry operation will continue.

  ```go
  func OnRetryAttempt(fn func(attempt int, err error) error) DoOption
  ```

- `OnRetry` specifies a function to call after each retry. The function
  receives the most recent error returned by the function. If the function
  returns an error, the retry operation will stop and the error will be
    returned. Otherwise, the retry operation will continue.

    ```go
    func OnRetry(fn func(err error) error) DoOption
    ```

## Strategy

The `Strategy` interface is used to determine how long to wait before retrying
an operation.

```go
// Strategy is a retry strategy.
type Strategy interface {
	// Next returns the duration to wait before the next retry. It is passed
	// the error returned by the last attempt. It returns a negative duration
	// to stop retrying.
	Next(err error) time.Duration
}
```

You can implement any retry strategy you want by implementing the `Strategy`
interface. The `retry` package provides a few strategies that you can use.

### Exponential

The `Exponential` strategy will wait an exponentially increasing amount of time
between retries. The initial time to wait is the "base interval," which is
multiplied by a fixed factor after every retry. A maximum interval can be
specified, which will cap the amount of time to wait between retries. In
addition, a maximum number of attempts can be specified, which will cap the
total number of retries.

```go
// Exponential returns a retry strategy that uses exponential backoff.
func Exponential(opts ...StrategyOption) Strategy
```

The `Exponential` strategy supports the following options:

- `WithBaseInterval` (default 100 milliseconds)
- `WithFactor` (default 2)
- `WithMaxAttempts` (default unlimited)
- `WithMaxInterval` (default 10 seconds)

### Linear

The `Linear` strategy will wait a linearly increasing amount of time between
retries. The initial time to wait is the "base interval," which is increased
by a fixed amount after every retry. A maximum interval can be specified,
which will cap the amount of time to wait between retries. In addition, a
maximum number of attempts can be specified, which will cap the total number
of retries.

```go
// Linear returns a retry strategy that uses linear backoff.
func Linear(opts ...StrategyOption) Strategy
```

The `Linear` strategy supports the following options:

- `WithBaseInterval` (default 100 milliseconds)
- `WithIncrement` (default 100 milliseconds)
- `WithMaxAttempts` (default unlimited)
- `WithMaxInterval` (default 10 seconds)

### Fixed

The `Fixed` strategy will wait a fixed amount of time between retries. A maximum
number of attempts can be specified, which will cap the total number of retries.
`Fixed` is a special case of `Linear` with an increment of 0.

```go
// Fixed returns a linear retry strategy with an increment of 0.
func Fixed(opts ...StrategyOption) Strategy
```

The `Fixed` strategy supports the following options:

- `WithBaseInterval` (default 100 milliseconds)
- `WithMaxAttempts` (default unlimited)

## Strategy Options

A `StrategyOption` is used to configure a strategy. Not all strategies support
all options. If an option is not supported, it will be ignored and `bug.Bug`
will be called.

```go
// StrategyOption is an option for a Strategy.
type StrategyOption func(Strategy)
```

### WithBaseInterval

The `WithBaseInterval` option specifies the base interval for a strategy. The
base interval is the amount of time to wait before the first retry. This can
be any non-negative duration. For `Exponential` strategies, it must be strictly
positive.

```go
func WithBaseInterval(d time.Duration) StrategyOption
```

### WithFactor

The `WithFactor` option specifies the factor for a strategy. The factor is
used to increase the amount of time to wait between retries. This can be any
number greater than or equal to 1.

```go
func WithFactor(f float64) StrategyOption
```

### WithIncrement

The `WithIncrement` option specifies the increment for a strategy. The
increment is used to increase the amount of time to wait between retries. This
can be any non-negative duration.

```go
func WithIncrement(d time.Duration) StrategyOption
```

### WithMaxAttempts

The `WithMaxAttempts` option specifies the maximum number of attempts for a
strategy. The maximum number of attempts is the total number of retries. This
can be any integer. Negative values mean there is no maximum number of
attempts.

```go
func WithMaxAttempts(n int) StrategyOption
```

### WithMaxInterval

The `WithMaxInterval` option specifies the maximum interval for a strategy. The
maximum interval is the maximum amount of time to wait between retries. This
can be any duration. Negative values mean there is no maximum interval. For
convenience, `WithNoMaxInterval` is provided as a synonym for `WithMaxInterval`
with a negative value.

```go
func WithMaxInterval(d time.Duration) StrategyOption
```

### WithNoMaxInterval

The `WithNoMaxInterval` option specifies that there is no maximum interval for a
strategy. This is a synonym for `WithMaxInterval` with a negative value.

```go
func WithNoMaxInterval() StrategyOption
```

## Specifying a timeout

Timeout is not implemented with any option. Instead, set a timeout or deadline
on the context passed to `Do`.

## Examples

The implementation of `HTTPGet` serves as a simple example.

```go
// HTTPGet is a convenience function that performs an HTTP GET request with
// the given URL and returns the response body. It retries the request according
// to the provided strategy, until the context is canceled or the request
// succeeds. "Success" is defined as any response with a non-5xx status code.
// The caller is responsible for closing the response body.
func HTTPGet(ctx context.Context, url string, strategy Strategy, opts ...DoOption) (*ht.Response, error) {
	var resp *ht.Response
	err := Do(ctx, strategy, func() error {
		resp, err := ht.Get(url)
		return status.HTTPError(resp, err)
	}, opts...)
	return resp, err
}
```
