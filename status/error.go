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

package status

type Error string

func (e Error) Error() string {
	return string(e)
}

// ErrInvalidVersion is returned when an invalid version is provided.
var ErrInvalidVersion = Error("invalid version")

// ErrInvalidRange is returned when an invalid range is provided.
var ErrInvalidRange = Error("invalid range")

// ErrOutOfRange is returned when a value is out of range.
var ErrOutOfRange = Error("out of range")

// ErrAlreadyExists is returned when a value already exists.
var ErrAlreadyExists = Error("already exists")

// ErrNotFound is returned when a value is not found.
var ErrNotFound = Error("not found")

// ErrInvalidArgument is returned when an invalid argument is provided.
var ErrInvalidArgument = Error("invalid argument")

// ErrAlreadyStarted is returned when a service is already started.
var ErrAlreadyStarted = Error("already started")

// ErrNotStarted is returned when a service is not started.
var ErrNotStarted = Error("not started")

// ErrNotImplemented is returned when a method is not implemented.
var ErrNotImplemented = Error("not implemented")

// ErrServerError is returned when a server error occurs.
var ErrServerError = Error("server error")

// ErrClientError is returned when a client error occurs.
var ErrClientError = Error("client error")
