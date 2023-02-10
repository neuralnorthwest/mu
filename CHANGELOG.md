## v0.1.14

TODO

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
