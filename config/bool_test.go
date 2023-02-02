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
	"github.com/neuralnorthwest/mu/bug"
	"github.com/neuralnorthwest/mu/status"
)

// Test_NewBool_Case is a test case for the Test_NewBool function.
type Test_NewBool_Case struct {
	// name is the name of the test case.
	name string
	// varName is the name of the variable.
	varName string
	// defaultValue is the default value of the variable.
	defaultValue bool
	// description is the description of the variable.
	description string
	// expected is the expected value of the variable.
	expected *Bool
	// err is the expected error.
	err error
}

// Test_NewBool tests the NewBool function.
func Test_NewBool(t *testing.T) {
	for _, tc := range []Test_NewBool_Case{
		{
			name:         "default false",
			varName:      "test",
			defaultValue: false,
			description:  "test",
			expected: &Bool{
				name:         "test",
				value:        false,
				defaultValue: false,
				description:  "test",
			},
			err: nil,
		},
		{
			name:         "default true",
			varName:      "test",
			defaultValue: true,
			description:  "test",
			expected: &Bool{
				name:         "test",
				value:        true,
				defaultValue: true,
				description:  "test",
			},
			err: nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := New().(*configImpl)
			err := c.NewBool(tc.varName, tc.defaultValue, tc.description)
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
			b, ok := c.bools[tc.varName]
			if !ok {
				t.Errorf("variable not found: %v", tc.varName)
				return
			}
			opts := []cmp.Option{
				cmp.AllowUnexported(Bool{}),
			}
			if diff := cmp.Diff(tc.expected, b, opts...); diff != "" {
				t.Errorf("unexpected variable (-want +got):\n%s", diff)
			}
		})
	}
}

// Test_NewBool_AlreadyExists tests that NewBool() returns an error if the variable
// name already exists.
func Test_NewBool_AlreadyExists(t *testing.T) {
	t.Parallel()
	c := New().(*configImpl)
	if err := c.NewBool("test", false, "test"); err != nil {
		t.Fatalf("NewBool() = %v; want nil", err)
	}
	if err := c.NewBool("test", false, "test"); !errors.Is(err, status.ErrAlreadyExists) {
		t.Errorf("NewBool() = %v; want %v", err, status.ErrAlreadyExists)
	}
}

// Test_Bool tests that Bool() returns the correct value.
func Test_Bool(t *testing.T) {
	t.Parallel()
	c := New().(*configImpl)
	if err := c.NewBool("test", true, "test"); err != nil {
		t.Fatalf("NewBool() = %v; want nil", err)
	}
	if v := c.Bool("test"); !v {
		t.Errorf("Bool() = %v; want true", v)
	}
}

// Test_Bool_FromTestSource tests that Bool() returns the correct value when the
// variable is set in the test source.
func Test_Bool_FromTestSource(t *testing.T) {
	t.Parallel()
	src := newTestSource()
	src.SetValue("test", "true")
	c := New(WithSource(src)).(*configImpl)
	if err := c.NewBool("test", false, "test"); err != nil {
		t.Fatalf("NewBool() = %v; want nil", err)
	}
	if v := c.Bool("test"); !v {
		t.Errorf("Bool() = %v; want true", v)
	}
}

// Test_Bool_FromErrorSource tests that Bool() returns the correct value when
// the variable is set in the error source.
func Test_Bool_FromErrorSource(t *testing.T) {
	t.Parallel()
	src := newErrorSource(status.ErrInvalidArgument)
	c := New(WithSource(src)).(*configImpl)
	if err := c.NewBool("test", false, "test"); !errors.Is(err, status.ErrInvalidArgument) {
		t.Errorf("NewBool() = %v; want %v", err, status.ErrInvalidArgument)
	}
}

// Test_Bool_FromEnv tests that Bool() returns the correct value when the
// variable is set in the environment.
func Test_Bool_FromEnv(t *testing.T) {
	t.Setenv("test", "true")
	c := New().(*configImpl)
	if err := c.NewBool("test", false, "test"); err != nil {
		t.Fatalf("NewBool() = %v; want nil", err)
	}
	if v := c.Bool("test"); !v {
		t.Errorf("Bool() = %v; want true", v)
	}
}

// Test_Bool_FromEnv_LoadPrefix tests that Bool() returns the correct value when
// the variable is set in the environment with a load prefix.
func Test_Bool_FromEnv_LoadPrefix(t *testing.T) {
	t.Setenv("prefix-test", "true")
	c := New(WithLoadPrefix("prefix-")).(*configImpl)
	if err := c.NewBool("test", false, "test"); err != nil {
		t.Fatalf("NewBool() = %v; want nil", err)
	}
	if v := c.Bool("test"); !v {
		t.Errorf("Bool() = %v; want true", v)
	}
}

// Test_Bool_Unknown tests that Bool() calls bug.Bug if the variable name is unknown.
func Test_Bool_Unknown(t *testing.T) {
	message := ""
	expectedMessage := "config: bool variable \"test\" does not exist"
	oldHandler := bug.Handler()
	defer bug.SetHandler(oldHandler)
	bug.SetHandler(func(msg string) {
		message = msg
	})
	c := New().(*configImpl)
	_ = c.Bool("test")
	if message != expectedMessage {
		t.Errorf("bug.Bug() = %v; want %v", message, expectedMessage)
	}
}

// Test_DescribeBool tests that DescribeBool() returns the correct value.
func Test_DescribeBool(t *testing.T) {
	t.Parallel()
	c := New().(*configImpl)
	if err := c.NewBool("test", true, "test"); err != nil {
		t.Fatalf("NewBool() = %v; want nil", err)
	}
	if v := c.DescribeBool("test"); v != "test" {
		t.Errorf("DescribeBool() = %v; want test", v)
	}
}

// Test_DescribeBool_Unknown tests that DescribeBool() calls bug.Bug if the variable
// name is unknown.
func Test_DescribeBool_Unknown(t *testing.T) {
	message := ""
	expectedMessage := "config: bool variable \"test\" does not exist"
	oldHandler := bug.Handler()
	defer bug.SetHandler(oldHandler)
	bug.SetHandler(func(msg string) {
		message = msg
	})
	c := New().(*configImpl)
	_ = c.DescribeBool("test")
	if message != expectedMessage {
		t.Errorf("bug.Bug() = %v; want %v", message, expectedMessage)
	}
}
