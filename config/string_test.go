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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/neuralnorthwest/mu/bug"
	"github.com/neuralnorthwest/mu/status"
)

// Test_NewString_Case is a test case for the Test_NewString function.
type Test_NewString_Case struct {
	// name is the name of the test case.
	name string
	// varName is the name of the variable.
	varName string
	// defaultValue is the default value of the variable.
	defaultValue string
	// description is the description of the variable.
	description string
	// options are the options for the variable.
	options []StringOption
	// expected is the expected value of the variable.
	expected *String
	// err is the expected error.
	err error
}

// stringOptionError return an StringOption that returns an error.
func stringOptionError(err error) StringOption {
	return func(c *configImpl, s *String) error {
		return err
	}
}

// Test_NewString tests the NewString function.
func Test_NewString(t *testing.T) {
	for _, tc := range []Test_NewString_Case{
		{
			name:         "basic case",
			varName:      "test",
			defaultValue: "test",
			description:  "test",
			options:      []StringOption{},
			expected: &String{
				name:         "test",
				value:        "test",
				defaultValue: "test",
				description:  "test",
			},
			err: nil,
		},
		{
			name:         "with validator",
			varName:      "test",
			defaultValue: "test",
			description:  "test",
			options: []StringOption{
				WithStringValidator(func(s string) error {
					if s != "test" {
						return status.ErrInvalidArgument
					}
					return nil
				}),
			},
			expected: &String{
				name:         "test",
				value:        "test",
				defaultValue: "test",
				description:  "test",
			},
			err: nil,
		},
		{
			name:         "with validator, invalid value",
			varName:      "test",
			defaultValue: "invalid",
			description:  "test",
			options: []StringOption{
				WithStringValidator(func(s string) error {
					if s != "test" {
						return status.ErrInvalidArgument
					}
					return nil
				}),
			},
			expected: nil,
			err:      status.ErrInvalidArgument,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := New().(*configImpl)
			err := c.NewString(tc.varName, tc.defaultValue, tc.description, tc.options...)
			if err != nil {
				if tc.err == nil {
					t.Errorf("unexpected error: %v", err)
				} else if !errors.Is(err, tc.err) {
					t.Errorf("unexpected error: %v, expected: %v", err, tc.err)
				}
				return
			}
			if tc.err != nil {
				t.Errorf("expected error: %v", tc.err)
				return
			}
			s, ok := c.strings[tc.varName]
			if !ok {
				t.Errorf("variable not found: %v", tc.varName)
				return
			}
			opts := []cmp.Option{
				cmp.AllowUnexported(String{}),
				cmpopts.IgnoreFields(String{}, "validator"),
			}
			if diff := cmp.Diff(tc.expected, s, opts...); diff != "" {
				t.Errorf("unexpected variable (-want +got):\n%s", diff)
			}
		})
	}
}

// Test_NewString_AlreadyExists tests that NewString() returns an error if the variable
// name already exists.
func Test_NewString_AlreadyExists(t *testing.T) {
	t.Parallel()
	c := New().(*configImpl)
	if err := c.NewString("test", "test", "test"); err != nil {
		t.Fatalf("NewString() = %v; want nil", err)
	}
	if err := c.NewString("test", "test", "test"); !errors.Is(err, status.ErrAlreadyExists) {
		t.Errorf("NewString() = %v; want %v", err, status.ErrAlreadyExists)
	}
}

// Test_NewString_OptionError tests that NewString() returns an error if an option
// returns an error.
func Test_NewString_OptionError(t *testing.T) {
	t.Parallel()
	c := New().(*configImpl)
	if err := c.NewString("test", "test", "test", stringOptionError(status.ErrInvalidArgument)); !errors.Is(err, status.ErrInvalidArgument) {
		t.Errorf("NewString() = %v; want %v", err, status.ErrInvalidArgument)
	}
}

// Test_String tests that String() returns the correct value.
func Test_String(t *testing.T) {
	t.Parallel()
	c := New().(*configImpl)
	if err := c.NewString("test", "test", "test"); err != nil {
		t.Fatalf("NewString() = %v; want nil", err)
	}
	if v := c.String("test"); v != "test" {
		t.Errorf("String() = %v; want test", v)
	}
}

