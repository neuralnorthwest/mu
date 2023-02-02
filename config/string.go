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

// String is a configuration value that is a string.
type String struct {
	// name is the name of the variable.
	name string
	// value is the value of the variable.
	value string
	// defaultValue is the default value of the variable.
	defaultValue string
	// description is the description of the variable.
	description string
	// validator is the validator for the variable.
	validator func(string) error
}

// StringOption is an option for a string variable.
type StringOption func(*configImpl, *String) error

// WithStringValidator returns an option that sets the validator for a string variable.
func WithStringValidator(f func(string) error) StringOption {
	return func(c *configImpl, s *String) error {
		s.validator = f
		return nil
	}
}

// NewString creates a new string variable.
func (c *configImpl) NewString(name string, defaultValue string, description string, options ...StringOption) error {
	if _, ok := c.strings[name]; ok {
		return fmt.Errorf("%w: %s", status.ErrAlreadyExists, name)
	}
	s := &String{
		name:         name,
		value:        defaultValue,
		defaultValue: defaultValue,
		description:  description,
	}
	for _, o := range options {
		if err := o(c, s); err != nil {
			return err
		}
	}
	if s.validator != nil {
		if err := s.validator(s.value); err != nil {
			return err
		}
	}
	v, err := c.source.LoadString(name)
	if err == nil {
		if s.validator != nil {
			if err := s.validator(v); err != nil {
				return err
			}
		}
		s.value = v
	} else if !errors.Is(err, status.ErrNotFound) {
		return err
	}
	c.strings[name] = s
	return nil
}

// String returns the value of the string variable with the given name. If the variable does not exist, it calls bug.Bug.
// If bug.Bug does not panic, String returns "".
func (c *configImpl) String(name string) string {
	if s, ok := c.strings[name]; ok {
		return s.value
	}
	defer bug.Bugf("config: string variable %q does not exist", name)
	return ""
}

// DescribeString returns the description of the string variable with the given name. If the variable does not exist, it calls bug.Bug.
// If bug.Bug does not panic, DescribeString returns "".
func (c *configImpl) DescribeString(name string) string {
	if s, ok := c.strings[name]; ok {
		return s.description
	}
	defer bug.Bugf("config: string variable %q does not exist", name)
	return ""
}
