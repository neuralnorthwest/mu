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

package retry

import (
	"testing"
	"time"
)

// Test_Linear tests the linear retry strategy with the default options.
func Test_Linear(t *testing.T) {
	t.Parallel()
	s := Linear()
	expectedDurations := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		300 * time.Millisecond,
		400 * time.Millisecond,
		500 * time.Millisecond,
		600 * time.Millisecond,
		700 * time.Millisecond,
		800 * time.Millisecond,
		900 * time.Millisecond,
		1000 * time.Millisecond,
	}
	testStrategy(t, s, expectedDurations, 10)
}

// Test_Linear_WithBaseInterval tests the linear retry strategy with a custom
// base interval.
func Test_Linear_WithBaseInterval(t *testing.T) {
	t.Parallel()
	s := Linear(WithBaseInterval(200 * time.Millisecond))
	expectedDurations := []time.Duration{
		200 * time.Millisecond,
		300 * time.Millisecond,
		400 * time.Millisecond,
		500 * time.Millisecond,
		600 * time.Millisecond,
		700 * time.Millisecond,
		800 * time.Millisecond,
		900 * time.Millisecond,
		1000 * time.Millisecond,
		1100 * time.Millisecond,
	}
	testStrategy(t, s, expectedDurations, 10)
}

// Test_Linear_WithMaxInterval tests the linear retry strategy with a custom
// max interval.
func Test_Linear_WithMaxInterval(t *testing.T) {
	t.Parallel()
	s := Linear(WithMaxInterval(500 * time.Millisecond))
	expectedDurations := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		300 * time.Millisecond,
		400 * time.Millisecond,
		500 * time.Millisecond,
		500 * time.Millisecond,
		500 * time.Millisecond,
		500 * time.Millisecond,
		500 * time.Millisecond,
		500 * time.Millisecond,
	}
	testStrategy(t, s, expectedDurations, 10)
}

// Test_Linear_WithIncrement tests the linear retry strategy with a custom
// increment.
func Test_Linear_WithIncrement(t *testing.T) {
	t.Parallel()
	s := Linear(WithIncrement(200 * time.Millisecond))
	expectedDurations := []time.Duration{
		100 * time.Millisecond,
		300 * time.Millisecond,
		500 * time.Millisecond,
		700 * time.Millisecond,
		900 * time.Millisecond,
		1100 * time.Millisecond,
		1300 * time.Millisecond,
		1500 * time.Millisecond,
		1700 * time.Millisecond,
		1900 * time.Millisecond,
	}
	testStrategy(t, s, expectedDurations, 10)
}

// Test_Linear_WithMaxAttempts tests the linear retry strategy with a custom
// max attempts.
func Test_Linear_WithMaxAttempts(t *testing.T) {
	t.Parallel()
	s := Linear(WithMaxAttempts(5))
	expectedDurations := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		300 * time.Millisecond,
		400 * time.Millisecond,
		500 * time.Millisecond,
		-1,
		-1,
		-1,
		-1,
		-1,
	}
	testStrategy(t, s, expectedDurations, 10)
}

// Test_Linear_WithInvalidBaseInterval tests the linear retry strategy with an
// invalid base interval.
func Test_Linear_WithInvalidBaseInterval(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	Linear(WithBaseInterval(-1 * time.Millisecond))
}

// Test_Linear_WithInvalidMaxInterval tests the linear retry strategy with an
// invalid max interval.
func Test_Linear_WithInvalidMaxInterval(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	Linear(WithMaxInterval(-1 * time.Millisecond))
}

// Test_Linear_WithInvalidIncrement tests the linear retry strategy with an
// invalid increment.
func Test_Linear_WithInvalidIncrement(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	Linear(WithIncrement(-1 * time.Millisecond))
}

// Test_Linear_WithInvalidBaseAndMaxInterval tests the linear retry strategy
// with an invalid base and max interval.
func Test_Linear_WithInvalidBaseAndMaxInterval(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	Linear(WithBaseInterval(2*time.Millisecond), WithMaxInterval(1*time.Millisecond))
}
