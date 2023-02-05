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
	"bytes"
	"log"
)

// AdaptedLevel is a log level that can be adapted to the standard library
// logger.
type AdaptedLevel int

const (
	// AdaptedLevelDebug is the debug level.
	AdaptedLevelDebug AdaptedLevel = iota
	// AdaptedLevelInfo is the info level.
	AdaptedLevelInfo
	// AdaptedLevelWarn is the warn level.
	AdaptedLevelWarn
	// AdaptedLevelError is the error level.
	AdaptedLevelError
)

// Adapter adapts the Mu logger to the standard library logger.
type Adapter struct {
	// level is the log level to adapt to.
	level AdaptedLevel
	// logger is the Mu logger.
	logger Logger
	// buffer is a buffer for accumulating log lines.
	buffer bytes.Buffer
}

// NewAdapter creates a new Adapter and returns a log.Logger that writes to it.
// The level controls how the log messages are written to the Mu logger.
func NewAdapter(logger Logger, level AdaptedLevel) *log.Logger {
	return log.New(&Adapter{
		level:  level,
		logger: logger,
	}, "", 0)
}

// Write implements the io.Writer interface.
func (a *Adapter) Write(b []byte) (int, error) {
	a.buffer.Write(b)
	// Pump out all the lines in the buffer
	for {
		line, err := a.buffer.ReadBytes('\n')
		if err != nil {
			break
		}
		// Remove the trailing newline
		line = line[:len(line)-1]
		switch a.level {
		case AdaptedLevelDebug:
			a.logger.Debug(string(line))
		case AdaptedLevelInfo:
			a.logger.Info(string(line))
		case AdaptedLevelWarn:
			a.logger.Warn(string(line))
		case AdaptedLevelError:
			a.logger.Error(string(line))
		}
	}
	return len(b), nil
}
