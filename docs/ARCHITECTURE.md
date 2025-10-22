# Wayframe Architecture

This document describes the architecture and design decisions of the Wayframe framework.

## Overview

Wayframe is a modular Go framework built with Bazel. It consists of three core packages that can be used independently or together:

1. **config**: Configuration management
2. **logger**: Structured logging
3. **server**: HTTP server utilities

## Design Principles

### 1. Package Independence

Each package is designed to work standalone without requiring other Wayframe packages. This allows developers to:
- Use only the packages they need
- Avoid unnecessary dependencies
- Integrate Wayframe packages into existing projects

### 2. Idiomatic Go

Wayframe follows Go best practices:
- Standard library patterns (http.Handler, context.Context)
- Interface-based design where appropriate
- Clear error handling
- No unnecessary abstractions

### 3. Pragmatic Defaults

Packages provide sensible defaults while allowing customization:
- config: Default values for all configuration options
- logger: Default log level (Info), default output (stdout)
- server: Standard timeouts, graceful shutdown

### 4. Build System Flexibility

Wayframe supports both build systems:
- **Go Modules**: Standard Go tooling (`go build`, `go test`)
- **Bazel**: Reproducible builds with explicit dependencies

## Package Design

### Config Package

**Purpose**: Load configuration from environment variables with type safety.

**Key Features**:
- Prefix support for namespacing
- Type conversion (string, int, bool, duration)
- Default values
- Required value validation

**Design Decisions**:
- No file parsing (environment variables only) - keeps it simple
- Panic on missing required values - fail fast principle
- Immutable loader - configuration shouldn't change at runtime

### Logger Package

**Purpose**: Provide structured logging with levels and contextual fields.

**Key Features**:
- Standard log levels (Debug, Info, Warn, Error)
- Contextual fields (WithField, WithFields)
- Formatted logging (Infof, Errorf, etc.)
- Thread-safe operations

**Design Decisions**:
- Immutable field additions - prevents accidental state sharing
- Structured output format - easy to parse by log aggregators
- No external dependencies - uses only standard library
- Simple implementation - ~150 lines of code

### Server Package

**Purpose**: HTTP server with graceful shutdown and middleware support.

**Key Features**:
- Graceful shutdown with timeout
- Middleware chain support
- Built-in logging and recovery middleware
- Signal handling (SIGINT, SIGTERM)

**Design Decisions**:
- Wraps http.Server - doesn't reinvent the wheel
- Standard http.Handler interface - works with any HTTP library
- Middleware applied in registration order - intuitive behavior
- Blocking Start method - simplifies application lifecycle

## Module Structure

```
github.com/Waryway/Wayframe
├── pkg/                    # Core packages
│   ├── config/            # Configuration management
│   │   ├── config.go      # Main implementation
│   │   ├── config_test.go # Unit tests
│   │   ├── doc.go         # Package documentation
│   │   └── BUILD.bazel    # Bazel build file
│   ├── logger/            # Structured logging
│   │   ├── logger.go      # Main implementation
│   │   ├── logger_test.go # Unit tests
│   │   ├── doc.go         # Package documentation
│   │   └── BUILD.bazel    # Bazel build file
│   └── server/            # HTTP server utilities
│       ├── server.go      # Main implementation
│       ├── server_test.go # Unit tests
│       ├── doc.go         # Package documentation
│       └── BUILD.bazel    # Bazel build file
├── examples/              # Example applications
│   └── basic/             # Basic usage example
│       ├── main.go        # Example implementation
│       └── BUILD.bazel    # Bazel build file
├── docs/                  # Documentation
│   ├── ARCHITECTURE.md    # This file
│   ├── USAGE.md          # Usage examples
│   └── CONTRIBUTING.md   # Contribution guidelines
├── BUILD.bazel           # Root build file
├── WORKSPACE             # Bazel workspace
├── deps.bzl             # External dependencies
├── go.mod               # Go module definition
└── README.md            # Project overview
```

## Testing Strategy

Each package includes comprehensive unit tests:

- **config**: Tests all type conversions, prefix handling, defaults
- **logger**: Tests log levels, fields, formatting, output
- **server**: Tests middleware, graceful shutdown, routing

Test coverage: >80% for all packages

## Dependencies

Wayframe has zero external dependencies. It uses only the Go standard library:
- `os` - Environment variable access
- `time` - Duration parsing and timing
- `net/http` - HTTP server functionality
- `context` - Request context and cancellation
- `log` - Basic logging primitives

## Future Considerations

Potential additions while maintaining the core principles:

1. **Database Package**: Connection pooling and query helpers
2. **Metrics Package**: Application metrics and monitoring
3. **Cache Package**: In-memory and distributed caching
4. **Validation Package**: Input validation utilities

Each would follow the same principles:
- Independent and modular
- Zero external dependencies (or minimal)
- Idiomatic Go
- Well-tested and documented

## Performance Characteristics

### Config Package
- O(1) environment variable lookups
- No runtime allocations after initialization
- Negligible overhead

### Logger Package
- Thread-safe with minimal lock contention
- Structured format with low serialization overhead
- Suitable for high-throughput applications

### Server Package
- Standard net/http performance
- Minimal middleware overhead
- Efficient graceful shutdown

## Security Considerations

1. **Config**: Environment variables may contain secrets - handle appropriately
2. **Logger**: Avoid logging sensitive data (passwords, tokens)
3. **Server**: Uses standard http.Server security features
4. **Dependencies**: Zero external dependencies = minimal supply chain risk

## Versioning

Wayframe follows semantic versioning (SemVer):
- MAJOR: Breaking API changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes (backward compatible)

Current version: 0.1.0 (initial release)
