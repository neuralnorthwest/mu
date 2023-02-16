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
	"fmt"
	"time"

	"github.com/neuralnorthwest/mu/bug"
)

// StrategyOption is an option for a Strategy.
type StrategyOption func(Strategy)

// StrategyWithBaseInterval is a Strategy that supports a base duration.
type StrategyWithBaseInterval interface {
	Strategy
	WithBaseInterval(d time.Duration)
}

// StrategyWithMaxInterval is a Strategy that supports a maximum duration.
type StrategyWithMaxInterval interface {
	Strategy
	WithMaxInterval(d time.Duration)
}

// StrategyWithFactor is a Strategy that supports a factor.
type StrategyWithFactor interface {
	Strategy
	WithFactor(f float64)
}

// StrategyWithMaxAttempts is a Strategy that supports a maximum number of
// attempts.
type StrategyWithMaxAttempts interface {
	Strategy
	WithMaxAttempts(n int)
}

// StrategyWithIncrement is a Strategy that supports an increment.
type StrategyWithIncrement interface {
	Strategy
	WithIncrement(d time.Duration)
}

// WithBaseInterval sets the base interval for the Strategy. If this option is
// not set, the Strategy will use a default base interval. This option is
// supported by the following Strategy types:
//   - Exponential
func WithBaseInterval(d time.Duration) StrategyOption {
	return func(s Strategy) {
		if sb, ok := s.(StrategyWithBaseInterval); ok {
			sb.WithBaseInterval(d)
		} else {
			bug.Bug(fmt.Sprintf("Strategy %T does not support WithBaseInterval", s))
		}
	}
}

// WithMaxInterval sets the maximum interval for the Strategy. If this option is
// not set, the Strategy may increase the interval indefinitely. This option is
// supported by the following Strategy types:
//   - Exponential
func WithMaxInterval(d time.Duration) StrategyOption {
	return func(s Strategy) {
		if si, ok := s.(StrategyWithMaxInterval); ok {
			si.WithMaxInterval(d)
		} else {
			bug.Bug(fmt.Sprintf("Strategy %T does not support WithMaxInterval", s))
		}
	}
}

// WithFactor sets the factor for the Strategy. If this option is not set, the
// Strategy will use a default factor. This option is supported by the following
// Strategy types:
//   - Exponential
func WithFactor(f float64) StrategyOption {
	return func(s Strategy) {
		if sf, ok := s.(StrategyWithFactor); ok {
			sf.WithFactor(f)
		} else {
			bug.Bug(fmt.Sprintf("Strategy %T does not support WithFactor", s))
		}
	}
}

// WithMaxAttempts sets the maximum number of attempts for the Strategy. If
// this option is not set, the Strategy will retry indefinitely. This option is
// supported by the following Strategy types:
//   - Exponential
func WithMaxAttempts(n int) StrategyOption {
	return func(s Strategy) {
		if sa, ok := s.(StrategyWithMaxAttempts); ok {
			sa.WithMaxAttempts(n)
		} else {
			bug.Bug(fmt.Sprintf("Strategy %T does not support WithMaxAttempts", s))
		}
	}
}

// WithIncrement sets the increment for the Strategy. If this option is not
// set, the Strategy will use a default increment. This option is supported by
// the following Strategy types:
//   - Linear
func WithIncrement(d time.Duration) StrategyOption {
	return func(s Strategy) {
		if si, ok := s.(StrategyWithIncrement); ok {
			si.WithIncrement(d)
		} else {
			bug.Bug(fmt.Sprintf("Strategy %T does not support WithIncrement", s))
		}
	}
}

// doOptions is the set of options for Do.
type doOptions struct {
	// onRetryAttempt is a function to call before each retry attempt. If
	// onRetryAttempt returns an error, the retry will be aborted and the error
	// will be returned. attempt is the number of the retry attempt, starting
	// with 1. err is the error returned by the previous attempt.
	onRetryAttempt func(attempt int, err error) error
	// onRetry is a function to call after each retry attempt. If onRetry
	// returns an error, the retry will be aborted and the error will be
	// returned. err is the error returned by the previous attempt.
	onRetry func(err error) error
}

// DoOption is an option for Do.
type DoOption func(*doOptions)

// OnRetryAttempt sets a function to call before each retry attempt. If
// onRetryAttempt returns an error, the retry will be aborted and the error
// will be returned. attempt is the number of the retry attempt, starting with
// 1. err is the error returned by the previous attempt. If this option is
// given multiple times, only the last one will be used. If this option is used
// together with OnRetry, OnRetryAttempt will be called first.
func OnRetryAttempt(f func(attempt int, err error) error) DoOption {
	return func(o *doOptions) {
		o.onRetryAttempt = f
	}
}

// OnRetry sets a function to call after each retry attempt. If onRetry
// returns an error, the retry will be aborted and the error will be returned.
// err is the error returned by the previous attempt. If this option is given
// multiple times, only the last one will be used. If this option is used
// together with OnRetryAttempt, OnRetryAttempt will be called first.
//
// This option is useful for identifying non-retryable errors.
func OnRetry(f func(err error) error) DoOption {
	return func(o *doOptions) {
		o.onRetry = f
	}
}
