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
	"errors"
	"fmt"
	"testing"

	"github.com/neuralnorthwest/mu/status"
)

type register func(name string, defaultValue interface{}) error

// Test_Names tests that an error is returned whenever a variable is registered
// with a name that has already been registered.
func Test_Names(t *testing.T) {
	c := New()
	defaultvals := []interface{}{"hello", 1, true}
	regfuncs := []register{
		func(name string, defaultValue interface{}) error {
			return c.NewString(name, defaultvals[0].(string), "")
		},
		func(name string, defaultValue interface{}) error {
			return c.NewInt(name, defaultvals[1].(int), "")
		},
		func(name string, defaultValue interface{}) error {
			return c.NewBool(name, defaultvals[2].(bool), "")
		},
	}
	for i, regfunc := range regfuncs[:len(regfuncs)-1] {
		for j, regfunc2 := range regfuncs[i+1:] {
			name := fmt.Sprintf("test%d", i*len(regfuncs)+j)
			// Test regfunc followed by regfunc2.
			err := regfunc(name+"_a", defaultvals[i])
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			err = regfunc2(name+"_a", defaultvals[j])
			if err == nil {
				t.Errorf("expected error")
			}
			if err != nil && !errors.Is(err, status.ErrAlreadyExists) {
				t.Errorf("unexpected error: %v", err)
			}
			// Test regfunc2 followed by regfunc.
			err = regfunc2(name+"_b", defaultvals[j])
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			err = regfunc(name+"_b", defaultvals[i])
			if err == nil {
				t.Errorf("expected error")
			}
			if err != nil && !errors.Is(err, status.ErrAlreadyExists) {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}
