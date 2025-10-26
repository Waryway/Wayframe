package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

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

func TestDuration(t *testing.T) {
	loader := New("")
	
	// Test default value
	val := loader.Duration("NONEXISTENT_VAR", 5*time.Second)
	if val != 5*time.Second {
		t.Errorf("expected 5s, got %v", val)
	}
	
	// Test environment variable
	os.Setenv("TEST_DURATION", "10s")
	defer os.Unsetenv("TEST_DURATION")
	
	val = loader.Duration("TEST_DURATION", 5*time.Second)
	if val != 10*time.Second {
		t.Errorf("expected 10s, got %v", val)
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
