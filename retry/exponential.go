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

// exponential is a retry strategy that uses exponential backoff.
type exponential struct {
	// baseInterval is the baseInterval duration to use for the exponential backoff.
	baseInterval time.Duration
	// maxInterval is the maximum duration to wait between retries.
	maxInterval time.Duration
	// factor is the factor to use for the exponential backoff.
	factor float64
	// maxAttempts is the maximum number of attempts to make.
	maxAttempts int
}

var _ Strategy = (*exponential)(nil)
var _ StrategyWithBaseInterval = (*exponential)(nil)
var _ StrategyWithMaxInterval = (*exponential)(nil)
var _ StrategyWithFactor = (*exponential)(nil)
var _ StrategyWithMaxAttempts = (*exponential)(nil)

// Exponential returns a retry strategy that uses exponential backoff.
func Exponential(opts ...Option) Strategy {
	e := &exponential{
		baseInterval: 100 * time.Millisecond,
		maxInterval:  10 * time.Second,
		factor:       2,
		maxAttempts:  -1,
	}
	for _, opt := range opts {
		opt(e)
	}
	// Validate. If any of these are invalid, we'll report a bug and return
	// a default strategy.
	if e.baseInterval < 0 {
		bug.Bug("baseInterval must be >= 0")
		return Exponential()
	}
	if e.maxInterval < 0 {
		bug.Bug("maxInterval must be >= 0")
		return Exponential()
	}
	if e.factor < 1 {
		bug.Bug("factor must be >= 1")
		return Exponential()
	}
	if e.maxInterval < e.baseInterval {
		bug.Bug("maxInterval must be >= baseInterval")
		return Exponential()
	}
	return e
}

// Next implements Strategy.Next.
func (e *exponential) Next(err error) time.Duration {
	if e.maxAttempts == 0 {
		return -1
	} else if e.maxAttempts > 0 {
		e.maxAttempts--
	}
	dur := e.baseInterval
	e.baseInterval = time.Duration(float64(e.baseInterval) * e.factor)
	if e.baseInterval > e.maxInterval {
		e.baseInterval = e.maxInterval
	}
	return dur
}

// WithBaseInterval implements StrategyWithBaseInterval.
func (e *exponential) WithBaseInterval(d time.Duration) {
	e.baseInterval = d
}

// WithMaxInterval implements StrategyWithMaxInterval.
func (e *exponential) WithMaxInterval(d time.Duration) {
	e.maxInterval = d
}

// WithFactor implements StrategyWithFactor.
func (e *exponential) WithFactor(f float64) {
	e.factor = f
}

// WithMaxAttempts implements StrategyWithMaxAttempts.
func (e *exponential) WithMaxAttempts(n int) {
	e.maxAttempts = n
}
