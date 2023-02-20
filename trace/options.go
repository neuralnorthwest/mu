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

package trace

// options is an interface to apply options to a tracer.
type options interface {
	collectorEndpoint(endpoint string)
	insecure()
}

// baseOptions is a group of options for a tracer.
type baseOptions struct{}

// collectorEndpoint sets the collector endpoint.
func (o baseOptions) collectorEndpoint(endpoint string) {}

// insecure sets the insecure flag.
func (o baseOptions) insecure() {}

// TracerOption applies an option to an options.
type TracerOption interface {
	apply(options)
}

// TracerOptionFunc is a function that applies an option to an options.
type TracerOptionFunc func(options)

// apply applies the option to the options.
func (f TracerOptionFunc) apply(o options) {
	f(o)
}

// WithCollectorEndpoint returns an option that sets the collector endpoint.
func WithCollectorEndpoint(endpoint string) TracerOption {
	return TracerOptionFunc(func(o options) {
		o.collectorEndpoint(endpoint)
	})
}

// WithInsecure returns an option that sets the insecure flag.
func WithInsecure() TracerOption {
	return TracerOptionFunc(func(o options) {
		o.insecure()
	})
}
