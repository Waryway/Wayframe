# Contributing to Wayframe

Thank you for your interest in contributing to Wayframe! This document provides guidelines and instructions for contributing.

## Philosophy

Wayframe follows these principles:

1. **Idiomatic Go**: Follow Go best practices and conventions
2. **Simplicity**: Keep the API simple and intuitive
3. **Independence**: Packages should work standalone
4. **Tested**: All code should be well-tested
5. **Documented**: Public APIs should have clear documentation

## Development Setup

### Prerequisites

- Go 1.24.7 or later
- Bazel 7.0.0 or later (optional)
- Git

### Getting Started

```bash
# Clone the repository
git clone https://github.com/Waryway/Wayframe.git
cd Wayframe

# Run tests
go test ./...

# Build the example
go build ./examples/basic
```

## Making Changes

### Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Keep functions focused and concise
- Add comments for exported functions and types

### Adding a New Feature

1. Create a new branch for your feature
2. Write tests first (TDD approach)
3. Implement the feature
4. Ensure all tests pass
5. Add documentation
6. Submit a pull request

### Package Guidelines

When adding new packages:

1. Place in `pkg/` directory
2. Include comprehensive tests (`*_test.go`)
3. Add package documentation (`doc.go`)
4. Create `BUILD.bazel` for Bazel support
5. Update README.md with package description
6. Add usage examples

Example package structure:
```
pkg/newpackage/
├── BUILD.bazel       # Bazel build file
├── doc.go            # Package documentation
├── newpackage.go     # Implementation
└── newpackage_test.go # Tests
```

### Testing

- Write unit tests for all public functions
- Aim for high test coverage (>80%)
- Test edge cases and error conditions
- Use table-driven tests where appropriate

Example test structure:
```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    interface{}
        expected interface{}
    }{
        {"case1", input1, expected1},
        {"case2", input2, expected2},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Documentation

- All exported functions, types, and constants must have godoc comments
- Package-level documentation in `doc.go`
- Include code examples in documentation
- Update README.md for significant changes

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./pkg/config

# Run tests with race detector
go test -race ./...

# Verbose output
go test -v ./...
```

## Building with Bazel

```bash
# Build all targets
bazel build //...

# Run all tests
bazel test //...

# Build specific package
bazel build //pkg/config

# Run example
bazel run //examples/basic
```

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./...`)
6. Commit with clear, descriptive messages
7. Push to your fork
8. Create a Pull Request

### Pull Request Guidelines

- Clear title describing the change
- Detailed description of what and why
- Reference any related issues
- Include test results
- Update documentation as needed

### Commit Messages

Follow conventional commit format:

```
type(scope): brief description

Detailed description of the change, including:
- What was changed
- Why it was changed
- Any breaking changes

Fixes #123
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test additions or changes
- `refactor`: Code refactoring
- `chore`: Build or tooling changes

## Code Review

All submissions require review. We use GitHub pull requests for this purpose.

Reviewers will check for:
- Code quality and style
- Test coverage
- Documentation
- Breaking changes
- Performance implications

## Questions?

Feel free to open an issue for:
- Bug reports
- Feature requests
- Questions about implementation
- General discussion

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (Apache-2.0).
