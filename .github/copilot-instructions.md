# Wayframe Development Instructions

## Technology Stack

### Bazel 8.0.0 with Bzlmod
- **DO NOT use WORKSPACE files** - Bazel 8 uses MODULE.bazel with bzlmod
- Use `MODULE.bazel` for dependency management
- Enable bzlmod in `.bazelrc` with `common --enable_bzlmod`
- Use `bazel_dep()` for external dependencies
- Use extensions for Go SDK and dependencies

### Go 1.25
- Use Go 1.25 for all development
- Leverage modern Go features including log/slog
- Follow idiomatic Go practices

### Internal Package Structure
- Use `internal/` directory for framework internals
- Packages in `internal/` are not importable by external projects
- Key internal packages:
  - `internal/config` - Configuration management
  - `internal/env` - Environment initialization
  - `internal/web` - Web server abstractions
  - `internal/web/stdlib` - Standard library HTTP server
  - `internal/web/fiber` - Fiber web framework adapter
  - `internal/web/gorilla` - Gorilla Mux adapter

## Configuration Management

- Support multiple file formats: JSON, YAML, key-value pairs
- Use struct tags:
  - `config:"key"` - Configuration key name
  - `env:"ENV_VAR"` - Environment variable name
  - `default:"value"` - Default value
  - `file:"path"` - Configuration file path
- Priority: Environment variables → File values → Default values

## Web Server Interface

All web servers implement `internal/web.Server` interface:
- stdlib - Standard library net/http
- fiber - Fiber v2 framework
- gorilla - Gorilla Mux router

## Build Commands

```bash
# Bazel
bazel build //...
bazel test //...

# Go
go build ./...
go test ./...
```
