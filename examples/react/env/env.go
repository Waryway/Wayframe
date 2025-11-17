// Package env provides a tiny environment/config helper for the React example.
package env

import (
	"time"

	"github.com/Waryway/Wayframe/pkg/config"
	"github.com/Waryway/Wayframe/pkg/logger"
)

// Env wraps a config loader and a structured logger.
type Env struct {
	cfg *config.Loader
	log *logger.Logger
}

// New creates a new Env with the given environment variable prefix (e.g., APP).
func New(prefix string) *Env {
	cfg := config.New(prefix)
	levelStr := cfg.String("LOG_LEVEL", "INFO")
	level := logger.InfoLevel
	switch levelStr {
	case "DEBUG":
		level = logger.DebugLevel
	case "WARN":
		level = logger.WarnLevel
	case "ERROR":
		level = logger.ErrorLevel
	}
	return &Env{cfg: cfg, log: logger.New(level)}
}

// LoadFiles attempts to load one or more config files. Non-existent files are ignored.
func (e *Env) LoadFiles(paths ...string) {
	for _, p := range paths {
		_ = e.cfg.LoadFile(p)
	}
}

// Logger returns the configured logger.
func (e *Env) Logger() *logger.Logger { return e.log }

// String returns a string value with default.
func (e *Env) String(key, def string) string { return e.cfg.String(key, def) }

// Duration returns a duration value with default.
func (e *Env) Duration(key string, def time.Duration) time.Duration { return e.cfg.Duration(key, def) }
