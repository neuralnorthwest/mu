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

import (
	"testing"

	"github.com/neuralnorthwest/mu/config"
	"github.com/neuralnorthwest/mu/http"
	"github.com/neuralnorthwest/mu/worker"
	"github.com/stretchr/testify/assert"
)

// Test_Cleanups tests that cleanups are run in the reverse order they were
// added.
func Test_Cleanups(t *testing.T) {
	t.Parallel()
	svc, err := New("test-service")
	assert.NoError(t, err)
	var cleanupOrder []int
	svc.Cleanup(func() {
		cleanupOrder = append(cleanupOrder, 1)
	})
	svc.SetupConfig(func(config.Config) error {
		svc.Cleanup(func() {
			cleanupOrder = append(cleanupOrder, 2)
		})
		return nil
	})
	svc.SetupWorkers(func(wg worker.Group) error {
		svc.Cleanup(func() {
			cleanupOrder = append(cleanupOrder, 3)
		})
		return nil
	})
	svc.SetupHTTP(func(*http.Server) error {
		svc.Cleanup(func() {
			cleanupOrder = append(cleanupOrder, 4)
		})
		return nil
	})
	svc.Cleanup(func() {
		cleanupOrder = append(cleanupOrder, 5)
	})
	svc.PreRun(func() error {
		svc.Cancel()
		return nil
	})
	err = svc.Run()
	assert.NoError(t, err)
	assert.Equal(t, []int{4, 3, 2, 5, 1}, cleanupOrder)
}
