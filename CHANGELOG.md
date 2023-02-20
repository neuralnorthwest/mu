## v0.1.17

TODO

### Added

### Changed

### Deprecated

### Removed

### Fixed

### Security

## v0.1.16

### Added

* `complete` sample added.
* `grpc` package added with `grpc.Server`.
* `retry` package added.
* Protobuf support:
  * `setup-dev` target installs `buf`.
  * `make generate-proto` generates all proto files.
  * `make lint-proto` lints all proto files.
* `http.Server` adds new method `Address` to get the address the server is
  listening on.
* `status` updates:
  * New errors `ErrServerError` and `ErrClientError`.
  * `HTTPError` function to map a response status code and error to a Go
    error.
* Other enhancements:
  * `make roll-version` updates the version number in `version.go` and adds
    a new entry to the changelog.
* Main `README.md` is enhanced with links to package documentation.

### Changed

* `metrics.MetricsOptions` field `Server` is now a function instead of a
  `*http.Server`.
* `SetupHTTP` hook now accepts `opts ...http.ServerOption` to configure the
  server.

### Fixed

* `config.NewBool`, `config.NewInt`, and `config.NewString` now properly check
  for conflicting variables of other types, not just their own type.

### Security

## v0.1.15

* Bug fix:
  * `http.WithPanicAndErrorLogging` was installing the middleware in the wrong
    order, causing error logging of panics to not be effective.

## v0.1.14

* Semantics of `Cleanup` have been changed. It can now be called multiple times
  to register multiple cleanup functions. The cleanup functions are called in
  the reverse order in which they were registered. Cleanup functions can no
  longer return errors.

## v0.1.13

* Added `PanicMiddleware` to `http` package.

## v0.1.12

* Dependency upgrades

## v0.1.11

* The `ConfigSetup` hook has been renamed to `SetupConfig`. **This is a
  breaking change.** To upgrade, call `SetupConfig` instead of `ConfigSetup`
  in your application.
* New `Level` and `SetLevel` methods for `logging.Logger`.
* New `WithLevel` option for `logging.New`.

## v0.1.10

* New `PreRun` hook can be used to register a function to run immediately
  before worker start.
* The `Setup` hook has been renamed to `SetupWorkers`. **This is a breaking
  change.** To upgrade, call `SetupWorkers` instead of `Setup` in your
  application.

## v0.1.9

* Mu now has a CODE_OF_CONDUCT.md and CONTRIBUTING.md.

## v0.1.8

* OpenTelemetry tracing via `WithOpenTelemetryTracing` http.ServerOption.
* TLS support via `WithTLS` http.ServerOption.

## v0.1.7

* Metrics improvements
* Lots more tests

* Bug fixes:
  * `http.ErrorLoggingMiddleware` now correctly logs the error.
  * Fixes to `Run` method of `http.Server`.

## v0.1.6

* Prometheus metrics

## v0.1.5

* Add SetupHTTP hook.

## v0.1.4

* Improve release process

## v0.1.3

* Starting documentation

## v0.1.2

* README.md updates

## v0.1.1

* Rename pr.yaml workflow to cicd.yaml
* Fix a bug in version tagging

## v0.1.0

This is the first alpha release of `mu`.
