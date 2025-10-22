// Package logger provides structured logging for Wayframe applications.
// It offers a simple, leveled logging interface with contextual fields.
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
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

var levelNames = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
}

// Logger provides structured logging capabilities.
type Logger struct {
	level  Level
	out    io.Writer
	mu     sync.Mutex
	fields map[string]interface{}
}

// New creates a new Logger with the specified minimum level.
// Logs with a level lower than the minimum will be discarded.
func New(level Level) *Logger {
	return &Logger{
		level:  level,
		out:    os.Stdout,
		fields: make(map[string]interface{}),
	}
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

// WithField creates a new logger with an additional contextual field.
func (l *Logger) WithField(key string, value interface{}) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	fields := make(map[string]interface{}, len(l.fields)+1)
	for k, v := range l.fields {
		fields[k] = v
	}
	fields[key] = value
	
	return &Logger{
		level:  l.level,
		out:    l.out,
		fields: fields,
	}
}

// WithFields creates a new logger with multiple contextual fields.
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	newFields := make(map[string]interface{}, len(l.fields)+len(fields))
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}
	
	return &Logger{
		level:  l.level,
		out:    l.out,
		fields: newFields,
	}
}

// Debug logs a message at DebugLevel.
func (l *Logger) Debug(msg string) {
	l.log(DebugLevel, msg)
}

// Debugf logs a formatted message at DebugLevel.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(format, args...))
}

// Info logs a message at InfoLevel.
func (l *Logger) Info(msg string) {
	l.log(InfoLevel, msg)
}

// Infof logs a formatted message at InfoLevel.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(format, args...))
}

// Warn logs a message at WarnLevel.
func (l *Logger) Warn(msg string) {
	l.log(WarnLevel, msg)
}

// Warnf logs a formatted message at WarnLevel.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(format, args...))
}

// Error logs a message at ErrorLevel.
func (l *Logger) Error(msg string) {
	l.log(ErrorLevel, msg)
}

// Errorf logs a formatted message at ErrorLevel.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...))
}

func (l *Logger) log(level Level, msg string) {
	if level < l.level {
		return
	}
	
	l.mu.Lock()
	defer l.mu.Unlock()
	
	timestamp := time.Now().Format(time.RFC3339)
	levelName := levelNames[level]
	
	// Build log message with fields
	logMsg := fmt.Sprintf("%s [%s] %s", timestamp, levelName, msg)
	
	if len(l.fields) > 0 {
		logMsg += " |"
		for k, v := range l.fields {
			logMsg += fmt.Sprintf(" %s=%v", k, v)
		}
	}
	
	log.New(l.out, "", 0).Println(logMsg)
}
