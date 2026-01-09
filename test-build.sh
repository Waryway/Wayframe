#!/bin/bash
set -euo pipefail

echo "Testing Bazel build with custom React rule..."
echo

cd /c/Users/kawie/GolandProjects/Wayframe

echo "1. Cleaning previous build..."
bazel clean

echo
echo "2. Building frontend only..."
if bazel build //examples/react/frontend:build_react; then
    echo "✅ Frontend build succeeded"
else
    echo "❌ Frontend build failed"
    exit 1
fi

echo
echo "3. Building full React example..."
if bazel build //examples/react:react; then
    echo "✅ Full React example build succeeded"
else
    echo "❌ Full React example build failed"
    exit 1
fi

echo
echo "4. Building all targets..."
if bazel build //...; then
    echo "✅ All targets built successfully"
else
    echo "❌ Some targets failed to build"
    exit 1
fi

echo
echo "✅ All builds completed successfully!"

