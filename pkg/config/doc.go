// Package config provides configuration management for Wayframe applications.
//
// The config package offers a simple, type-safe way to load configuration
// values from environment variables with support for defaults. It's designed
// to work with the 12-factor app methodology.
//
// # Basic Usage
//
//	cfg := config.New("APP")
//	port := cfg.String("PORT", "8080")
//	debug := cfg.Bool("DEBUG", false)
//
// # Environment Variable Naming
//
// When a prefix is provided, it's prepended to all keys with an underscore.
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
// # Example
//
//	cfg := config.New("MYAPP")
//	host := cfg.String("HOST", "localhost")
//	port := cfg.Int("PORT", 8080)
//	timeout := cfg.Duration("TIMEOUT", 30*time.Second)
//	apiKey := cfg.Required("API_KEY")  // Panics if not set
package config
