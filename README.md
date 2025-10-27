# Wayframe

A core framework for Golang applications built with Bazel. Wayframe provides pragmatic, idiomatic, and opinionated packages that can be used independently or together to build robust Go applications.

## Philosophy

While the Go community tends to avoid frameworks, Wayframe exists to prevent reinventing the wheel. It provides:

- **Modular packages**: Each package is independently usable
- **Idiomatic Go**: Follows Go best practices and conventions
- **Pragmatic defaults**: Sensible defaults that work out of the box
- **Bazel integration**: Built with Bazel for reproducible builds

## Packages

### config

Configuration management with support for JSON files, environment variables, and defaults. Environment variable names match configuration keys.

```go
import "github.com/Waryway/Wayframe/pkg/config"

cfg := config.New("APP")

// Optionally load from file (falls back to env vars and defaults)
cfg.LoadFile("config.json")

port := cfg.String("PORT", "8080")
timeout := cfg.Duration("TIMEOUT", 30*time.Second)
debug := cfg.Bool("DEBUG", false)
```

**Priority order**: Environment variables → File values → Default values

### logger

Structured logging based on Go's standard `log/slog` package with a simplified interface.

```go
import "github.com/Waryway/Wayframe/pkg/logger"

log := logger.New(logger.InfoLevel)
log.Info("Application started")
log.WithField("user", "john").Info("User logged in")
log.Errorf("Failed to connect: %v", err)
```

### server

HTTP server with graceful shutdown, middleware support, and common patterns.

```go
import "github.com/Waryway/Wayframe/pkg/server"

srv := server.New(server.Config{
    Addr:         ":8080",
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
})

srv.Use(server.LoggingMiddleware(log))
srv.HandleFunc("/", handler)
srv.Start(30 * time.Second)
```

## Example

See [examples/basic](examples/basic/main.go) for a complete example demonstrating all packages working together.

## Building with Go

```bash
# Run tests
go test ./...

# Build example
go build ./examples/basic

# Run example
./basic
```

## Building with Bazel

```bash
# Build all packages
bazel build //...

# Run tests
bazel test //...

# Build and run example
bazel run //examples/basic
```

## Project Structure

```
.
├── pkg/              # Core packages
│   ├── config/       # Configuration management
│   ├── logger/       # Structured logging
│   └── server/       # HTTP server utilities
├── examples/         # Example applications
│   └── basic/        # Basic usage example
├── BUILD.bazel       # Root Bazel build file
├── WORKSPACE         # Bazel workspace configuration
└── go.mod            # Go module definition
```

## License

See [LICENSE](LICENSE) for details.
