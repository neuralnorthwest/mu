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

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package mock -destination mock/logger.go github.com/neuralnorthwest/mu/logging Logger

// Logger is the interface for logging.
type Logger interface {
	// DPanic logs a panic message and panics in development mode.
	DPanic(args ...interface{})
	// DPanicf logs a panic message with formatting and panics in development mode.
	DPanicf(template string, args ...interface{})
	// DPanicln logs a panic message with a newline and panics in development mode.
	DPanicln(args ...interface{})
	// DPanicw logs a panic message with key-value pairs and panics in development mode.
	DPanicw(msg string, keysAndValues ...interface{})
	// Debug logs a debug message.
	Debug(args ...interface{})
	// Debugf logs a debug message with formatting.
	Debugf(template string, args ...interface{})
	// Debugln logs a debug message with a newline.
	Debugln(args ...interface{})
	// Debugw logs a debug message with key-value pairs.
	Debugw(msg string, keysAndValues ...interface{})
	// Error logs an error message.
	Error(args ...interface{})
	// Errorf logs an error message with formatting.
	Errorf(template string, args ...interface{})
	// Errorln logs an error message with a newline.
	Errorln(args ...interface{})
	// Errorw logs an error message with key-value pairs.
	Errorw(msg string, keysAndValues ...interface{})
	// Fatal logs a fatal message.
	Fatal(args ...interface{})
	// Fatalf logs a fatal message with formatting.
	Fatalf(template string, args ...interface{})
	// Fatalln logs a fatal message with a newline.
	Fatalln(args ...interface{})
	// Fatalw logs a fatal message with key-value pairs.
	Fatalw(msg string, keysAndValues ...interface{})
	// Info logs an info message.
	Info(args ...interface{})
	// Infof logs an info message with formatting.
	Infof(template string, args ...interface{})
	// Infoln logs an info message with a newline.
	Infoln(args ...interface{})
	// Infow logs an info message with key-value pairs.
	Infow(msg string, keysAndValues ...interface{})
	// Panic logs a panic message.
	Panic(args ...interface{})
	// Panicf logs a panic message with formatting.
	Panicf(template string, args ...interface{})
	// Panicln logs a panic message with a newline.
	Panicln(args ...interface{})
	// Panicw logs a panic message with key-value pairs.
	Panicw(msg string, keysAndValues ...interface{})
	// Sync flushes the logger.
	Sync() error
	// Warn logs a warning message.
	Warn(args ...interface{})
	// Warnf logs a warning message with formatting.
	Warnf(template string, args ...interface{})
	// Warnln logs a warning message with a newline.
	Warnln(args ...interface{})
	// Warnw logs a warning message with key-value pairs.
	Warnw(msg string, keysAndValues ...interface{})
	// With creates a child logger and adds structured context to it.
	With(args ...interface{}) Logger
}

// New returns a new logger based on zap.
func New() (Logger, error) {
	return NewZapLogger()
}
