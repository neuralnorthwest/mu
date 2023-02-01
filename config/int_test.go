package config

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/neuralnorthwest/mu/bug"
	"github.com/neuralnorthwest/mu/status"
)

// Test_NewInt_Case is a test case for NewInt.
type Test_NewInt_Case struct {
	// name is the name of the test case.
	name string
	// varName is the name of the variable.
	varName string
	// defaultValue is the default value of the variable.
	defaultValue int
	// description is the description of the variable.
	description string
	// options are the options for the variable.
	options []IntOption
	// expected is the expected result.
	expected *Int
	// err is the expected error.
	err error
}

// validatorIsNonNegative is an int validator that returns an error if the value is negative.
func validatorIsNonNegative(value int) error {
	if value < 0 {
		return fmt.Errorf("value must be greater than or equal to 0: %w", status.ErrOutOfRange)
	}
	return nil
}

// intOptionError return an IntOption that returns an error.
func intOptionError(err error) IntOption {
	return func(c *configImpl, i *Int) error {
		return err
	}
}

// Test_NewInt tests the NewInt function.
func Test_NewInt(t *testing.T) {
	t.Parallel()
	for _, tc := range []Test_NewInt_Case{
		{
			name:         "basic case",
			varName:      "test",
			defaultValue: 1,
			description:  "test",
			options:      []IntOption{},
			expected: &Int{
				name:         "test",
				value:        1,
				defaultValue: 1,
				minimumValue: math.MinInt32,
				maximumValue: math.MaxInt32,
				description:  "test",
			},
			err: nil,
		},
		{
			name:         "with minimum value",
			varName:      "test",
			defaultValue: 1,
			description:  "test",
			options: []IntOption{
				WithMinimumValue(0),
			},
			expected: &Int{
				name:         "test",
				value:        1,
				defaultValue: 1,
				minimumValue: 0,
				maximumValue: math.MaxInt32,
				description:  "test",
			},
			err: nil,
		},
		{
			name:         "with minimum value, default value less than minimum value",
			varName:      "test",
			defaultValue: -1,
			description:  "test",
			options: []IntOption{
				WithMinimumValue(0),
			},
			expected: nil,
			err:      status.ErrOutOfRange,
		},
		{
			name:         "with minimum value, default value less than minimum value, with clamping",
			varName:      "test",
			defaultValue: -1,
			description:  "test",
			options: []IntOption{
				WithMinimumValue(0),
				WithClamping(),
			},
			expected: &Int{
				name:         "test",
				defaultValue: 0,
				minimumValue: 0,
				maximumValue: math.MaxInt32,
				clampValue:   true,
				description:  "test",
			},
			err: nil,
		},
		{
			name:         "with minimum value, default value equal to minimum value",
			varName:      "test",
			defaultValue: 0,
			description:  "test",
			options: []IntOption{
				WithMinimumValue(0),
			},
			expected: &Int{
				name:         "test",
				defaultValue: 0,
				minimumValue: 0,
				maximumValue: math.MaxInt32,
				description:  "test",
			},
			err: nil,
		},
		{
			name:         "with maximum value",
			varName:      "test",
			defaultValue: 1,
			description:  "test",
			options: []IntOption{
				WithMaximumValue(2),
			},
			expected: &Int{
				name:         "test",
				value:        1,
				defaultValue: 1,
				minimumValue: math.MinInt32,
				maximumValue: 2,
				description:  "test",
			},
			err: nil,
		},
		{
			name:         "with maximum value, default value greater than maximum value",
			varName:      "test",
			defaultValue: 3,
			description:  "test",
			options: []IntOption{
				WithMaximumValue(2),
			},
			expected: nil,
			err:      status.ErrOutOfRange,
		},
		{
			name:         "with maximum value, default value greater than maximum value, with clamping",
			varName:      "test",
			defaultValue: 3,
			description:  "test",
			options: []IntOption{
				WithMaximumValue(2),
				WithClamping(),
			},
			expected: &Int{
				name:         "test",
				value:        2,
				defaultValue: 2,
				minimumValue: math.MinInt32,
				maximumValue: 2,
				clampValue:   true,
				description:  "test",
			},
			err: nil,
		},
		{
			name:         "with maximum value, default value equal to maximum value",
			varName:      "test",
			defaultValue: 2,
			description:  "test",
			options: []IntOption{
				WithMaximumValue(2),
			},
			expected: &Int{
				name:         "test",
				value:        2,
				defaultValue: 2,
				minimumValue: math.MinInt32,
				maximumValue: 2,
				description:  "test",
			},
			err: nil,
		},
		{
			name:         "with validator",
			varName:      "test",
			defaultValue: 1,
			description:  "test",
			options: []IntOption{
				WithIntValidator(validatorIsNonNegative),
			},
			expected: &Int{
				name:         "test",
				value:        1,
				defaultValue: 1,
				minimumValue: math.MinInt32,
				maximumValue: math.MaxInt32,
				description:  "test",
			},
			err: nil,
		},
		{
			name:         "with validator, default value invalid",
			varName:      "test",
			defaultValue: -1,
			description:  "test",
			options: []IntOption{
				WithIntValidator(validatorIsNonNegative),
			},
			expected: nil,
			err:      status.ErrOutOfRange,
		},
		{
			name:         "minimum value greater than maximum value",
			varName:      "test",
			defaultValue: 1,
			description:  "test",
			options: []IntOption{
				WithMinimumValue(2),
				WithMaximumValue(1),
			},
			expected: nil,
			err:      status.ErrInvalidRange,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Use a null source to ensure we don't read anything from the environment.
			c := New(WithSource(newNullSource())).(*configImpl)
			err := c.NewInt(tc.varName, tc.defaultValue, tc.description, tc.options...)
			if err != nil {
				if tc.err == nil {
					t.Errorf("NewInt() = %v; want nil", err)
				} else if !errors.Is(err, tc.err) {
					t.Errorf("NewInt() = %v; want %v", err, tc.err)
				}
				return
			}
			if tc.err != nil {
				t.Errorf("NewInt() = nil; want %v", tc.err)
				return
			}
			if c.ints[tc.varName] == nil {
				t.Errorf("NewInt() = nil; want %v", tc.expected)
				return
			}
			// get diffs, allowing unexported fields to be compared
			opts := []cmp.Option{
				cmp.AllowUnexported(Int{}),
				cmpopts.IgnoreFields(Int{}, "validator"),
			}
			diffs := cmp.Diff(tc.expected, c.ints[tc.varName], opts...)
			if diffs != "" {
				t.Errorf("NewInt() = %v; want %v; diffs: %v", c.ints[tc.varName], tc.expected, diffs)
			}
		})
	}
}

