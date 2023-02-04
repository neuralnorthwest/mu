# mu/config

`config` is a package for declaring and reading configuration variables.

To use `config`, create a `Config` object with `New`, and then call its
`NewInt`, `NewString`, or `NewBool` methods to declare configuration variables.
You can then call `Int`, `String`, or `Bool` to read the value of a variable.

## Motivation

Cloud-native applications are usually configured using environment variables.
However, accessing environment variables directly with `os.Getenv` is primitive
and error-prone. `config` provides a simple, type-safe interface for declaring
and reading configuration variables.

* Configuration variables are loaded, type-checked, and validated when they are
  registered, not when they are accessed. This ensures that configuration errors
  are caught early, instead of causing application errors, crashes, or
  misbehavior at runtime.
* Because validation happens at registration time, accessor methods like
  `Int` and `String` do not need to return an error. This simplifies the
  application code, and makes it easier to use configuration variables in
  places where error handling is not convenient.
* All configuration variables are declared in a single place, instead of
  scattered throughout the application code. This makes it easy to see all
  configuration variables, and to find the code that uses them.
* Description strings are attached to each configuration variable. This enables
  automatic documentation generation, and makes it easy to see what each
  variable is used for.

## Example

```go
package main

import (
    "fmt"
    "os"

    "github.com/neuralnorthwest/mu/config"
)

func main() {
    c := config.New()
    c.NewInt("PORT", 8080, "The port to listen on")
    c.NewString("MESSAGE", "Hello world!", "The message to print")
    c.NewBool("DEBUG", false, "Enable debug mode")

    fmt.Printf("Listening on port %d\n", c.Int("PORT"))
    fmt.Printf("Message: %s\n", c.String("MESSAGE"))
    if c.Bool("DEBUG") {
        fmt.Println("Debug mode is enabled")
    }
}
```

## Options

Various options are available for customizing the behavior of `config`. Some
options can be passed to `New`, while others are used with `NewInt`, `NewString`,
and `NewBool`.

### Options for `New`

* `WithEnvPrefix` - Sets the prefix for environment variables. For example, if
  the prefix is `MYAPP_`, then the environment variable `MYAPP_PORT` will be
  used to set the value of the `PORT` configuration variable. Example:

```go
package main

import (
    "fmt"
    "os"

    "github.com/neuralnorthwest/mu/config"
)

func main() {
    c := config.New(config.WithEnvPrefix("MYAPP_"))
    c.NewInt("PORT", 8080, "The port to listen on")
    fmt.Printf("Listening on port %d\n", c.Int("PORT"))
}
```

* `WithSource` - Sets the source for configuration variables. The default
  source is `os.Getenv`, but you can also use a custom source that reads
  configuration variables from a file, or from a database.

```go
package main

import (
    "fmt"
    "os"

    "github.com/neuralnorthwest/mu/config"
)

type customSource struct {
    prefix string
    // Other fields
}

func newCustomSource() *customSource {
    return &customSource{}
}

// Implement the Source interface
func (s *customSource) SetPrefix(prefix string) {
    s.prefix = prefix
}

func (s *customSource) LoadInt(name string) (int, error) {
    // Read the value of the configuration variable from a file or database
    value := ...
    return value, nil
}

func (s *customSource) LoadString(name string) (string, error) {
    // Read the value of the configuration variable from a file or database
    value := ...
    return value, nil
}

func (s *customSource) LoadBool(name string) (bool, error) {
    // Read the value of the configuration variable from a file or database
    value := ...
    return value, nil
}

func main() {
    c := config.New(config.WithSource(newCustomSource()))
    c.NewInt("PORT", 8080, "The port to listen on")
    fmt.Printf("Listening on port %d\n", c.Int("PORT"))
}
```

### Options for `NewInt`

* `WithMinimumValue` and `WithMaximumValue` - Sets the minimum and maximum
  allowed values for the configuration variable. If the value is outside the
  allowed range, `NewInt` will either return an error, or clamp the value to
  the nearest allowed value.
* `WithClamping` - Enables clamping for the configuration variable. If the
  value is outside the allowed range, it will be clamped to the nearest
  allowed value. If clamping is disabled, `NewInt` will return an error if the
  value is outside the allowed range.
* `WithIntValidator` - Sets a custom validator for the configuration variable.
  The validator is a function that takes an integer value and returns an error
  if the value is invalid. If the validator returns an error, `NewInt` will
  return the same error.
