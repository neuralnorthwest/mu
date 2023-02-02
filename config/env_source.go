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

package config

import (
	"os"
	"strconv"

	"github.com/neuralnorthwest/mu/status"
)

// envSource is a source that reads from the environment.
type envSource struct {
	// prefix is the prefix for the environment variables.
	prefix string
}

// SetPrefix sets the prefix for the environment variables.
func (s *envSource) SetPrefix(prefix string) {
	s.prefix = prefix
}

// LoadInt loads the value of the int variable with the given name.
func (s *envSource) LoadInt(name string) (int, error) {
	str, ok := os.LookupEnv(s.prefix + name)
	if !ok {
		return 0, status.ErrNotFound
	}
	return strconv.Atoi(str)
}

// LoadString loads the value of the string variable with the given name.
func (s *envSource) LoadString(name string) (string, error) {
	str, ok := os.LookupEnv(s.prefix + name)
	if !ok {
		return "", status.ErrNotFound
	}
	return str, nil
}

// LoadBool loads the value of the bool variable with the given name.
func (s *envSource) LoadBool(name string) (bool, error) {
	str, ok := os.LookupEnv(s.prefix + name)
	if !ok {
		return false, status.ErrNotFound
	}
	return strconv.ParseBool(str)
}
