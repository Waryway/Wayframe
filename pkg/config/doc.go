// Package config provides configuration management for Wayframe applications.
//
// The config package offers a flexible, type-safe way to load configuration
// values from JSON files and environment variables with support for defaults.
// It's designed to work with the 12-factor app methodology while also supporting
// file-based configuration.
//
// # Basic Usage
//
//	cfg := config.New("APP")
//	port := cfg.String("PORT", "8080")
//	debug := cfg.Bool("DEBUG", false)
//
// # Loading from Files
//
// Configuration can be loaded from JSON files. If the file fails to load,
// the loader will fall back to environment variables and defaults:
//
//	cfg := config.New("APP")
//	if err := cfg.LoadFile("/path/to/config.json"); err != nil {
//	    // Log error but continue - env vars and defaults will be used
//	}
//	port := cfg.String("PORT", "8080")
//
// # Priority Order
//
// Configuration values are resolved in this priority order:
//  1. Environment variables (highest priority)
//  2. Values from loaded JSON file
//  3. Default values provided in the code (lowest priority)
//
// # Environment Variable Naming
//
// When a prefix is provided, it's prepended to all keys with an underscore.
// Environment variable names match the configuration key names.
// For example, with prefix "APP":
//   - Key "PORT" becomes environment variable "APP_PORT"
//   - Key "DEBUG" becomes environment variable "APP_DEBUG"
//
// # Type Safety
//
// The package provides type-safe loading methods:
//   - String: Load string values
//   - Int: Load integer values (with validation)
//   - Bool: Load boolean values (supports true/false, 1/0, yes/no, on/off)
//   - Duration: Load time.Duration values (e.g., "30s", "5m", "1h")
//   - Required: Load required string values (panics if not set)
//
// # Example with File Loading
//
//	cfg := config.New("MYAPP")
//
//	// Try to load from file (optional)
//	if err := cfg.LoadFile("config.json"); err != nil {
//	    log.Printf("Config file not found, using env vars and defaults: %v", err)
//	}
//
//	host := cfg.String("HOST", "localhost")
//	port := cfg.Int("PORT", 8080)
//	timeout := cfg.Duration("TIMEOUT", 30*time.Second)
//	apiKey := cfg.Required("API_KEY")  // Panics if not set anywhere
package config
