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
	"time"
)

// Strategy is a retry strategy.
type Strategy interface {
	// Next returns the duration to wait before the next retry. It is passed
	// the error returned by the last attempt. It returns a negative duration
	// to stop retrying.
	Next(err error) time.Duration
}

// Do retries the provided function according to the provided strategy, until
// the context is canceled or the function returns nil.
func Do(ctx context.Context, strategy Strategy, fn func() error) error {
	for {
		err := fn()
		if err == nil {
			return nil
		}
		dur := strategy.Next(err)
		if dur < 0 {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
		}
	}
}

// DoAsync retries the provided function according to the provided strategy,
// until the context is canceled or the function returns nil. It returns a
// channel that will receive the error returned by the function, or nil if the
// function succeeded.
func DoAsync(ctx context.Context, strategy Strategy, fn func() error) <-chan error {
	ch := make(chan error, 1)
	go func() {
		ch <- Do(ctx, strategy, fn)
	}()
	return ch
}
