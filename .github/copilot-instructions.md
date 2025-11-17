# Wayframe Development Instructions

## Technology Stack

### Bazel 8.0.0 with Bzlmod
- **DO NOT use WORKSPACE files** - Bazel 8 uses MODULE.bazel with bzlmod
- Use `MODULE.bazel` for dependency management
- Enable bzlmod in `.bazelrc` with `common --enable_bzlmod`
- Use `bazel_dep()` for external dependencies
- Use extensions for Go SDK and dependencies
- Do not use io_bazel_rules_go, just use rules_go
### Go 1.25
- Use Go 1.25 for all development
- Leverage modern Go features including log/slog
- Follow idiomatic Go practices
- Prefer pkg/logger over fmt.Println for logging

### React | NODE | Javascript
- Use React for frontend applications
- Use pnpm for package management
- Use bazel rules for building React apps
- Embed React builds into Go servers for deployment

### Internal Package Structure
- Prefer public packages under `pkg/` for APIs
- Examples live under `examples/`
- Avoid long-lived framework code under `internal/` unless truly module-private.
- Key packages and locations:
  - `examples/react/env` - Example-scoped environment/config + logger helper
  - `pkg/config` - Configuration loader and tag-based binding
  - `pkg/logger` - Structured logging wrapper (slog)
  - `pkg/server` - Web server utilities and interfaces
  - `pkg/server/stdlib` - Standard library HTTP adapter
  - `pkg/server/fiber` - Fiber adapter
  - `pkg/server/gorilla` - Gorilla Mux adapter

## Configuration Management

- Support multiple file formats: JSON, YAML, key-value pairs
- Use struct tags:
  - `config:"key"` - Configuration key name
  - `env:"ENV_VAR"` - Environment variable name
  - `default:"value"` - Default value
  - `file:"path"` - Configuration file path
- Priority: Environment variables → File values → Default values

## Web Server Interface

All web servers implement a common adapter under `pkg/server`:
- stdlib - Standard library net/http (`pkg/server/stdlib`)
- fiber - Fiber v2 framework (`pkg/server/fiber`)
- gorilla - Gorilla Mux router (`pkg/server/gorilla`)

## Build Commands

```bash
# Bazel
bazel build //...
bazel test //...

# Go
go build ./...
go test ./...
```
