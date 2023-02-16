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

package complete

import "github.com/neuralnorthwest/mu/metrics"

// met holds the various metrics that we want to expose.
type met struct {
	// metrics is the mu metrics object.
	metrics metrics.Metrics
	// helloCounter is a counter metric that counts the number of times we've
	// said hello.
	helloCounter metrics.Counter
}

// setup sets up the service metrics.
func (s *complete) setupMetrics() error {
	m, err := metrics.New()
	if err != nil {
		return err
	}
	s.metrics.metrics = m
	s.metrics.helloCounter = m.NewCounter("hello_counter", "The number of times we've said hello.")
	return nil
}
