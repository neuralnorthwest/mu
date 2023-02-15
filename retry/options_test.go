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

// noOptionsStrategy is a Strategy that doesn't support any of the options.
type noOptionsStrategy struct {
}

// Next implements Strategy.Next.
func (s *noOptionsStrategy) Next(err error) time.Duration {
	return 0
}

// Test_WithBaseInterval_NotSupported tests that WithBaseInterval bugs when
// called on a strategy that doesn't support it.
func Test_WithBaseInterval_NotSupported(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	s := &noOptionsStrategy{}
	WithBaseInterval(100 * time.Millisecond)(s)
}

// Test_WithMaxInterval_NotSupported tests that WithMaxInterval bugs when
// called on a strategy that doesn't support it.
func Test_WithMaxInterval_NotSupported(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	s := &noOptionsStrategy{}
	WithMaxInterval(100 * time.Millisecond)(s)
}

// Test_WithIncrement_NotSupported tests that WithIncrement bugs when called
// on a strategy that doesn't support it.
func Test_WithIncrement_NotSupported(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	s := &noOptionsStrategy{}
	WithIncrement(100 * time.Millisecond)(s)
}

// Test_WithFactor_NotSupported tests that WithFactor bugs when called on a
// strategy that doesn't support it.
func Test_WithFactor_NotSupported(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	s := &noOptionsStrategy{}
	WithFactor(2)(s)
}

// Test_WithMaxAttempts_NotSupported tests that WithMaxAttempts bugs when
// called on a strategy that doesn't support it.
func Test_WithMaxAttempts_NotSupported(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	s := &noOptionsStrategy{}
	WithMaxAttempts(5)(s)
}
