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

package metrics

import (
	ht "net/http"

	"github.com/neuralnorthwest/mu/logging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler returns an ht.Handler that serves the metrics. The provided logger
// will be used to log any errors that occur while serving the metrics.
func (m *metrics) Handler(logger logging.Logger) ht.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{
		ErrorLog:      logging.NewAdapter(logger, logging.AdaptedLevelError),
		ErrorHandling: promhttp.ContinueOnError,
	})
}