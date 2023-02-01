package logging

import (
	"go.uber.org/zap"
)

// zapLogger is a logger that uses the zap library.
type zapLogger struct {
	logger *zap.SugaredLogger
}

var _ Logger = (*zapLogger)(nil)

// NewZapLogger creates a new zap logger.
func NewZapLogger() (Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &zapLogger{
		logger: logger.Sugar(),
	}, nil
}

// GetZap returns the underlying zap logger.
func (l *zapLogger) GetZap() *zap.SugaredLogger {
	return l.logger
}

// DPanic logs a panic message and panics in development mode.
func (l *zapLogger) DPanic(args ...interface{}) {
	l.logger.DPanic(args...)
}

// DPanicf logs a panic message with formatting and panics in development mode.
func (l *zapLogger) DPanicf(template string, args ...interface{}) {
	l.logger.DPanicf(template, args...)
}

// DPanicln logs a panic message with a newline and panics in development mode.
func (l *zapLogger) DPanicln(args ...interface{}) {
	l.logger.DPanicln(args...)
}

// DPanicw logs a panic message with key-value pairs and panics in development mode.
func (l *zapLogger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.logger.DPanicw(msg, keysAndValues...)
}

// Debug logs a debug message.
func (l *zapLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Debugf logs a debug message with formatting.
func (l *zapLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

// Debugln logs a debug message with a newline.
func (l *zapLogger) Debugln(args ...interface{}) {
	l.logger.Debugln(args...)
}

// Debugw logs a debug message with key-value pairs.
func (l *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

// Error logs an error message.
func (l *zapLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Errorf logs an error message with formatting.
func (l *zapLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

// Errorln logs an error message with a newline.
func (l *zapLogger) Errorln(args ...interface{}) {
	l.logger.Errorln(args...)
}

// Errorw logs an error message with key-value pairs.
func (l *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

// Fatal logs a fatal message and exits.
func (l *zapLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf logs a fatal message with formatting and exits.
func (l *zapLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

// Fatalln logs a fatal message with a newline and exits.
func (l *zapLogger) Fatalln(args ...interface{}) {
	l.logger.Fatalln(args...)
}

// Fatalw logs a fatal message with key-value pairs and exits.
func (l *zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

// Info logs an info message.
func (l *zapLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof logs an info message with formatting.
func (l *zapLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

// Infoln logs an info message with a newline.
func (l *zapLogger) Infoln(args ...interface{}) {
	l.logger.Infoln(args...)
}

// Infow logs an info message with key-value pairs.
func (l *zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

// Panic logs a panic message and panics.
func (l *zapLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

// Panicf logs a panic message with formatting and panics.
func (l *zapLogger) Panicf(template string, args ...interface{}) {
	l.logger.Panicf(template, args...)
}

// Panicln logs a panic message with a newline and panics.
func (l *zapLogger) Panicln(args ...interface{}) {
	l.logger.Panicln(args...)
}

// Panicw logs a panic message with key-value pairs and panics.
func (l *zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.logger.Panicw(msg, keysAndValues...)
}

// Sync flushes the logger.
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}

// Warn logs a warning message.
func (l *zapLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Warnf logs a warning message with formatting.
func (l *zapLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

// Warnln logs a warning message with a newline.
func (l *zapLogger) Warnln(args ...interface{}) {
	l.logger.Warnln(args...)
}

// Warnw logs a warning message with key-value pairs.
func (l *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

// With creates a child logger and adds structured context to it.
func (l *zapLogger) With(keysAndValues ...interface{}) Logger {
	return &zapLogger{
		logger: l.logger.With(keysAndValues...),
	}
}
