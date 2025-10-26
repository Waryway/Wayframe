// Package logger provides structured logging for Wayframe applications.
//
// The logger package wraps Go's standard log/slog package with a simplified
// interface. It offers leveled logging with contextual fields support and
// is designed to be simple yet powerful enough for production use.
//
// # Basic Usage
//
//	log := logger.New(logger.InfoLevel)
//	log.Info("Application started")
//	log.Errorf("Failed to connect: %v", err)
//
// # Log Levels
//
// The package supports four log levels:
//   - DebugLevel: Verbose information, typically disabled in production
//   - InfoLevel: General informational messages (default)
//   - WarnLevel: Warning messages for potentially harmful situations
//   - ErrorLevel: Error messages for serious problems
//
// Only messages at or above the configured level will be logged.
//
// # Contextual Logging
//
// Add contextual fields to log messages:
//
//	log.WithField("user_id", 123).Info("User logged in")
//	log.WithFields(map[string]interface{}{
//	    "user_id": 123,
//	    "ip": "192.168.1.1",
//	}).Info("User logged in")
//
// # Formatted Logging
//
// All levels support formatted messages:
//
//	log.Infof("Processing %d items", count)
//	log.Errorf("Connection failed: %v", err)
//
// # Using slog Directly
//
// For advanced use cases, you can create a logger with a custom slog.Handler:
//
//	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
//	    Level: slog.LevelInfo,
//	})
//	log := logger.NewWithHandler(handler)
//
// # Output Format
//
// By default, log messages use slog's text format:
//   time=2025-10-22T16:00:00.000Z level=INFO msg="message" field1=value1 field2=value2
package logger
