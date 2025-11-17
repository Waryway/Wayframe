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

### react

Serve React applications from Go servers with SPA routing, environment variable injection, and smart caching.

```go
import "github.com/Waryway/Wayframe/pkg/react"

reactHandler, err := react.NewHandler(react.Config{
    BuildDir: "./build",
    BasePath: "/",
    EnvVars: map[string]string{
        "REACT_APP_API_URL": "http://localhost:8080/api",
        "REACT_APP_VERSION": "1.0.0",
    },
})

srv := server.New(server.Config{Addr: ":8080"})
srv.Handle("/", reactHandler)
srv.Start(30 * time.Second)
```

**Features**:
- Static asset serving with proper cache headers
- SPA fallback routing (all non-asset routes serve index.html)
- Runtime environment variable injection into React builds
- Content-hash based cache optimization
- Compatible with embedded filesystems
- **Build React apps with npm or Bazel**
- Bootstrap new React apps (CRA, Vite, Next.js)
- Bazel build file generation and management

## Example

See [examples/basic](examples/basic/main.go) for a complete example demonstrating all packages working together.

See [examples/react](examples/react/README.md) for a React SPA serving example with environment variable injection and API routes.

See [docs/REACT_BAZEL_GUIDE.md](docs/REACT_BAZEL_GUIDE.md) for building React applications with Bazel.

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
