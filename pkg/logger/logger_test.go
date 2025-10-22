package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoggerLevels(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(InfoLevel)
	log.SetOutput(buf)
	
	// Debug should not be logged at InfoLevel
	log.Debug("debug message")
	if strings.Contains(buf.String(), "debug message") {
		t.Error("Debug message should not be logged at InfoLevel")
	}
	
	// Info should be logged
	buf.Reset()
	log.Info("info message")
	if !strings.Contains(buf.String(), "info message") {
		t.Error("Info message should be logged")
	}
	if !strings.Contains(buf.String(), "[INFO]") {
		t.Error("Log should contain INFO level")
	}
	
	// Warn should be logged
	buf.Reset()
	log.Warn("warn message")
	if !strings.Contains(buf.String(), "warn message") {
		t.Error("Warn message should be logged")
	}
	if !strings.Contains(buf.String(), "[WARN]") {
		t.Error("Log should contain WARN level")
	}
	
	// Error should be logged
	buf.Reset()
	log.Error("error message")
	if !strings.Contains(buf.String(), "error message") {
		t.Error("Error message should be logged")
	}
	if !strings.Contains(buf.String(), "[ERROR]") {
		t.Error("Log should contain ERROR level")
	}
}

func TestDebugLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(DebugLevel)
	log.SetOutput(buf)
	
	log.Debug("debug message")
	if !strings.Contains(buf.String(), "debug message") {
		t.Error("Debug message should be logged at DebugLevel")
	}
	if !strings.Contains(buf.String(), "[DEBUG]") {
		t.Error("Log should contain DEBUG level")
	}
}

func TestWithField(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(InfoLevel)
	log.SetOutput(buf)
	
	log.WithField("key", "value").Info("message with field")
	
	output := buf.String()
	if !strings.Contains(output, "message with field") {
		t.Error("Should contain message")
	}
	if !strings.Contains(output, "key=value") {
		t.Error("Should contain field key=value")
	}
}

func TestWithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(InfoLevel)
	log.SetOutput(buf)
	
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	log.WithFields(fields).Info("message with fields")
	
	output := buf.String()
	if !strings.Contains(output, "message with fields") {
		t.Error("Should contain message")
	}
	if !strings.Contains(output, "key1=value1") {
		t.Error("Should contain field key1=value1")
	}
	if !strings.Contains(output, "key2=42") {
		t.Error("Should contain field key2=42")
	}
}

func TestFormattedLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(InfoLevel)
	log.SetOutput(buf)
	
	log.Infof("formatted message: %s, %d", "test", 123)
	
	output := buf.String()
	if !strings.Contains(output, "formatted message: test, 123") {
		t.Error("Should contain formatted message")
	}
}
