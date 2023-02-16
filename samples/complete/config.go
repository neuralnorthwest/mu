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

package complete

import "github.com/neuralnorthwest/mu/config"

const MessageConfigVar = "MESSAGE"
const DevModeConfigVar = "DEV_MODE"

// Config is a convenience wrapper around the service configuration.
type Config struct {
	// config is the service configuration.
	config config.Config
}

// setupConfig sets up the service configuration. This is where we register
// configuration variables.
func (s *complete) setupConfig(c config.Config) error {
	// Register a configuration variable. The first argument is the name of the
	// variable. The second argument is the default value. The third argument is
	// the description of the variable.
	//
	// It's good practice to use constants for the variable names. This makes it
	// easier to find all references to a variable, as well as preventing
	// typos.
	err := c.NewString(MessageConfigVar, "Hello, World!", "The message to display.")
	// The NewString, NewInt, and NewBool methods return an error that must be
	// checked. An error could occur if the variable has already been
	// registered, the default value is invalid, or the value set in the
	// environment is invalid.
	if err != nil {
		return err
	}
	// DEV_MODE enables development mode. The application will behave slightly
	// differently in development mode. For instance:
	//
	//   - The bug handler will print bugs to the console and panic.
	err = c.NewBool(DevModeConfigVar, false, "Enable development mode.")
	if err != nil {
		return err
	}
	s.config = &Config{
		config: c,
	}
	return nil
}

// Message returns the message to display.
func (c *Config) Message() string {
	return c.config.String(MessageConfigVar)
}

// DevMode returns true if development mode is enabled.
func (c *Config) DevMode() bool {
	return c.config.Bool(DevModeConfigVar)
}