// Test_NewInt_AlreadyExists tests that NewInt() returns an error if the variable
// name already exists.
func Test_NewInt_AlreadyExists(t *testing.T) {
	t.Parallel()
	c := New().(*configImpl)
	if err := c.NewInt("test", 1, "test"); err != nil {
		t.Fatalf("NewInt() = %v; want nil", err)
	}
	if err := c.NewInt("test", 1, "test"); !errors.Is(err, status.ErrAlreadyExists) {
		t.Errorf("NewInt() = %v; want %v", err, status.ErrAlreadyExists)
	}
}

// Test_NewInt_OptionError tests that NewInt() returns an error if an option
// returns an error.
func Test_NewInt_OptionError(t *testing.T) {
	t.Parallel()
	c := New().(*configImpl)
	if err := c.NewInt("test", 1, "test", intOptionError(status.ErrOutOfRange)); !errors.Is(err, status.ErrOutOfRange) {
		t.Errorf("NewInt() = %v; want %v", err, status.ErrOutOfRange)
	}
}

// Test_Int tests that Int() returns the correct value.
func Test_Int(t *testing.T) {
	t.Parallel()
	c := New().(*configImpl)
	if err := c.NewInt("test", 1, "test"); err != nil {
		t.Fatalf("NewInt() = %v; want nil", err)
	}
	if v := c.Int("test"); v != 1 {
		t.Errorf("Int() = %v; want 1", v)
	}
}

