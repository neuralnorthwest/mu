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

package main

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/neuralnorthwest/mu/logging"
	mock_logging "github.com/neuralnorthwest/mu/logging/mock"
)

// outputLogCall is a call to output a log message.
type outputLogCall struct {
	// level is the level of the log message.
	level logging.Level
	// message is the log message.
	message string
}

// Test_Adapter_Case is a test case for Test_Adapter.
type Test_Adapter_Case struct {
	// name is the name of the test case.
	name string
	// level is the level to use for the adapter
	level logging.Level
	// inputLogs are the input logs to write to the adapter.
	inputLogs []string
	// expectedCalls are the expected calls to the logger.
	expectedCalls []outputLogCall
}

// Test_Adapter is a test for Adapter.
func Test_Adapter(t *testing.T) {
	t.Parallel()
	// Test cases cover:
	//   - Testing that the logger is called with the indicated AdaptedLevel.
	//   - Testing that the logger is called with the indicated message.
	for _, testCase := range []Test_Adapter_Case{
		{
			name:          "Debug",
			level:         logging.DebugLevel,
			inputLogs:     []string{"debug message"},
			expectedCalls: []outputLogCall{{level: logging.DebugLevel, message: "debug message"}},
		},
		{
			name:          "Info",
			level:         logging.InfoLevel,
			inputLogs:     []string{"info message"},
			expectedCalls: []outputLogCall{{level: logging.InfoLevel, message: "info message"}},
		},
		{
			name:          "Warning",
			level:         logging.WarnLevel,
			inputLogs:     []string{"warning message"},
			expectedCalls: []outputLogCall{{level: logging.WarnLevel, message: "warning message"}},
		},
		{
			name:          "Error",
			level:         logging.ErrorLevel,
			inputLogs:     []string{"error message"},
			expectedCalls: []outputLogCall{{level: logging.ErrorLevel, message: "error message"}},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			logger := mock_logging.NewMockLogger(mockCtrl)
			adapter := logging.NewAdapter(logger, testCase.level)
			for _, expectedCall := range testCase.expectedCalls {
				switch expectedCall.level {
				case logging.DebugLevel:
					logger.EXPECT().Debug(expectedCall.message)
				case logging.InfoLevel:
					logger.EXPECT().Info(expectedCall.message)
				case logging.WarnLevel:
					logger.EXPECT().Warn(expectedCall.message)
				case logging.ErrorLevel:
					logger.EXPECT().Error(expectedCall.message)
				}
			}
			for _, inputLog := range testCase.inputLogs {
				adapter.Print(inputLog)
			}
		})
	}
}
