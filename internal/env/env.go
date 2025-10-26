// Package env provides environment initialization for Wayframe applications.
// It coordinates configuration loading and logger setup.
package env

import (
	"log/slog"
	"os"
	"time"

	"github.com/Waryway/Wayframe/internal/config"
	"github.com/Waryway/Wayframe/pkg/logger"
)

// Config represents the standard application configuration structure.
// Applications can embed this or use it directly for common configuration needs.
// Note: env tags are not specified to allow the prefix to be applied automatically.
type Config struct {
	// Server configuration
	Port            int           `config:"port" default:"8080"`
	Host            string        `config:"host" default:"0.0.0.0"`
	ReadTimeout     time.Duration `config:"read_timeout" default:"10s"`
	WriteTimeout    time.Duration `config:"write_timeout" default:"10s"`
	IdleTimeout     time.Duration `config:"idle_timeout" default:"120s"`
	ShutdownTimeout time.Duration `config:"shutdown_timeout" default:"30s"`
	
	// Logging configuration
	LogLevel string `config:"log_level" default:"INFO"`
	LogFile  string `config:"log_file" default:""`
	
	// Application configuration
	Environment string `config:"environment" default:"development"`
	Debug       bool   `config:"debug" default:"false"`
	
	// Optional config file path
	ConfigFile string `config:"config_file" default:""`
}

// Env represents the application environment with initialized config and logger.
type Env struct {
	config       *config.Config
	Logger       *logger.Logger
	AppConfig    *Config
	customConfig interface{}
}

// New creates a new environment with the given prefix for environment variables.
func New(prefix string) *Env {
	return &Env{
		config:    config.New(prefix),
		Logger:    logger.New(logger.InfoLevel),
		AppConfig: &Config{},
	}
}

// LoadConfig loads configuration into the provided struct.
// Uses struct tags for configuration: config, env, default, file
func (e *Env) LoadConfig(configStruct interface{}) error {
	e.customConfig = configStruct
	return e.config.Load(configStruct)
}

// LoadStandardConfig loads the standard Wayframe configuration.
// This should be called to populate the AppConfig field with values from
// environment variables, config files, and defaults.
func (e *Env) LoadStandardConfig() error {
	// If a config file is specified via env var, load it first
	if configFile := e.config.String("CONFIG_FILE", ""); configFile != "" {
		e.config.LoadFile(configFile)
	}
	
	// Load the standard config structure
	if err := e.config.Load(e.AppConfig); err != nil {
		return err
	}
	
	// Initialize logger based on config
	e.InitLoggerFromConfig()
	
	return nil
}

// InitLoggerFromConfig initializes the logger based on the AppConfig settings.
func (e *Env) InitLoggerFromConfig() {
	level := logger.InfoLevel
	switch e.AppConfig.LogLevel {
	case "DEBUG":
		level = logger.DebugLevel
	case "INFO":
		level = logger.InfoLevel
	case "WARN":
		level = logger.WarnLevel
	case "ERROR":
		level = logger.ErrorLevel
	}
	
	e.Logger = logger.New(level)
	
	// Set log file if specified
	if e.AppConfig.LogFile != "" {
		if f, err := os.OpenFile(e.AppConfig.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
			e.Logger.SetOutput(f)
		}
	}
}

// InitLogger initializes the logger with the specified level.
func (e *Env) InitLogger(level logger.Level) {
	e.Logger = logger.New(level)
}

// InitLoggerWithHandler initializes the logger with a custom slog handler.
func (e *Env) InitLoggerWithHandler(handler slog.Handler) {
	e.Logger = logger.NewWithHandler(handler)
}

// SetLogOutput sets the output for the logger.
func (e *Env) SetLogOutput(output *os.File) {
	if output != nil {
		e.Logger.SetOutput(output)
	}
}

// GetConfig returns the configuration manager for direct access.
func (e *Env) GetConfig() *config.Config {
	return e.config
}

// GetLogger returns the logger.
func (e *Env) GetLogger() *logger.Logger {
	return e.Logger
}

// GetAppConfig returns the standard application configuration.
func (e *Env) GetAppConfig() *Config {
	return e.AppConfig
}

// GetCustomConfig returns the custom configuration struct if one was loaded.
func (e *Env) GetCustomConfig() interface{} {
	return e.customConfig
}
