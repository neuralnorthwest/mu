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

	"github.com/neuralnorthwest/mu/status"
)

// testStrategy tests the given strategy with the given number of attempts.
func testStrategy(t *testing.T, s Strategy, expectedDurations []time.Duration, attempts int) {
	for i := 0; i < attempts; i++ {
		d := s.Next(status.ErrInvalidArgument)
		if d != expectedDurations[i] {
			t.Errorf("expected duration %v, got %v", expectedDurations[i], d)
		}
	}
}

// Test_Exponential tests the exponential retry strategy with the default
// options.
func Test_Exponential(t *testing.T) {
	t.Parallel()
	s := Exponential()
	expectedDurations := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		800 * time.Millisecond,
		1600 * time.Millisecond,
		3200 * time.Millisecond,
		6400 * time.Millisecond,
		10000 * time.Millisecond,
		10000 * time.Millisecond,
		10000 * time.Millisecond,
	}
	testStrategy(t, s, expectedDurations, 10)
}

// Test_Exponential_WithBaseInterval tests the exponential retry strategy with
// a custom base interval.
func Test_Exponential_WithBaseInterval(t *testing.T) {
	t.Parallel()
	s := Exponential(WithBaseInterval(200 * time.Millisecond))
	expectedDurations := []time.Duration{
		200 * time.Millisecond,
		400 * time.Millisecond,
		800 * time.Millisecond,
		1600 * time.Millisecond,
		3200 * time.Millisecond,
		6400 * time.Millisecond,
		10000 * time.Millisecond,
		10000 * time.Millisecond,
		10000 * time.Millisecond,
		10000 * time.Millisecond,
	}
	testStrategy(t, s, expectedDurations, 10)
}

// Test_Exponential_WithMaxInterval tests the exponential retry strategy with
// a custom max interval.
func Test_Exponential_WithMaxInterval(t *testing.T) {
	t.Parallel()
	s := Exponential(WithMaxInterval(500 * time.Millisecond))
	expectedDurations := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		500 * time.Millisecond,
		500 * time.Millisecond,
		500 * time.Millisecond,
		500 * time.Millisecond,
		500 * time.Millisecond,
		500 * time.Millisecond,
		500 * time.Millisecond,
	}
	testStrategy(t, s, expectedDurations, 10)
}

// Test_Exponential_NoMaxInterval tests the exponential retry strategy with
// no max interval.
func Test_Exponential_NoMaxInterval(t *testing.T) {
	t.Parallel()
	s := Exponential(WithNoMaxInterval())
	expectedDurations := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		800 * time.Millisecond,
		1600 * time.Millisecond,
		3200 * time.Millisecond,
		6400 * time.Millisecond,
		12800 * time.Millisecond,
		25600 * time.Millisecond,
		51200 * time.Millisecond,
	}
	testStrategy(t, s, expectedDurations, 10)
}

// Test_Exponential_WithFactor tests the exponential retry strategy with a
// custom factor.
func Test_Exponential_WithFactor(t *testing.T) {
	t.Parallel()
	s := Exponential(WithFactor(3))
	expectedDurations := []time.Duration{
		100 * time.Millisecond,
		300 * time.Millisecond,
		900 * time.Millisecond,
		2700 * time.Millisecond,
		8100 * time.Millisecond,
		10000 * time.Millisecond,
		10000 * time.Millisecond,
		10000 * time.Millisecond,
		10000 * time.Millisecond,
		10000 * time.Millisecond,
	}
	testStrategy(t, s, expectedDurations, 10)
}

// Test_Exponential_WithMaxAttempts tests the exponential retry strategy with
// a custom max attempts.
func Test_Exponential_WithMaxAttempts(t *testing.T) {
	t.Parallel()
	s := Exponential(WithMaxAttempts(5))
	expectedDurations := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		800 * time.Millisecond,
		1600 * time.Millisecond,
		-1,
		-1,
		-1,
		-1,
		-1,
	}
	testStrategy(t, s, expectedDurations, 10)
}

// Test_Exponential_WithInvalidBaseInterval tests the exponential retry
// strategy with an invalid base interval.
func Test_Exponential_WithInvalidBaseInterval(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic")
		}
	}()
	Exponential(WithBaseInterval(-1 * time.Millisecond))
}

// Test_Exponential_WithInvalidFactor tests the exponential retry strategy
// with an invalid factor.
func Test_Exponential_WithInvalidFactor(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic")
		}
	}()
	Exponential(WithFactor(-1))
}

// Test_Exponential_WithInvalidBaseAndMaxInterval tests the exponential retry
// strategy with an invalid base and max interval.
func Test_Exponential_WithInvalidBaseAndMaxInterval(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic")
		}
	}()
	Exponential(WithBaseInterval(1000*time.Millisecond), WithMaxInterval(100*time.Millisecond))
}
