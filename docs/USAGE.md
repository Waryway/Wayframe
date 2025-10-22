# Wayframe Usage Guide

This guide demonstrates how each Wayframe package can be used independently or together.

## Using Packages Independently

Each Wayframe package is designed to be used standalone without dependencies on other Wayframe packages.

### Config Package Only

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/Waryway/Wayframe/pkg/config"
)

func main() {
    // Create a loader with prefix
    cfg := config.New("MYAPP")
    
    // Load with defaults
    port := cfg.String("PORT", "3000")
    maxConns := cfg.Int("MAX_CONNECTIONS", 100)
    timeout := cfg.Duration("TIMEOUT", 30*time.Second)
    debug := cfg.Bool("DEBUG", false)
    
    // Load required values
    apiKey := cfg.Required("API_KEY")
    
    fmt.Printf("Port: %s, Max Connections: %d\n", port, maxConns)
    fmt.Printf("Timeout: %v, Debug: %v\n", timeout, debug)
    fmt.Printf("API Key: %s\n", apiKey)
}
```

### Logger Package Only

```go
package main

import (
    "github.com/Waryway/Wayframe/pkg/logger"
)

func main() {
    // Create logger
    log := logger.New(logger.InfoLevel)
    
    // Simple logging
    log.Info("Application started")
    log.Warn("This is a warning")
    log.Error("An error occurred")
    
    // Formatted logging
    count := 42
    log.Infof("Processing %d items", count)
    
    // Contextual logging
    log.WithField("user_id", 123).Info("User action")
    log.WithFields(map[string]interface{}{
        "request_id": "abc-123",
        "method": "POST",
    }).Info("Request received")
}
```

### Server Package Only

```go
package main

import (
    "fmt"
    "net/http"
    "time"
    
    "github.com/Waryway/Wayframe/pkg/server"
)

func main() {
    // Create server
    srv := server.New(server.Config{
        Addr:         ":8080",
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    })
    
    // Register handlers
    srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })
    
    srv.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "OK")
    })
    
    // Start with graceful shutdown
    if err := srv.Start(30 * time.Second); err != nil {
        fmt.Printf("Server error: %v\n", err)
    }
}
```

## Using Packages Together

When used together, the packages complement each other to provide a complete application foundation.

```go
package main

import (
    "fmt"
    "net/http"
    "time"
    
    "github.com/Waryway/Wayframe/pkg/config"
    "github.com/Waryway/Wayframe/pkg/logger"
    "github.com/Waryway/Wayframe/pkg/server"
)

func main() {
    // 1. Load configuration
    cfg := config.New("APP")
    port := cfg.String("PORT", "8080")
    logLevel := cfg.String("LOG_LEVEL", "INFO")
    shutdownTimeout := cfg.Duration("SHUTDOWN_TIMEOUT", 30*time.Second)
    
    // 2. Setup logger
    level := logger.InfoLevel
    if logLevel == "DEBUG" {
        level = logger.DebugLevel
    }
    log := logger.New(level)
    log.Info("Starting application")
    
    // 3. Create and configure server
    srv := server.New(server.Config{
        Addr:         fmt.Sprintf(":%s", port),
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    })
    
    // Add middleware
    srv.Use(server.LoggingMiddleware(log))
    srv.Use(server.RecoveryMiddleware(log))
    
    // Register routes
    srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Debug("Handling root request")
        fmt.Fprintf(w, "Welcome!")
    })
    
    srv.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
        log.WithField("endpoint", "/api/status").Info("Status check")
        fmt.Fprintf(w, `{"status":"ok"}`)
    })
    
    // 4. Start server
    log.Infof("Server starting on port %s", port)
    if err := srv.Start(shutdownTimeout); err != nil {
        log.Errorf("Server error: %v", err)
    }
}
```

## Custom Middleware

Create custom middleware following the same pattern:

```go
func TimingMiddleware(log *logger.Logger) server.Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            next.ServeHTTP(w, r)
            duration := time.Since(start)
            
            log.WithFields(map[string]interface{}{
                "method": r.Method,
                "path": r.URL.Path,
                "duration_ms": duration.Milliseconds(),
            }).Info("Request completed")
        })
    }
}

// Use it
srv.Use(TimingMiddleware(log))
```

## Configuration Environment Variables

When using prefix "APP", the config package looks for these environment variables:

```bash
# String values
export APP_PORT=8080
export APP_HOST=localhost

# Integer values
export APP_MAX_CONNECTIONS=100

# Boolean values (true: true/1/yes/on, false: false/0/no/off)
export APP_DEBUG=true
export APP_VERBOSE=1

# Duration values
export APP_TIMEOUT=30s
export APP_IDLE_TIMEOUT=2m
export APP_SHUTDOWN_TIMEOUT=1h

# Required values (panic if not set)
export APP_API_KEY=secret-key-123
```

## Testing

Each package includes comprehensive tests:

```bash
# Test individual packages
go test github.com/Waryway/Wayframe/pkg/config
go test github.com/Waryway/Wayframe/pkg/logger
go test github.com/Waryway/Wayframe/pkg/server

# Test all packages
go test ./...

# Test with coverage
go test -cover ./...

# Verbose test output
go test -v ./...
```

## Bazel Usage

```bash
# Build specific package
bazel build //pkg/config
bazel build //pkg/logger
bazel build //pkg/server

# Test specific package
bazel test //pkg/config:config_test
bazel test //pkg/logger:logger_test
bazel test //pkg/server:server_test

# Build example
bazel build //examples/basic

# Run example
bazel run //examples/basic
```
