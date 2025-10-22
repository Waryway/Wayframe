package config

import (
	"os"
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
