// Package config provides advanced configuration management for Wayframe applications.
// It supports struct tags, multiple file formats (JSON, YAML, key-value), environment variables, and defaults.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Loader provides methods to load configuration values from files and environment variables
// with support for defaults, type conversion, and struct tag-based loading.
//
// Use struct tags to define configuration sources and defaults:
//
//	type AppConfig struct {
//	    Port     int    `config:"port" env:"APP_PORT" default:"8080" file:"config.json"`
//	    LogLevel string `config:"log_level" env:"LOG_LEVEL" default:"INFO" file:"config.yaml"`
//	}
type Loader struct {
	values    map[string]string
	durations map[string]time.Duration
	prefix    string
}

// New creates a new configuration loader with an optional prefix for environment variables.
// The prefix is prepended to all environment variable names (e.g., "APP" -> "APP_PORT").
func New(prefix string) *Loader {
	return &Loader{
		values:    make(map[string]string),
		durations: make(map[string]time.Duration),
		prefix:    strings.ToUpper(prefix),
	}
}

// LoadFile loads configuration from a file. Supports JSON, YAML, and key-value formats.
// The format is auto-detected based on file extension or content.
func (l *Loader) LoadFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Detect format from extension
	ext := strings.ToLower(path[strings.LastIndex(path, ".")+1:])

	switch ext {
	case "json":
		return l.loadJSON(data)
	case "yaml", "yml":
		return l.loadYAML(data)
	case "env", "txt", "conf":
		return l.loadKeyValue(data)
	default:
		// Try to auto-detect
		if err := l.loadJSON(data); err == nil {
			return nil
		}
		if err := l.loadYAML(data); err == nil {
			return nil
		}
		return l.loadKeyValue(data)
	}
}

func (l *Loader) loadJSON(data []byte) error {
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	l.flattenMap("", config)
	return nil
}

func (l *Loader) loadYAML(data []byte) error {
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	l.flattenMap("", config)
	return nil
}

func (l *Loader) loadKeyValue(data []byte) error {
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			value = strings.Trim(value, `"'`)
			l.values[strings.ToUpper(key)] = value
		}
	}
	return nil
}

func (l *Loader) flattenMap(prefix string, m map[string]interface{}) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		switch val := v.(type) {
		case map[string]interface{}:
			l.flattenMap(key, val)
		default:
			l.values[strings.ToUpper(key)] = fmt.Sprintf("%v", val)
		}
	}
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
// Successfully parsed durations from config sources are cached to avoid repeated parsing.
func (l *Loader) Duration(key string, defaultValue time.Duration) time.Duration {
	key = strings.ToUpper(key)

	// Check if we already parsed this duration
	if cached, ok := l.durations[key]; ok {
		return cached
	}

	val := l.String(key, "")
	if val == "" {
		// No config value, return default without caching
		return defaultValue
	}

	duration, err := time.ParseDuration(val)
	if err != nil {
		// Parse error, return default without caching
		return defaultValue
	}

	// Cache the successfully parsed duration from config
	l.durations[key] = duration
	return duration
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

// Load populates a struct with configuration values from files, environment variables, and defaults.
// Uses struct tags: `config:"key"`, `env:"ENV_VAR"`, `default:"value"`, `file:"path"`
func (l *Loader) Load(configStruct interface{}) error {
	v := reflect.ValueOf(configStruct)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config must be a pointer to a struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		// Load file if specified
		if filePath := field.Tag.Get("file"); filePath != "" {
			l.LoadFile(filePath)
		}

		// Get configuration key
		configKey := field.Tag.Get("config")
		if configKey == "" {
			configKey = strings.ToLower(field.Name)
		}

		// Handle time.Duration fields specially using Duration() method
		if fieldValue.Kind() == reflect.Int64 && fieldValue.Type() == reflect.TypeOf(time.Duration(0)) {
			defaultValue := field.Tag.Get("default")
			var defaultDur time.Duration
			if defaultValue != "" {
				var err error
				defaultDur, err = time.ParseDuration(defaultValue)
				if err != nil {
					return fmt.Errorf("failed to parse default duration for field %s: %w", field.Name, err)
				}
				// Store default in values so Duration() can cache it properly
				upperKey := strings.ToUpper(configKey)
				l.values[upperKey] = defaultValue
			}
			// Use Duration() method which handles priority and caching
			dur := l.Duration(configKey, defaultDur)
			fieldValue.SetInt(int64(dur))
			continue
		}

		// Get environment variable name
		envKey := field.Tag.Get("env")
		if envKey == "" && l.prefix != "" {
			envKey = l.prefix + "_" + strings.ToUpper(configKey)
		} else if envKey == "" {
			envKey = strings.ToUpper(configKey)
		}

		// Get default value
		defaultValue := field.Tag.Get("default")

		// Priority: env var > file > default
		var value string
		if envVal := os.Getenv(envKey); envVal != "" {
			value = envVal
		} else if fileVal, ok := l.values[strings.ToUpper(configKey)]; ok {
			value = fileVal
		} else {
			value = defaultValue
		}

		if value == "" {
			continue
		}

		// Set the field based on its type
		if err := l.setField(fieldValue, value); err != nil {
			return fmt.Errorf("failed to set field %s: %w", field.Name, err)
		}
	}

	return nil
}

func (l *Loader) setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Note: time.Duration fields are handled separately in Load()
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			// Try common boolean strings
			value = strings.ToLower(value)
			b = value == "yes" || value == "on" || value == "1"
		}
		field.SetBool(b)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	default:
		return fmt.Errorf("unsupported field type: %v", field.Kind())
	}

	return nil
}
