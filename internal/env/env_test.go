package env

import (
	"os"
	"testing"
	"time"
)

func TestLoadStandardConfig(t *testing.T) {
	// Clean up environment
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("HOST")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("DEBUG")
	}()
	
	e := New("")
	
	// Load with defaults
	if err := e.LoadStandardConfig(); err != nil {
		t.Fatalf("failed to load standard config: %v", err)
	}
	
	// Check defaults
	if e.AppConfig.Port != 8080 {
		t.Errorf("expected port 8080, got %d", e.AppConfig.Port)
	}
	if e.AppConfig.Host != "0.0.0.0" {
		t.Errorf("expected host 0.0.0.0, got %s", e.AppConfig.Host)
	}
	if e.AppConfig.LogLevel != "INFO" {
		t.Errorf("expected log level INFO, got %s", e.AppConfig.LogLevel)
	}
	if e.AppConfig.Debug != false {
		t.Errorf("expected debug false, got %v", e.AppConfig.Debug)
	}
}

func TestLoadStandardConfigWithEnvVars(t *testing.T) {
	// Set environment variables (no prefix)
	os.Setenv("PORT", "9000")
	os.Setenv("HOST", "localhost")
	os.Setenv("LOG_LEVEL", "DEBUG")
	os.Setenv("DEBUG", "true")
	
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("HOST")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("DEBUG")
	}()
	
	e := New("")
	
	if err := e.LoadStandardConfig(); err != nil {
		t.Fatalf("failed to load standard config: %v", err)
	}
	
	// Check env var values
	if e.AppConfig.Port != 9000 {
		t.Errorf("expected port 9000 from env, got %d", e.AppConfig.Port)
	}
	if e.AppConfig.Host != "localhost" {
		t.Errorf("expected host localhost from env, got %s", e.AppConfig.Host)
	}
	if e.AppConfig.LogLevel != "DEBUG" {
		t.Errorf("expected log level DEBUG from env, got %s", e.AppConfig.LogLevel)
	}
	if e.AppConfig.Debug != true {
		t.Errorf("expected debug true from env, got %v", e.AppConfig.Debug)
	}
}

func TestLoadCustomConfig(t *testing.T) {
	type CustomConfig struct {
		AppName string `config:"app_name" env:"APP_NAME" default:"myapp"`
		Version string `config:"version" env:"VERSION" default:"1.0.0"`
	}
	
	e := New("")
	var cfg CustomConfig
	
	if err := e.LoadConfig(&cfg); err != nil {
		t.Fatalf("failed to load custom config: %v", err)
	}
	
	if cfg.AppName != "myapp" {
		t.Errorf("expected app_name myapp, got %s", cfg.AppName)
	}
	if cfg.Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", cfg.Version)
	}
	
	// Check that custom config is stored
	if e.GetCustomConfig() == nil {
		t.Error("expected custom config to be stored")
	}
}

func TestInitLoggerFromConfig(t *testing.T) {
	e := New("")
	e.AppConfig.LogLevel = "DEBUG"
	
	e.InitLoggerFromConfig()
	
	if e.Logger == nil {
		t.Error("expected logger to be initialized")
	}
}

func TestConfigStructTags(t *testing.T) {
	// Set some env vars with prefix
	os.Setenv("APP_PORT", "3000")
	os.Setenv("APP_LOG_LEVEL", "WARN")
	
	defer func() {
		os.Unsetenv("APP_PORT")
		os.Unsetenv("APP_LOG_LEVEL")
	}()
	
	e := New("APP")
	if err := e.LoadStandardConfig(); err != nil {
		t.Fatalf("failed to load config with prefix: %v", err)
	}
	
	if e.AppConfig.Port != 3000 {
		t.Errorf("expected port 3000 with prefix, got %d", e.AppConfig.Port)
	}
	if e.AppConfig.LogLevel != "WARN" {
		t.Errorf("expected log level WARN with prefix, got %s", e.AppConfig.LogLevel)
	}
}

func TestTimeoutDefaults(t *testing.T) {
	e := New("")
	if err := e.LoadStandardConfig(); err != nil {
		t.Fatalf("failed to load standard config: %v", err)
	}
	
	if e.AppConfig.ReadTimeout != 10*time.Second {
		t.Errorf("expected read timeout 10s, got %v", e.AppConfig.ReadTimeout)
	}
	if e.AppConfig.WriteTimeout != 10*time.Second {
		t.Errorf("expected write timeout 10s, got %v", e.AppConfig.WriteTimeout)
	}
	if e.AppConfig.IdleTimeout != 120*time.Second {
		t.Errorf("expected idle timeout 120s, got %v", e.AppConfig.IdleTimeout)
	}
	if e.AppConfig.ShutdownTimeout != 30*time.Second {
		t.Errorf("expected shutdown timeout 30s, got %v", e.AppConfig.ShutdownTimeout)
	}
}
