// Package logger provides structured logging for Wayframe applications.
// It wraps Go's standard slog package with a simplified interface.
package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

// Level represents the severity of a log message.
type Level int

const (
	// DebugLevel logs are typically voluminous and are usually disabled in production.
	DebugLevel Level = iota
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual human review.
	WarnLevel
	// ErrorLevel logs are high-priority and should be addressed.
	ErrorLevel
)

// Logger provides structured logging capabilities using slog.
type Logger struct {
	logger *slog.Logger
}

// New creates a new Logger with the specified minimum level using slog.
// Logs with a level lower than the minimum will be discarded.
func New(level Level) *Logger {
	slogLevel := levelToSlogLevel(level)
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
	})
	return &Logger{
		logger: slog.New(handler),
	}
}

// NewWithHandler creates a new Logger with a custom slog.Handler.
func NewWithHandler(handler slog.Handler) *Logger {
	return &Logger{
		logger: slog.New(handler),
	}
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(w io.Writer) {
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	l.logger = slog.New(handler)
}

// WithField creates a new logger with an additional contextual field.
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		logger: l.logger.With(key, value),
	}
}

// WithFields creates a new logger with multiple contextual fields.
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		logger: l.logger.With(args...),
	}
}

// Debug logs a message at DebugLevel.
func (l *Logger) Debug(msg string) {
	l.logger.Debug(msg)
}

// Debugf logs a formatted message at DebugLevel.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debug(sprintf(format, args...))
}

// Info logs a message at InfoLevel.
func (l *Logger) Info(msg string) {
	l.logger.Info(msg)
}

// Infof logs a formatted message at InfoLevel.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Info(sprintf(format, args...))
}

// Warn logs a message at WarnLevel.
func (l *Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

// Warnf logs a formatted message at WarnLevel.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warn(sprintf(format, args...))
}

// Error logs a message at ErrorLevel.
func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

// Errorf logs a formatted message at ErrorLevel.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Error(sprintf(format, args...))
}

// levelToSlogLevel converts our Level to slog.Level.
func levelToSlogLevel(level Level) slog.Level {
	switch level {
	case DebugLevel:
		return slog.LevelDebug
	case InfoLevel:
		return slog.LevelInfo
	case WarnLevel:
		return slog.LevelWarn
	case ErrorLevel:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// sprintf is a helper to format strings.
func sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
