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

// Package bug provides an interface for handling bugs.
//
// A bug is a condition that should never happen. If a bug is encountered, the
// program should respond in an appropriate way. For example, if a bug is
// encountered in a web server, the server should return a 500 Internal Server
// Error response. If a bug is encountered in a command-line program, the
// program should exit with a non-zero exit code.
//
// To handle a bug, call bug.Bug with a message. To format the message, use
// bug.Bugf.
//
// Example:
//
//	// Override the handler to log the message and exit.
//	bug.SetHandler(func(message string) {
//	  log.Print(message)
//	  os.Exit(1)
//	})
//	// Handle a bug with Bug.
//	bug.Bug("this should never happen")
//	// Handle a bug with Bugf.
//	bug.Bugf("this should never happen: %s", "foo")
//
// The default behavior of bug.Bug is to panic. Applications can and should
// change the behavior of bug.Bug by calling bug.SetHandler. The bug handler
// might be called from multiple goroutines, so it must be thread-safe. Library
// code should not call bug.SetHandler.
//
// When calling bug.Bug from library code, always do so within a defer
// statement (unless it is being called from an infinite loop). This avoids
// making assumptions about the behavior of the handler (in particular, whether
// it returns or not).
package bug