// Test_Int_FromTestSource tests that Int() returns the correct value when
// the variable is set in the test source.
func Test_Int_FromTestSource(t *testing.T) {
	t.Parallel()
	src := newTestSource()
	src.SetValue("test", "2")
	c := New(WithSource(src)).(*configImpl)
	if err := c.NewInt("test", 1, "test"); err != nil {
		t.Fatalf("NewInt() = %v; want nil", err)
	}
	if v := c.Int("test"); v != 2 {
		t.Errorf("Int() = %v; want 2", v)
	}
}

// Test_Int_FromTestSource_OutOfRange tests that Int() returns an error when
// the variable is set in the test source and the value is out of range.
func Test_Int_FromTestSource_OutOfRange(t *testing.T) {
	t.Parallel()
	src := newTestSource()
	src.SetValue("test", "2")
	c := New(WithSource(src)).(*configImpl)
	if err := c.NewInt("test", 3, "test", WithMinimumValue(3)); !errors.Is(err, status.ErrOutOfRange) {
		t.Errorf("NewInt() = %v; want %v", err, status.ErrOutOfRange)
	}
}

// Test_Int_FromTestSource_Validation tests that Int() returns an error
// when the variable is set in the test source and the value fails validation.
func Test_Int_FromTestSource_Validation(t *testing.T) {
	t.Parallel()
	src := newTestSource()
	src.SetValue("test", "1")
	src.SetValue("test2", "-1")
	validatorOpt := WithIntValidator(func(v int) error {
		if v < 0 {
			return status.ErrInvalidArgument
		}
		return nil
	})
	c := New(WithSource(src)).(*configImpl)
	if err := c.NewInt("test", 0, "test", validatorOpt); err != nil {
		t.Errorf("NewInt() = %v; want nil", err)
	}
	if err := c.NewInt("test2", 0, "test", validatorOpt); !errors.Is(err, status.ErrInvalidArgument) {
		t.Errorf("NewInt() = %v; want %v", err, status.ErrInvalidArgument)
	}
}

// Test_Int_FromErrorSource tests that Int() returns an error when the source
// returns an error when getting the value.
func Test_Int_FromErrorSource(t *testing.T) {
	t.Parallel()
	src := newErrorSource(status.ErrInvalidArgument)
	c := New(WithSource(src)).(*configImpl)
	if err := c.NewInt("test", 1, "test"); err != status.ErrInvalidArgument {
		t.Errorf("NewInt() = %v; want %v", err, status.ErrInvalidArgument)
	}
}

// Test_Int_FromEnv tests that Int() returns the correct value when the
// variable is set in the environment.
func Test_Int_FromEnv(t *testing.T) {
	t.Setenv("test", "2")
	c := New().(*configImpl)
	if err := c.NewInt("test", 1, "test"); err != nil {
		t.Fatalf("NewInt() = %v; want nil", err)
	}
	if v := c.Int("test"); v != 2 {
		t.Errorf("Int() = %v; want 2", v)
	}
}

// Test_Int_FromEnv_LoadPrefix tests that Int() returns the correct value
// when the variable is set in the environment with a load prefix.
func Test_Int_FromEnv_LoadPrefix(t *testing.T) {
	t.Setenv("prefix-test", "2")
	c := New(WithLoadPrefix("prefix-")).(*configImpl)
	if err := c.NewInt("test", 1, "test"); err != nil {
		t.Fatalf("NewInt() = %v; want nil", err)
	}
	if v := c.Int("test"); v != 2 {
		t.Errorf("Int() = %v; want 2", v)
	}
}

// Test_Int_Unknown tests that Int() calls bug.Bug if the variable name is unknown.
func Test_Int_Unknown(t *testing.T) {
	message := ""
	expectedMessage := "config: int variable \"test\" does not exist"
	oldHandler := bug.Handler()
	defer bug.SetHandler(oldHandler)
	bug.SetHandler(func(msg string) {
		message = msg
	})
	c := New().(*configImpl)
	_ = c.Int("test")
	if message != expectedMessage {
		t.Errorf("bug.Bug() = %v; want %v", message, expectedMessage)
	}
}

