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

// Package service implements the Service type. This type is the core of the
// Mu framework.
//
// Service is based on the concept of hooks and workers. Hooks are functions
// that are invoked at various points in the service lifecycle. Workers are
// goroutines that are started and stopped as part of the service lifecycle.
//
// Service also contains some core functionality that is useful for most
// services. This includes logging, configuration, and signal handling.
package service
