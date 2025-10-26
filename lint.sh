#!/bin/bash
# Bazel wrapper for golangci-lint
# This script runs golangci-lint on all Go code in the repository

set -e

# Get the workspace directory (Bazel sets this)
if [ -n "$BUILD_WORKSPACE_DIRECTORY" ]; then
  cd "$BUILD_WORKSPACE_DIRECTORY"
fi

# Check if golangci-lint is available
if ! command -v golangci-lint &> /dev/null; then
  echo "Error: golangci-lint is not installed"
  echo "Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
  echo "Or download from: https://golangci-lint.run/usage/install/"
  exit 1
fi

# Run golangci-lint
echo "Running golangci-lint..."
golangci-lint run --timeout=5m "$@"
echo "Linting completed successfully!"
