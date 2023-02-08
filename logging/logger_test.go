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

package logging

import "testing"

// Test_New tests that the New function returns a logger, and that it is a
// zapLogger.
func Test_New(t *testing.T) {
	t.Parallel()
	logger, err := New()
	if err != nil {
		t.Fatalf("New returned an error: %v", err)
	}
	if logger == nil {
		t.Fatal("New returned a nil logger")
	}
	if _, ok := logger.(*zapLogger); !ok {
		t.Fatal("New returned a logger that is not a zapLogger")
	}
	z := logger.(*zapLogger).GetZap()
	if z == nil {
		t.Fatal("zapLogger.GetZap returned a nil logger")
	}
}

// Test_With tests the With function.
func Test_With(t *testing.T) {
	t.Parallel()
	logger, err := New()
	if err != nil {
		t.Fatalf("New returned an error: %v", err)
	}
	logger = logger.With("key", "value")
	if logger == nil {
		t.Fatal("With returned a nil logger")
	}
	if _, ok := logger.(*zapLogger); !ok {
		t.Fatal("With returned a logger that is not a zapLogger")
	}
	z := logger.(*zapLogger).GetZap()
	if z == nil {
		t.Fatal("zapLogger.GetZap returned a nil logger")
	}
}

// Test_Level tests the Level/SetLevel functions.
func Test_Level(t *testing.T) {
	t.Parallel()
	logger, err := New()
	if err != nil {
		t.Fatalf("New returned an error: %v", err)
	}
	if logger.Level() != InfoLevel {
		t.Fatalf("Level returned %v, expected InfoLevel", logger.Level())
	}
	for _, level := range []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel} {
		logger.SetLevel(level)
		if logger.Level() != level {
			t.Fatalf("Level returned %v, expected %v", logger.Level(), level)
		}
	}
}

// Test_WithLevel tests the WithLevel option.
func Test_WithLevel(t *testing.T) {
	t.Parallel()
	logger, err := New(WithLevel(DebugLevel))
	if err != nil {
		t.Fatalf("New returned an error: %v", err)
	}
	if logger.Level() != DebugLevel {
		t.Fatalf("Level returned %v, expected DebugLevel", logger.Level())
	}
}
