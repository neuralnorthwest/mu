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
	level zap.AtomicLevel
}

var _ Logger = (*zapLogger)(nil)

// NewZapLogger creates a new zap logger.
func NewZapLogger(opts ...Option) (Logger, error) {
	zlogger := &zapLogger{
		level: zap.NewAtomicLevelAt(zap.InfoLevel),
	}
	config := zap.NewProductionConfig()
	config.Level = zlogger.level
	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	zlogger.SugaredLogger = logger.Sugar()
	for _, opt := range opts {
		opt(zlogger)
	}
	return zlogger, nil
}

// Level returns the current logging level.
func (l *zapLogger) Level() Level {
	var level Level
	switch l.level.Level() {
	case zap.DebugLevel:
		level = DebugLevel
	case zap.InfoLevel:
		level = InfoLevel
	case zap.WarnLevel:
		level = WarnLevel
	case zap.ErrorLevel:
		level = ErrorLevel
	}
	return level
}

// SetLevel sets the logging level.
func (l *zapLogger) SetLevel(level Level) {
	switch level {
	case DebugLevel:
		l.level.SetLevel(zap.DebugLevel)
	case InfoLevel:
		l.level.SetLevel(zap.InfoLevel)
	case WarnLevel:
		l.level.SetLevel(zap.WarnLevel)
	case ErrorLevel:
		l.level.SetLevel(zap.ErrorLevel)
	}
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
