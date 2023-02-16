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
	"time"

	"github.com/neuralnorthwest/mu/bug"
)

// linear is a retry strategy that uses linear backoff.
type linear struct {
	// baseInterval is the baseInterval duration to use for the linear backoff.
	baseInterval time.Duration
	// maxInterval is the maximum duration to wait between retries.
	maxInterval time.Duration
	// increment is the increment to use for the linear backoff.
	increment time.Duration
	// maxAttempts is the maximum number of attempts to make.
	maxAttempts int
}

var _ Strategy = (*linear)(nil)
var _ StrategyWithBaseInterval = (*linear)(nil)
var _ StrategyWithMaxInterval = (*linear)(nil)
var _ StrategyWithIncrement = (*linear)(nil)
var _ StrategyWithMaxAttempts = (*linear)(nil)

// Linear returns a retry strategy that uses linear backoff.
func Linear(opts ...StrategyOption) Strategy {
	l := &linear{
		baseInterval: 100 * time.Millisecond,
		maxInterval:  10 * time.Second,
		increment:    100 * time.Millisecond,
		maxAttempts:  -1,
	}
	for _, opt := range opts {
		opt(l)
	}
	// Validate. If any of these are invalid, we'll report a bug and return
	// a default strategy.
	if l.baseInterval < 0 {
		bug.Bug("baseInterval must be >= 0")
		return Linear()
	}
	if l.increment < 0 {
		bug.Bug("increment must be >= 0")
		return Linear()
	}
	if l.maxInterval >= 0 && l.maxInterval < l.baseInterval {
		bug.Bug("maxInterval must be >= baseInterval")
		return Linear()
	}
	return l
}

// Next implements Strategy.Next.
func (l *linear) Next(err error) time.Duration {
	if l.maxAttempts == 0 {
		return -1
	} else if l.maxAttempts > 0 {
		l.maxAttempts--
	}
	dur := l.baseInterval
	l.baseInterval += l.increment
	if l.maxInterval >= 0 && l.baseInterval > l.maxInterval {
		l.baseInterval = l.maxInterval
	}
	return dur
}

// WithBaseInterval implements StrategyWithBaseInterval.
func (l *linear) WithBaseInterval(d time.Duration) {
	l.baseInterval = d
}

// WithMaxInterval implements StrategyWithMaxInterval.
func (l *linear) WithMaxInterval(d time.Duration) {
	l.maxInterval = d
}

// WithIncrement implements StrategyWithIncrement.
func (l *linear) WithIncrement(d time.Duration) {
	l.increment = d
}

// WithMaxAttempts implements StrategyWithMaxAttempts.
func (l *linear) WithMaxAttempts(n int) {
	l.maxAttempts = n
}
