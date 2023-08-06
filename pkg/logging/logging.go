package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps a zapLogger
type Logger struct {
	*zap.Logger
}

// NewLogger creates a Logger
func NewLogger(name string) (*Logger, error) {
	// Configure
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.NameKey = name

	// Build & return
	zapLogger, err := cfg.Build()
	l := &Logger{
		Logger: zapLogger,
	}
	return l, err
}

// Cleanup flushes the logger's buffer before the program terminates
func (l *Logger) Cleanup() {
	l.Logger.Sync()
}

func (l *Logger) NewChildLogger(name string) *Logger {
	child := l.Logger.Named(name)
	return &Logger {
		Logger: child,
	}
}
