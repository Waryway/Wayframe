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

// Config represents a configuration structure that can be loaded from various sources.
// Use struct tags to define configuration sources and defaults:
//
//	type AppConfig struct {
//	    Port     int    `config:"port" env:"APP_PORT" default:"8080" file:"config.json"`
//	    LogLevel string `config:"log_level" env:"LOG_LEVEL" default:"INFO" file:"config.yaml"`
//	}
type Config struct {
	values map[string]string
	prefix string
}

// New creates a new configuration manager with an optional environment variable prefix.
func New(prefix string) *Config {
	return &Config{
		values: make(map[string]string),
		prefix: strings.ToUpper(prefix),
	}
}

// LoadFile loads configuration from a file. Supports JSON, YAML, and key-value formats.
// The format is auto-detected based on file extension or content.
func (c *Config) LoadFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Detect format from extension
	ext := strings.ToLower(path[strings.LastIndex(path, ".")+1:])
	
	switch ext {
	case "json":
		return c.loadJSON(data)
	case "yaml", "yml":
		return c.loadYAML(data)
	case "env", "txt", "conf":
		return c.loadKeyValue(data)
	default:
		// Try to auto-detect
		if err := c.loadJSON(data); err == nil {
			return nil
		}
		if err := c.loadYAML(data); err == nil {
			return nil
		}
		return c.loadKeyValue(data)
	}
}

func (c *Config) loadJSON(data []byte) error {
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	c.flattenMap("", config)
	return nil
}

func (c *Config) loadYAML(data []byte) error {
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}
	
	c.flattenMap("", config)
	return nil
}

func (c *Config) loadKeyValue(data []byte) error {
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
			c.values[strings.ToUpper(key)] = value
		}
	}
	return nil
}

func (c *Config) flattenMap(prefix string, m map[string]interface{}) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		
		switch val := v.(type) {
		case map[string]interface{}:
			c.flattenMap(key, val)
		default:
			c.values[strings.ToUpper(key)] = fmt.Sprintf("%v", val)
		}
	}
}

// Load populates a struct with configuration values from files, environment variables, and defaults.
// Uses struct tags: `config:"key"`, `env:"ENV_VAR"`, `default:"value"`, `file:"path"`
func (c *Config) Load(configStruct interface{}) error {
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
			c.LoadFile(filePath)
		}
		
		// Get configuration key
		configKey := field.Tag.Get("config")
		if configKey == "" {
			configKey = strings.ToLower(field.Name)
		}
		
		// Get environment variable name
		envKey := field.Tag.Get("env")
		if envKey == "" && c.prefix != "" {
			envKey = c.prefix + "_" + strings.ToUpper(configKey)
		} else if envKey == "" {
			envKey = strings.ToUpper(configKey)
		}
		
		// Get default value
		defaultValue := field.Tag.Get("default")
		
		// Priority: env var > file > default
		var value string
		if envVal := os.Getenv(envKey); envVal != "" {
			value = envVal
		} else if fileVal, ok := c.values[strings.ToUpper(configKey)]; ok {
			value = fileVal
		} else {
			value = defaultValue
		}
		
		if value == "" {
			continue
		}
		
		// Set the field based on its type
		if err := c.setField(fieldValue, value); err != nil {
			return fmt.Errorf("failed to set field %s: %w", field.Name, err)
		}
	}
	
	return nil
}

func (c *Config) setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			d, err := time.ParseDuration(value)
			if err != nil {
				return err
			}
			field.SetInt(int64(d))
		} else {
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(i)
		}
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

// String loads a string configuration value.
// Priority: 1) Environment variable, 2) File value, 3) Default value.
func (c *Config) String(key, defaultValue string) string {
	key = strings.ToUpper(key)
	
	// Check environment variable first
	envKey := c.buildKey(key)
	if val := os.Getenv(envKey); val != "" {
		return val
	}
	
	// Check loaded file values
	if val, ok := c.values[key]; ok {
		return val
	}
	
	// Return default
	return defaultValue
}

// Int loads an integer configuration value.
func (c *Config) Int(key string, defaultValue int) int {
	val := c.String(key, "")
	if val == "" {
		return defaultValue
	}
	
	if intVal, err := strconv.Atoi(val); err == nil {
		return intVal
	}
	
	return defaultValue
}

// Bool loads a boolean configuration value.
func (c *Config) Bool(key string, defaultValue bool) bool {
	val := c.String(key, "")
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
func (c *Config) Duration(key string, defaultValue time.Duration) time.Duration {
	val := c.String(key, "")
	if val == "" {
		return defaultValue
	}
	
	if duration, err := time.ParseDuration(val); err == nil {
		return duration
	}
	
	return defaultValue
}

// Required loads a required string configuration value.
func (c *Config) Required(key string) string {
	val := c.String(key, "")
	if val == "" {
		envKey := c.buildKey(strings.ToUpper(key))
		panic(fmt.Sprintf("required configuration %s is not set", envKey))
	}
	return val
}

func (c *Config) buildKey(key string) string {
	if c.prefix != "" {
		return c.prefix + "_" + key
	}
	return key
}
