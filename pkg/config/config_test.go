package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStructTagLoading(t *testing.T) {
	type TestConfig struct {
		Port     int           `config:"port" env:"TEST_PORT" default:"8080"`
		Host     string        `config:"host" env:"TEST_HOST" default:"localhost"`
		Debug    bool          `config:"debug" env:"TEST_DEBUG" default:"false"`
		Timeout  time.Duration `config:"timeout" env:"TEST_TIMEOUT" default:"30s"`
		MaxConns int           `config:"max_conns" default:"100"`
	}
	
	loader := New("TEST")
	var testCfg TestConfig
	
	if err := loader.Load(&testCfg); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	
	// Check defaults
	if testCfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", testCfg.Port)
	}
	if testCfg.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", testCfg.Host)
	}
	if testCfg.Debug != false {
		t.Errorf("expected debug false, got %v", testCfg.Debug)
	}
	if testCfg.Timeout != 30*time.Second {
		t.Errorf("expected timeout 30s, got %v", testCfg.Timeout)
	}
	if testCfg.MaxConns != 100 {
		t.Errorf("expected max_conns 100, got %d", testCfg.MaxConns)
	}
}

func TestEnvVarOverride(t *testing.T) {
	type TestConfig struct {
		Port int `config:"port" env:"TEST_PORT" default:"8080"`
	}
	
	os.Setenv("TEST_PORT", "9000")
	defer os.Unsetenv("TEST_PORT")
	
	loader := New("TEST")
	var testCfg TestConfig
	
	if err := loader.Load(&testCfg); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	
	if testCfg.Port != 9000 {
		t.Errorf("expected port 9000 from env var, got %d", testCfg.Port)
	}
}

