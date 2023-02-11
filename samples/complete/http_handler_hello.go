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

import ht "net/http"

// helloHandler is a HTTP handler that returns the configured message.
func (s *complete) helloHandler(w ht.ResponseWriter, r *ht.Request) {
	// Increment the hello counter.
	s.metrics.helloCounter.Inc()
	// Get the configured message.
	msg := s.config.Message()
	// Write the message to the response.
	_, err := w.Write([]byte(msg))
	if err != nil {
		s.Logger().Errorw("unable to write response", "err", err)
	}
}
