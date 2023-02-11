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

import (
	"github.com/neuralnorthwest/mu/http"
	"github.com/neuralnorthwest/mu/worker"
)

// setupWorkers sets up the service workers. This is where we start a secondary
// HTTP server that serves metrics.
func (s *complete) setupWorkers(workerGroup worker.Group) error {
	// Create a diagnostics HTTP server.
	diagnosticsServer, err := http.NewServer(http.WithAddress(":8081"))
	if err != nil {
		return err
	}
	s.diagnosticsServer = diagnosticsServer
	return workerGroup.Add("diagnostics_server", diagnosticsServer)
}
