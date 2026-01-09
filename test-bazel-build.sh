#!/bin/bash
set -e

echo "==================================="
echo "Testing Bazel React Build"
echo "==================================="
echo

cd /c/Users/kawie/GolandProjects/Wayframe

echo "Step 1: Clean build"
bazel clean

echo
echo "Step 2: Build frontend only"
bazel build //examples/react/frontend:build_react

if [ $? -eq 0 ]; then
    echo "✅ Frontend build SUCCESS"
    echo
    echo "Checking outputs:"
    ls -la bazel-bin/examples/react/frontend/dist/
else
    echo "❌ Frontend build FAILED"
    exit 1
fi

echo
echo "Step 3: Build everything"
bazel build //...

if [ $? -eq 0 ]; then
    echo "✅ Full build SUCCESS"
else
    echo "❌ Full build FAILED"
    exit 1
fi

echo
echo "==================================="
echo "✅ ALL BUILDS PASSED"
echo "==================================="

