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
	"context"
	"errors"
	"testing"
	"time"
)

// checkError checks that err is equal to expectedErr.
func checkError(t *testing.T, err, expectedErr error) {
	t.Helper()
	if err == nil && expectedErr == nil {
		return
	}
	if err == nil || expectedErr == nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err.Error() != expectedErr.Error() {
		t.Fatalf("unexpected error: %v", err)
	}
}

// failCall is a duration to wait and an error to return.
type failCall struct {
	wait time.Duration
	err  error
}

// failAfter describes a sequence of calls to fail.
type failAfter struct {
	calls []failCall
	i     int
}

// F is the function that fails after a number of attempts.
func (f *failAfter) F(t *testing.T) error {
	if f.i >= len(f.calls) {
		t.Fatalf("unexpected call to F")
	}
	call := f.calls[f.i]
	f.i++
	time.Sleep(call.wait)
	return call.err
}

// Test_Do_Case is a test case for Test_Do.
type Test_Do_Case struct {
	name              string
	construct         func(opts ...Option) Strategy
	opts              []Option
	fail              failAfter
	expectedIntervals []time.Duration
	expectedErr       error
	cancelAfter       time.Duration
}

var doCases = []Test_Do_Case{
	{
		name:      "linear, default, succeed after 5",
		construct: Linear,
		fail: failAfter{
			calls: []failCall{
				{wait: 100 * time.Millisecond, err: errors.New("1")},
				{wait: 100 * time.Millisecond, err: errors.New("2")},
				{wait: 100 * time.Millisecond, err: errors.New("3")},
				{wait: 100 * time.Millisecond, err: errors.New("4")},
				{wait: 100 * time.Millisecond, err: nil},
			},
		},
		expectedIntervals: []time.Duration{
			0 * time.Millisecond,
			200 * time.Millisecond,
			500 * time.Millisecond,
			900 * time.Millisecond,
			1400 * time.Millisecond,
		},
	},
	{
		name:      "linear, default, fail after 5",
		construct: Linear,
		opts: []Option{
			WithMaxAttempts(4),
		},
		fail: failAfter{
			calls: []failCall{
				{wait: 100 * time.Millisecond, err: errors.New("1")},
				{wait: 100 * time.Millisecond, err: errors.New("2")},
				{wait: 100 * time.Millisecond, err: errors.New("3")},
				{wait: 100 * time.Millisecond, err: errors.New("4")},
				{wait: 100 * time.Millisecond, err: errors.New("5")},
			},
		},
		expectedIntervals: []time.Duration{
			0,
			200 * time.Millisecond,
			500 * time.Millisecond,
			900 * time.Millisecond,
			1400 * time.Millisecond,
		},
		expectedErr: errors.New("5"),
	},
	{
		name:      "linear, default, cancel after 3",
		construct: Linear,
		fail: failAfter{
			calls: []failCall{
				{wait: 100 * time.Millisecond, err: errors.New("1")},
				{wait: 100 * time.Millisecond, err: errors.New("2")},
				{wait: 100 * time.Millisecond, err: errors.New("3")},
				{wait: 100 * time.Millisecond, err: errors.New("4")},
				{wait: 100 * time.Millisecond, err: errors.New("5")},
			},
		},
		expectedIntervals: []time.Duration{
			0,
			200 * time.Millisecond,
			500 * time.Millisecond,
		},
		cancelAfter: 600 * time.Millisecond,
		expectedErr: context.DeadlineExceeded,
	},
}

// Test_Do tests that Do executes retries as expected.
func Test_Do(t *testing.T) {
	t.Parallel()
	// increase tolerance if test fails
	const tolerance = 10 * time.Millisecond
	for _, c := range doCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			start := time.Now()
			var intervals []time.Duration
			if c.cancelAfter > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, c.cancelAfter)
				defer cancel()
			}
			err := Do(ctx, c.construct(c.opts...), func() error {
				intervals = append(intervals, time.Since(start))
				return c.fail.F(t)
			})
			checkError(t, err, c.expectedErr)
			if len(intervals) != len(c.expectedIntervals) {
				t.Fatalf("unexpected number of intervals: %d", len(intervals))
			}
			for i := range intervals {
				deviation := intervals[i] - c.expectedIntervals[i]
				if deviation < 0 || deviation > tolerance {
					t.Fatalf("unexpected interval: i=%d, %v (deviation: %v)", i, intervals[i], deviation)
				}
			}
		})
	}
}

// Test_DoAsync tests that DoAsync executes works as expected.
func Test_DoAsync(t *testing.T) {
	t.Parallel()
	// increase tolerance if test fails
	const tolerance = 10 * time.Millisecond
	for _, c := range doCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			start := time.Now()
			var intervals []time.Duration
			if c.cancelAfter > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, c.cancelAfter)
				defer cancel()
			}
			errCh := DoAsync(ctx, c.construct(c.opts...), func() error {
				intervals = append(intervals, time.Since(start))
				return c.fail.F(t)
			})
			err := <-errCh
			checkError(t, err, c.expectedErr)
			if len(intervals) != len(c.expectedIntervals) {
				t.Fatalf("unexpected number of intervals: %d", len(intervals))
			}
			for i := range intervals {
				deviation := intervals[i] - c.expectedIntervals[i]
				if deviation < 0 || deviation > tolerance {
					t.Fatalf("unexpected interval: i=%d, %v (deviation: %v)", i, intervals[i], deviation)
				}
			}
		})
	}
}
