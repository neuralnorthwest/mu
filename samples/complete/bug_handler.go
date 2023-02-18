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

import "runtime/debug"

// bugHandler is the complete service bug handler.
func (s *complete) bugHandler(msg string) {
	stack := string(debug.Stack())
	s.Logger().Errorw("BUG", "message", msg, "stack", stack)
	// If development mode is enabled, we also log the bug to the console
	// and panic.
	if s.config != nil && s.config.DevMode() {
		s.Logger().Errorw("BUG", "message", msg, "stack", stack)
		panic(msg)
	}
}
