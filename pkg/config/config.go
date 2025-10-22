// Package config provides configuration management for Wayframe applications.
// It supports environment variables, defaults, and type-safe configuration loading.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Loader provides methods to load configuration values from environment variables
// with support for defaults and type conversion.
type Loader struct {
	prefix string
}

// New creates a new configuration loader with an optional prefix for environment variables.
// The prefix is prepended to all environment variable names (e.g., "APP" -> "APP_PORT").
func New(prefix string) *Loader {
	return &Loader{
		prefix: strings.ToUpper(prefix),
	}
}

// String loads a string configuration value from the environment.
// Returns the default value if the environment variable is not set.
func (l *Loader) String(key, defaultValue string) string {
	envKey := l.buildKey(key)
	if val := os.Getenv(envKey); val != "" {
		return val
	}
	return defaultValue
}

// Int loads an integer configuration value from the environment.
// Returns the default value if the environment variable is not set or cannot be parsed.
func (l *Loader) Int(key string, defaultValue int) int {
	envKey := l.buildKey(key)
	if val := os.Getenv(envKey); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// Bool loads a boolean configuration value from the environment.
// Accepts: "true", "1", "yes", "on" (case-insensitive) as true.
// Returns the default value if the environment variable is not set or cannot be parsed.
func (l *Loader) Bool(key string, defaultValue bool) bool {
	envKey := l.buildKey(key)
	if val := os.Getenv(envKey); val != "" {
		val = strings.ToLower(val)
		if val == "true" || val == "1" || val == "yes" || val == "on" {
			return true
		}
		if val == "false" || val == "0" || val == "no" || val == "off" {
			return false
		}
	}
	return defaultValue
}

// Duration loads a duration configuration value from the environment.
// Accepts values like "1s", "5m", "1h" as per time.ParseDuration.
// Returns the default value if the environment variable is not set or cannot be parsed.
func (l *Loader) Duration(key string, defaultValue time.Duration) time.Duration {
	envKey := l.buildKey(key)
	if val := os.Getenv(envKey); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}
	return defaultValue
}

// Required loads a required string configuration value.
// Panics if the environment variable is not set.
func (l *Loader) Required(key string) string {
	envKey := l.buildKey(key)
	val := os.Getenv(envKey)
	if val == "" {
		panic(fmt.Sprintf("required configuration %s is not set", envKey))
	}
	return val
}

// buildKey constructs the full environment variable name with prefix.
func (l *Loader) buildKey(key string) string {
	key = strings.ToUpper(key)
	if l.prefix != "" {
		return l.prefix + "_" + key
	}
	return key
}
