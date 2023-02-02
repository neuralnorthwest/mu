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

package config

import (
	"errors"
	"fmt"

	"github.com/neuralnorthwest/mu/bug"
	"github.com/neuralnorthwest/mu/status"
)

// Bool is a configuration value that is a bool.
type Bool struct {
	// name is the name of the variable.
	name string
	// value is the value of the variable.
	value bool
	// defaultValue is the default value of the variable.
	defaultValue bool
	// description is the description of the variable.
	description string
}

// NewBool creates a new bool variable.
func (c *configImpl) NewBool(name string, defaultValue bool, description string) error {
	if _, ok := c.bools[name]; ok {
		return fmt.Errorf("%w: %s", status.ErrAlreadyExists, name)
	}
	b := &Bool{
		name:         name,
		value:        defaultValue,
		defaultValue: defaultValue,
		description:  description,
	}
	v, err := c.source.LoadBool(name)
	if err == nil {
		b.value = v
	}
	if err != nil && !errors.Is(err, status.ErrNotFound) {
		return err
	}
	c.bools[name] = b
	return nil
}

// Bool returns the value of the bool variable with the given name. If the variable does not exist, it calls bug.Bug.
// If bug.Bug does not panic, Bool returns false.
func (c *configImpl) Bool(name string) bool {
	if b, ok := c.bools[name]; ok {
		return b.value
	}
	defer bug.Bugf("config: bool variable %q does not exist", name)
	return false
}

// DescribeBool returns the description of the bool variable with the given name. If the variable does not exist, it calls bug.Bug.
// If bug.Bug does not panic, DescribeBool returns "".
func (c *configImpl) DescribeBool(name string) string {
	if b, ok := c.bools[name]; ok {
		return b.description
	}
	defer bug.Bugf("config: bool variable %q does not exist", name)
	return ""
}
