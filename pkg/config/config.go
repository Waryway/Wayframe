// Package config provides configuration management for Wayframe applications.
// It supports loading from files, environment variables, and defaults with type-safe loading.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Loader provides methods to load configuration values from files and environment variables
// with support for defaults and type conversion.
type Loader struct {
	prefix string
	values map[string]string
}

// New creates a new configuration loader with an optional prefix for environment variables.
// The prefix is prepended to all environment variable names (e.g., "APP" -> "APP_PORT").
func New(prefix string) *Loader {
	return &Loader{
		prefix: strings.ToUpper(prefix),
		values: make(map[string]string),
	}
}

// LoadFile loads configuration from a JSON file at the given path.
// If the file cannot be loaded or parsed, it returns an error but the loader
// can still be used with environment variables and defaults as fallback.
func (l *Loader) LoadFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Flatten the config map into string values
	for k, v := range config {
		l.values[strings.ToUpper(k)] = fmt.Sprintf("%v", v)
	}

	return nil
}

// String loads a string configuration value.
// Priority: 1) Environment variable, 2) File value, 3) Default value.
// The environment variable name matches the key name (with prefix if set).
func (l *Loader) String(key, defaultValue string) string {
	key = strings.ToUpper(key)
	
	// Check environment variable first
	envKey := l.buildKey(key)
	if val := os.Getenv(envKey); val != "" {
		return val
	}
	
	// Check loaded file values
	if val, ok := l.values[key]; ok {
		return val
	}
	
	// Return default
	return defaultValue
}

// Int loads an integer configuration value.
// Priority: 1) Environment variable, 2) File value, 3) Default value.
// Returns the default value if the value cannot be parsed.
func (l *Loader) Int(key string, defaultValue int) int {
	val := l.String(key, "")
	if val == "" {
		return defaultValue
	}
	
	if intVal, err := strconv.Atoi(val); err == nil {
		return intVal
	}
	
	return defaultValue
}

// Bool loads a boolean configuration value.
// Priority: 1) Environment variable, 2) File value, 3) Default value.
// Accepts: "true", "1", "yes", "on" (case-insensitive) as true.
// Returns the default value if the value cannot be parsed.
func (l *Loader) Bool(key string, defaultValue bool) bool {
	val := l.String(key, "")
	if val == "" {
		return defaultValue
	}
	
	val = strings.ToLower(val)
	if val == "true" || val == "1" || val == "yes" || val == "on" {
		return true
	}
	if val == "false" || val == "0" || val == "no" || val == "off" {
		return false
	}
	
	return defaultValue
}

// Duration loads a duration configuration value.
// Priority: 1) Environment variable, 2) File value, 3) Default value.
// Accepts values like "1s", "5m", "1h" as per time.ParseDuration.
// Returns the default value if the value cannot be parsed.
func (l *Loader) Duration(key string, defaultValue time.Duration) time.Duration {
	val := l.String(key, "")
	if val == "" {
		return defaultValue
	}
	
	if duration, err := time.ParseDuration(val); err == nil {
		return duration
	}
	
	return defaultValue
}

// Required loads a required string configuration value.
// Priority: 1) Environment variable, 2) File value.
// Panics if the value is not set in either location.
func (l *Loader) Required(key string) string {
	val := l.String(key, "")
	if val == "" {
		envKey := l.buildKey(strings.ToUpper(key))
		panic(fmt.Sprintf("required configuration %s is not set", envKey))
	}
	return val
}

// buildKey constructs the full environment variable name with prefix.
func (l *Loader) buildKey(key string) string {
	if l.prefix != "" {
		return l.prefix + "_" + key
	}
	return key
}
