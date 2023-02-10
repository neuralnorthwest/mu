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

package service

// CleanupFunc is a function that cleans up a service.
type CleanupFunc func()

// Cleanup registers a cleanup function that is invoked when the service is
// stopped. Multiple cleanup functions can be registered. Cleanups are invoked
// in the reverse order they are registered.
func (s *Service) Cleanup(f CleanupFunc) {
	s.cleanups = append(s.cleanups, f)
}

// invokeCleanups invokes the cleanups for the service.
func (s *Service) invokeCleanups() {
	for i := len(s.cleanups) - 1; i >= 0; i-- {
		s.cleanups[i]()
	}
}