func TestJSONFileLoading(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	
	jsonData := `{
		"port": 9999,
		"host": "example.com",
		"debug": true
	}`
	
	if err := os.WriteFile(configPath, []byte(jsonData), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	
	type TestConfig struct {
		Port  int    `config:"port" default:"8080"`
		Host  string `config:"host" default:"localhost"`
		Debug bool   `config:"debug" default:"false"`
	}
	
	loader := New("")
	// Load the file separately
	if err := loader.LoadFile(configPath); err != nil {
		t.Fatalf("failed to load file: %v", err)
	}
	
	var testCfg TestConfig
	if err := loader.Load(&testCfg); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	
	if testCfg.Port != 9999 {
		t.Errorf("expected port 9999 from file, got %d", testCfg.Port)
	}
	if testCfg.Host != "example.com" {
		t.Errorf("expected host example.com from file, got %s", testCfg.Host)
	}
	if testCfg.Debug != true {
		t.Errorf("expected debug true from file, got %v", testCfg.Debug)
	}
}

func TestYAMLFileLoading(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	
	yamlData := `port: 7777
host: yaml.example.com
debug: true
`
	
	if err := os.WriteFile(configPath, []byte(yamlData), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	
	loader := New("")
	if err := loader.LoadFile(configPath); err != nil {
		t.Fatalf("failed to load YAML file: %v", err)
	}
	
	port := loader.Int("port", 8080)
	if port != 7777 {
		t.Errorf("expected port 7777 from YAML, got %d", port)
	}
	
	host := loader.String("host", "localhost")
	if host != "yaml.example.com" {
		t.Errorf("expected host yaml.example.com from YAML, got %s", host)
	}
}

func TestKeyValueFileLoading(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.env")
	
	envData := `PORT=6666
HOST=env.example.com
DEBUG=yes
# Comment line
TIMEOUT=45s
`
	
	if err := os.WriteFile(configPath, []byte(envData), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	
	loader := New("")
	if err := loader.LoadFile(configPath); err != nil {
		t.Fatalf("failed to load key-value file: %v", err)
	}
	
	port := loader.Int("port", 8080)
	if port != 6666 {
		t.Errorf("expected port 6666 from env file, got %d", port)
	}
	
	host := loader.String("host", "localhost")
	if host != "env.example.com" {
		t.Errorf("expected host env.example.com from env file, got %s", host)
	}
	
	debug := loader.Bool("debug", false)
	if debug != true {
		t.Errorf("expected debug true from env file, got %v", debug)
	}
}

func TestDirectAccessMethods(t *testing.T) {
	loader := New("")
	
	// String
	val := loader.String("NONEXISTENT", "default")
	if val != "default" {
		t.Errorf("expected 'default', got '%s'", val)
	}
	
	// Int
	intVal := loader.Int("NONEXISTENT", 42)
	if intVal != 42 {
		t.Errorf("expected 42, got %d", intVal)
	}
	
	// Bool
	boolVal := loader.Bool("NONEXISTENT", true)
	if boolVal != true {
		t.Errorf("expected true, got %v", boolVal)
	}
	
	// Duration
	durVal := loader.Duration("NONEXISTENT", 5*time.Second)
	if durVal != 5*time.Second {
		t.Errorf("expected 5s, got %v", durVal)
	}
}

func TestString(t *testing.T) {
	loader := New("")
	
	// Test default value
	val := loader.String("NONEXISTENT_VAR", "default")
	if val != "default" {
		t.Errorf("expected 'default', got '%s'", val)
	}
	
	// Test environment variable
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")
	
	val = loader.String("TEST_VAR", "default")
	if val != "test_value" {
		t.Errorf("expected 'test_value', got '%s'", val)
	}
}

func TestInt(t *testing.T) {
	loader := New("")
	
	// Test default value
	val := loader.Int("NONEXISTENT_VAR", 42)
	if val != 42 {
		t.Errorf("expected 42, got %d", val)
	}
	
	// Test environment variable
	os.Setenv("TEST_INT", "123")
	defer os.Unsetenv("TEST_INT")
	
	val = loader.Int("TEST_INT", 42)
	if val != 123 {
		t.Errorf("expected 123, got %d", val)
	}
	
	// Test invalid int
	os.Setenv("TEST_INT_INVALID", "not_a_number")
	defer os.Unsetenv("TEST_INT_INVALID")
	
	val = loader.Int("TEST_INT_INVALID", 42)
	if val != 42 {
		t.Errorf("expected default 42 for invalid int, got %d", val)
	}
}

func TestBool(t *testing.T) {
	loader := New("")
	
	tests := []struct {
		envVal   string
		expected bool
	}{
		{"true", true},
		{"TRUE", true},
		{"1", true},
		{"yes", true},
		{"on", true},
		{"false", false},
		{"FALSE", false},
		{"0", false},
		{"no", false},
		{"off", false},
	}
	
	for _, tt := range tests {
		os.Setenv("TEST_BOOL", tt.envVal)
		val := loader.Bool("TEST_BOOL", false)
		if val != tt.expected {
			t.Errorf("for value '%s', expected %v, got %v", tt.envVal, tt.expected, val)
		}
		os.Unsetenv("TEST_BOOL")
	}
}

func TestPrefix(t *testing.T) {
	loader := New("APP")
	
	os.Setenv("APP_PORT", "8080")
	defer os.Unsetenv("APP_PORT")
	
	val := loader.String("PORT", "3000")
	if val != "8080" {
		t.Errorf("expected '8080', got '%s'", val)
	}
}

func TestRequired(t *testing.T) {
	loader := New("")
	
	os.Setenv("REQUIRED_VAR", "value")
	defer os.Unsetenv("REQUIRED_VAR")
	
	val := loader.Required("REQUIRED_VAR")
	if val != "value" {
		t.Errorf("expected 'value', got '%s'", val)
	}
	
	// Test panic for missing required
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for missing required variable")
		}
	}()
	loader.Required("MISSING_REQUIRED_VAR")
}

func TestLoadFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	
	configData := `{
		"port": "9000",
		"debug": "true",
		"timeout": "60s"
	}`
	
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("failed to create test config file: %v", err)
	}
	
	loader := New("")
	if err := loader.LoadFile(configPath); err != nil {
		t.Fatalf("failed to load config file: %v", err)
	}
	
	// Test values from file
	if val := loader.String("port", "8080"); val != "9000" {
		t.Errorf("expected '9000' from file, got '%s'", val)
	}
	
	if val := loader.Bool("debug", false); val != true {
		t.Errorf("expected true from file, got %v", val)
	}
	
	if val := loader.Duration("timeout", 30*time.Second); val != 60*time.Second {
		t.Errorf("expected 60s from file, got %v", val)
	}
}

func TestLoadFileWithEnvOverride(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	
	configData := `{
		"port": "9000"
	}`
	
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("failed to create test config file: %v", err)
	}
	
	loader := New("")
	if err := loader.LoadFile(configPath); err != nil {
		t.Fatalf("failed to load config file: %v", err)
	}
	
	// Environment variable should override file value
	os.Setenv("PORT", "8888")
	defer os.Unsetenv("PORT")
	
	if val := loader.String("port", "8080"); val != "8888" {
		t.Errorf("expected '8888' from env (override), got '%s'", val)
	}
}

func TestLoadFileNotFound(t *testing.T) {
	loader := New("")
	err := loader.LoadFile("/nonexistent/path/config.json")
	
	if err == nil {
		t.Error("expected error for non-existent file")
	}
	
	// Loader should still work with env vars and defaults
	if val := loader.String("test", "default"); val != "default" {
		t.Errorf("expected default value after failed file load, got '%s'", val)
	}
}
