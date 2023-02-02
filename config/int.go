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
	"math"

	"github.com/neuralnorthwest/mu/bug"
	"github.com/neuralnorthwest/mu/status"
)

// Int is a configuration value that is an integer.
type Int struct {
	// name is the name of the variable.
	name string
	// value is the value of the variable.
	value int
	// defaultValue is the default value of the variable.
	defaultValue int
	// description is the description of the variable.
	description string
	// minimumValue is the minimum value of the variable.
	minimumValue int
	// maximumValue is the maximum value of the variable.
	maximumValue int
	// clampValue is true if the value should be clamped to the minimum and maximum values.
	clampValue bool
	// validator is the validator for the variable.
	validator func(int) error
}

// IntOption is an option for an int variable.
type IntOption func(*configImpl, *Int) error

// WithMinimumValue returns an option that sets the minimum value for an int variable.
func WithMinimumValue(min int) IntOption {
	return func(c *configImpl, i *Int) error {
		i.minimumValue = min
		return nil
	}
}

// WithMaximumValue returns an option that sets the maximum value for an int variable.
func WithMaximumValue(max int) IntOption {
	return func(c *configImpl, i *Int) error {
		i.maximumValue = max
		return nil
	}
}

// WithIntValidator returns an option that sets the validator for an int variable.
func WithIntValidator(f func(int) error) IntOption {
	return func(c *configImpl, i *Int) error {
		i.validator = f
		return nil
	}
}

// WithClamping returns an option that enables clamping for an int variable.
func WithClamping() IntOption {
	return func(c *configImpl, i *Int) error {
		i.clampValue = true
		return nil
	}
}

// NewInt creates a new int variable.
func (c *configImpl) NewInt(name string, defaultValue int, description string, options ...IntOption) error {
	if _, ok := c.ints[name]; ok {
		return fmt.Errorf("%w: %s", status.ErrAlreadyExists, name)
	}
	i := &Int{
		name:         name,
		defaultValue: defaultValue,
		description:  description,
		minimumValue: math.MinInt32,
		maximumValue: math.MaxInt32,
	}
	for _, opt := range options {
		if err := opt(c, i); err != nil {
			return err
		}
	}
	if i.minimumValue > i.maximumValue {
		return status.ErrInvalidRange
	}
	var err error
	i.defaultValue, err = clampOrError(i.defaultValue, i.minimumValue, i.maximumValue, i.clampValue)
	if err != nil {
		return err
	}
	if i.validator != nil {
		if err := i.validator(i.defaultValue); err != nil {
			return err
		}
	}
	i.value = i.defaultValue
	v, err := c.source.LoadInt(name)
	if err == nil {
		v, err = clampOrError(v, i.minimumValue, i.maximumValue, i.clampValue)
		if err != nil {
			return err
		}
		if i.validator != nil {
			if err := i.validator(v); err != nil {
				return err
			}
		}
		i.value = v
	} else if !errors.Is(err, status.ErrNotFound) {
		return err
	}
	c.ints[name] = i
	return nil
}

// Int returns the value of the int variable with the given name. If the variable does not exist, it calls bug.Bug.
// If bug.Bug does not panic, Int returns 0.
func (c *configImpl) Int(name string) int {
	if i, ok := c.ints[name]; ok {
		return i.value
	}
	defer bug.Bugf("config: int variable %q does not exist", name)
	return 0
}

// DescribeInt returns the description of the int variable with the given name. If the variable does not exist, it calls bug.Bug.
// If bug.Bug does not panic, DescribeInt returns an empty string.
func (c *configImpl) DescribeInt(name string) string {
	if i, ok := c.ints[name]; ok {
		return i.description
	}
	defer bug.Bugf("config: int variable %q does not exist", name)
	return ""
}

// clampOrError clamps the value to the given range, or returns an error if the value is out of range.
func clampOrError(value, min, max int, clamp bool) (int, error) {
	if value < min {
		if clamp {
			return min, nil
		}
		return value, status.ErrOutOfRange
	}
	if value > max {
		if clamp {
			return max, nil
		}
		return value, status.ErrOutOfRange
	}
	return value, nil
}