// Test_DescribeInt tests that DescribeInt() returns the correct value.
func Test_DescribeInt(t *testing.T) {
	t.Parallel()
	c := New().(*configImpl)
	if err := c.NewInt("test", 1, "test"); err != nil {
		t.Fatalf("NewInt() = %v; want nil", err)
	}
	if v := c.DescribeInt("test"); v != "test" {
		t.Errorf("DescribeInt() = %v; want \"test\"", v)
	}
}

// Test_DescribeInt_Unknown tests that DescribeInt() calls bug.Bug if the variable
// name is unknown.
func Test_DescribeInt_Unknown(t *testing.T) {
	message := ""
	expectedMessage := "config: int variable \"test\" does not exist"
	oldHandler := bug.Handler()
	defer bug.SetHandler(oldHandler)
	bug.SetHandler(func(msg string) {
		message = msg
	})
	c := New().(*configImpl)
	_ = c.DescribeInt("test")
	if message != expectedMessage {
		t.Errorf("bug.Bug() = %v; want %v", message, expectedMessage)
	}
}

// Test_Int_clampOrError_Case is a test case for Test_Int_clampOrError.
type Test_Int_clampOrError_Case struct {
	name         string
	value        int
	minimumValue int
	maximumValue int
	clamp        bool
	expected     int
	err          error
}

// Test_Int_clampOrError tests that clampOrError() returns the correct value.
func Test_Int_clampOrError(t *testing.T) {
	t.Parallel()
	for _, tc := range []Test_Int_clampOrError_Case{
		{
			name:         "in range",
			value:        1,
			minimumValue: 0,
			maximumValue: 2,
			clamp:        false,
			expected:     1,
			err:          nil,
		},
		{
			name:         "below min",
			value:        -1,
			minimumValue: 0,
			maximumValue: 2,
			clamp:        false,
			expected:     -1,
			err:          status.ErrOutOfRange,
		},
		{
			name:         "at min",
			value:        0,
			minimumValue: 0,
			maximumValue: 2,
			clamp:        false,
			expected:     0,
			err:          nil,
		},
		{
			name:         "above max",
			value:        3,
			minimumValue: 0,
			maximumValue: 2,
			clamp:        false,
			expected:     3,
			err:          status.ErrOutOfRange,
		},
		{
			name:         "at max",
			value:        2,
			minimumValue: 0,
			maximumValue: 2,
			clamp:        false,
			expected:     2,
			err:          nil,
		},
		{
			name:         "below min, clamp",
			value:        -1,
			minimumValue: 0,
			maximumValue: 2,
			clamp:        true,
			expected:     0,
			err:          nil,
		},
		{
			name:         "at min, clamp",
			value:        0,
			minimumValue: 0,
			maximumValue: 2,
			clamp:        true,
			expected:     0,
			err:          nil,
		},
		{
			name:         "above max, clamp",
			value:        3,
			minimumValue: 0,
			maximumValue: 2,
			clamp:        true,
			expected:     2,
			err:          nil,
		},
		{
			name:         "at max, clamp",
			value:        2,
			minimumValue: 0,
			maximumValue: 2,
			clamp:        true,
			expected:     2,
			err:          nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v, err := clampOrError(tc.value, tc.minimumValue, tc.maximumValue, tc.clamp)
			if err != nil {
				if tc.err == nil {
					t.Errorf("clampOrError() = %v; want nil", err)
				} else if !errors.Is(err, tc.err) {
					t.Errorf("clampOrError() = %v; want %v", err, tc.err)
				}
				return
			}
			if tc.err != nil {
				t.Errorf("clampOrError() = nil; want %v", tc.err)
				return
			}
			if v != tc.expected {
				t.Errorf("clampOrError() = %v; want %v", v, tc.expected)
			}
		})
	}
}
