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
