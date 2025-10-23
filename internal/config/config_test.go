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
	
	cfg := New("TEST")
	var testCfg TestConfig
	
	if err := cfg.Load(&testCfg); err != nil {
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
	
	cfg := New("TEST")
	var testCfg TestConfig
	
	if err := cfg.Load(&testCfg); err != nil {
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
	
	cfg := New("")
	// Load the file separately
	if err := cfg.LoadFile(configPath); err != nil {
		t.Fatalf("failed to load file: %v", err)
	}
	
	var testCfg TestConfig
	if err := cfg.Load(&testCfg); err != nil {
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
	
	cfg := New("")
	if err := cfg.LoadFile(configPath); err != nil {
		t.Fatalf("failed to load YAML file: %v", err)
	}
	
	port := cfg.Int("port", 8080)
	if port != 7777 {
		t.Errorf("expected port 7777 from YAML, got %d", port)
	}
	
	host := cfg.String("host", "localhost")
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
	
	cfg := New("")
	if err := cfg.LoadFile(configPath); err != nil {
		t.Fatalf("failed to load key-value file: %v", err)
	}
	
	port := cfg.Int("port", 8080)
	if port != 6666 {
		t.Errorf("expected port 6666 from env file, got %d", port)
	}
	
	host := cfg.String("host", "localhost")
	if host != "env.example.com" {
		t.Errorf("expected host env.example.com from env file, got %s", host)
	}
	
	debug := cfg.Bool("debug", false)
	if debug != true {
		t.Errorf("expected debug true from env file, got %v", debug)
	}
}

func TestDirectAccessMethods(t *testing.T) {
	cfg := New("")
	
	// String
	val := cfg.String("NONEXISTENT", "default")
	if val != "default" {
		t.Errorf("expected 'default', got '%s'", val)
	}
	
	// Int
	intVal := cfg.Int("NONEXISTENT", 42)
	if intVal != 42 {
		t.Errorf("expected 42, got %d", intVal)
	}
	
	// Bool
	boolVal := cfg.Bool("NONEXISTENT", true)
	if boolVal != true {
		t.Errorf("expected true, got %v", boolVal)
	}
	
	// Duration
	durVal := cfg.Duration("NONEXISTENT", 5*time.Second)
	if durVal != 5*time.Second {
		t.Errorf("expected 5s, got %v", durVal)
	}
}