// Test_String_FromTestSource tests that String() returns the correct value when
// the variable is set in the test source.
func Test_String_FromTestSource(t *testing.T) {
	t.Parallel()
	src := newTestSource()
	src.SetValue("test", "test")
	c := New(WithSource(src)).(*configImpl)
	if err := c.NewString("test", "xxx", "test"); err != nil {
		t.Fatalf("NewString() = %v; want nil", err)
	}
	if v := c.String("test"); v != "test" {
		t.Errorf("String() = %v; want test", v)
	}
}

// Test_String_FromTestSource_Validation tests that String() returns an error
// when the variable is set in the test source and the value fails validation.
func Test_String_FromTestSource_Validation(t *testing.T) {
	t.Parallel()
	src := newTestSource()
	src.SetValue("test", "invalid")
	src.SetValue("test2", "test")
	validatorOption := WithStringValidator(func(s string) error {
		if s != "test" {
			return status.ErrInvalidArgument
		}
		return nil
	})
	c := New(WithSource(src)).(*configImpl)
	if err := c.NewString("test", "test", "test", validatorOption); !errors.Is(err, status.ErrInvalidArgument) {
		t.Fatalf("NewString() = nil; want %v", status.ErrInvalidArgument)
	}
	if err := c.NewString("test2", "test", "test", validatorOption); err != nil {
		t.Fatalf("NewString() = %v; want nil", err)
	}
}

// Test_String_FromErrorSource tests that String() returns the correct value when
// the variable is set in the error source.
func Test_String_FromErrorSource(t *testing.T) {
	t.Parallel()
	src := newErrorSource(status.ErrInvalidArgument)
	c := New(WithSource(src)).(*configImpl)
	if err := c.NewString("test", "xxx", "test"); err == nil {
		t.Fatalf("NewString() = nil; want %v", status.ErrInvalidArgument)
	}
}

// Test_String_FromEnv tests that String() returns the correct value when
// the variable is set in the environment.
func Test_String_FromEnv(t *testing.T) {
	t.Setenv("test", "test")
	c := New().(*configImpl)
	if err := c.NewString("test", "xxx", "test"); err != nil {
		t.Fatalf("NewString() = %v; want nil", err)
	}
	if v := c.String("test"); v != "test" {
		t.Errorf("String() = %v; want test", v)
	}
}

// Test_String_FromEnv_LoadPrefix tests that String() returns the correct value when
// the variable is set in the environment with a load prefix.
func Test_String_FromEnv_LoadPrefix(t *testing.T) {
	t.Setenv("prefix-test", "test")
	c := New(WithLoadPrefix("prefix-")).(*configImpl)
	if err := c.NewString("test", "xxx", "test"); err != nil {
		t.Fatalf("NewString() = %v; want nil", err)
	}
	if v := c.String("test"); v != "test" {
		t.Errorf("String() = %v; want test", v)
	}
}

// Test_String_Unknown tests that String() calls bug.Bug if the variable name is unknown.
func Test_String_Unknown(t *testing.T) {
	message := ""
	expectedMessage := "config: string variable \"test\" does not exist"
	oldHandler := bug.Handler()
	defer bug.SetHandler(oldHandler)
	bug.SetHandler(func(msg string) {
		message = msg
	})
	c := New().(*configImpl)
	_ = c.String("test")
	if message != expectedMessage {
		t.Errorf("bug.Bug() = %v; want %v", message, expectedMessage)
	}
}

// Test_DescribeString tests that DescribeString() returns the correct value.
func Test_DescribeString(t *testing.T) {
	t.Parallel()
	c := New().(*configImpl)
	if err := c.NewString("test", "test", "test"); err != nil {
		t.Fatalf("NewString() = %v; want nil", err)
	}
	if v := c.DescribeString("test"); v != "test" {
		t.Errorf("DescribeString() = %v; want \"test\"", v)
	}
}

// Test_DescribeString_Unknown tests that DescribeString() calls bug.Bug if the variable
// name is unknown.
func Test_DescribeString_Unknown(t *testing.T) {
	message := ""
	expectedMessage := "config: string variable \"test\" does not exist"
	oldHandler := bug.Handler()
	defer bug.SetHandler(oldHandler)
	bug.SetHandler(func(msg string) {
		message = msg
	})
	c := New().(*configImpl)
	_ = c.DescribeString("test")
	if message != expectedMessage {
		t.Errorf("bug.Bug() = %v; want %v", message, expectedMessage)
	}
}
