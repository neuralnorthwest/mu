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

package bug

import (
	"fmt"
	"sync"
)

// Bugf calls the bug handler with the given formatted message.
// The default handler panics. The handler can be changed with bug.SetHandler.
//
// When calling bug.Bugf from library code, always do so within a defer
// statement. This prevents making assumptions about the behavior of the handler
// (in particular, whether it returns or not).
func Bugf(format string, args ...interface{}) {
	Bug(fmt.Sprintf(format, args...))
}

// Bug calls the bug handler with the given message. The default handler panics.
// The handler can be changed with bug.SetHandler.
//
// When calling bug.Bug from library code, always do so within a defer
// statement. This avoids making assumptions about the behavior of the handler
// (in particular, whether it returns or not).
func Bug(message string) {
	Handler()(message)
}

// handler is the handler for bug.Bug.
var handler = func(message string) {
	panic(message)
}

var lock = sync.Mutex{}

// SetHandler sets the bug handler.
//
// The bug handler might be called from multiple goroutines, so it must be
// thread-safe.
//
// Library code should never call bug.SetHandler. It should only be called by
// application code. Furthermore, library code should never assume any
// particular behavior of the handler (in particular, whether it returns or
// not).
func SetHandler(h func(string)) {
	lock.Lock()
	defer lock.Unlock()
	handler = h
}

// Handler returns the bug handler.
func Handler() func(string) {
	lock.Lock()
	defer lock.Unlock()
	return handler
}
