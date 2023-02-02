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

import (
	"go.uber.org/zap"
)

// zapLogger is a logger that uses the zap library.
type zapLogger struct {
	*zap.SugaredLogger
}

var _ Logger = (*zapLogger)(nil)

// NewZapLogger creates a new zap logger.
func NewZapLogger() (Logger, error) {
	logger, err := zap.NewProduction(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &zapLogger{
		SugaredLogger: logger.Sugar(),
	}, nil
}

// GetZap returns the underlying zap logger.
func (l *zapLogger) GetZap() *zap.SugaredLogger {
	return l.SugaredLogger
}

// With returns a new logger with the given key-value pairs added to its context.
func (l *zapLogger) With(args ...interface{}) Logger {
	return &zapLogger{
		SugaredLogger: l.SugaredLogger.With(args...),
	}
}
