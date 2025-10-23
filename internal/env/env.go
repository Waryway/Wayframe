// Package env provides environment initialization for Wayframe applications.
// It coordinates configuration loading and logger setup.
package env

import (
	"log/slog"
	"os"

	"github.com/Waryway/Wayframe/internal/config"
	"github.com/Waryway/Wayframe/pkg/logger"
)

// Env represents the application environment with initialized config and logger.
type Env struct {
	Config *config.Config
	Logger *logger.Logger
}

// New creates a new environment with the given prefix for environment variables.
func New(prefix string) *Env {
	return &Env{
		Config: config.New(prefix),
		Logger: logger.New(logger.InfoLevel),
	}
}

// LoadConfig loads configuration into the provided struct.
// Uses struct tags for configuration: config, env, default, file
func (e *Env) LoadConfig(configStruct interface{}) error {
	return e.Config.Load(configStruct)
}

// InitLogger initializes the logger with the specified level and optional output.
// If output is nil, uses os.Stdout.
func (e *Env) InitLogger(level logger.Level) {
	e.Logger = logger.New(level)
}

// InitLoggerWithHandler initializes the logger with a custom slog handler.
func (e *Env) InitLoggerWithHandler(handler slog.Handler) {
	e.Logger = logger.NewWithHandler(handler)
}

// SetLogOutput sets the output for the logger.
func (e *Env) SetLogOutput(output *os.File) {
	if output != nil {
		e.Logger.SetOutput(output)
	}
}

// GetConfig returns the configuration manager.
func (e *Env) GetConfig() *config.Config {
	return e.Config
}

// GetLogger returns the logger.
func (e *Env) GetLogger() *logger.Logger {
	return e.Logger
}
