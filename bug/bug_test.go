package bug

import "testing"

// Test_Bugf_Case is a test case for the Test_Bugf function.
type Test_Bugf_Case struct {
	name     string
	format   string
	args     []interface{}
	expected string
}

// Test_Bugf tests the Bugf function.
func Test_Bugf(t *testing.T) {
	t.Parallel()
	for _, tc := range []Test_Bugf_Case{
		{
			name:     "basic case",
			format:   "test %d",
			args:     []interface{}{1},
			expected: "test 1",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			panicked := false
			defer func() {
				if r := recover(); r != tc.expected {
					t.Errorf("expected %s, got %s", tc.expected, r)
				} else {
					panicked = true
				}
			}()
			Bugf(tc.format, tc.args...)
			if !panicked {
				t.Errorf("expected panic")
			}
		})
	}
}

// Test_Bug_Case is a test case for the Test_Bug function.
type Test_Bug_Case struct {
	name     string
	message  string
	expected string
}

// Test_Bug tests the Bug function.
func Test_Bug(t *testing.T) {
	t.Parallel()
	for _, tc := range []Test_Bug_Case{
		{
			name:     "basic case",
			message:  "test",
			expected: "test",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			panicked := false
			defer func() {
				if r := recover(); r != tc.expected {
					t.Errorf("expected %s, got %s", tc.expected, r)
				} else {
					panicked = true
				}
			}()
			Bug(tc.message)
			if !panicked {
				t.Errorf("expected panic")
			}
		})
	}
}

// Test_Bug_Handler tests that Handler returns a function that calls panic.
func Test_Bug_Handler(t *testing.T) {
	t.Parallel()
	panicked := false
	defer func() {
		if r := recover(); r != "test" {
			t.Errorf("expected test, got %s", r)
		} else {
			panicked = true
		}
	}()
	Handler()("test")
	if !panicked {
		t.Errorf("expected panic")
	}
}

// Test_Bug_SetHandler tests that SetHandler sets the handler. It does this
// by replacing the handler with a function that sets a boolean to true, and
// then calling bug.Bug.
func Test_Bug_SetHandler(t *testing.T) {
	// NOT a parallel test because it changes the global state.
	message := ""
	oldHandler := Handler()
	defer SetHandler(oldHandler)
	SetHandler(func(msg string) {
		message = msg
	})
	Bug("test")
	if message != "test" {
		t.Errorf("expected test, got %s", message)
	}
}
