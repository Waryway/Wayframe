#!/bin/bash

echo "Testing React build with race condition fix..."
cd /c/Users/kawie/GolandProjects/Wayframe

echo "1. Cleaning..."
bazel clean > /dev/null 2>&1

echo "2. Building frontend..."
if bazel build //examples/react/frontend:build_react 2>&1 | tee /tmp/build-output.txt | grep -q "Build completed successfully"; then
    echo "✅ BUILD SUCCESS"
    echo
    echo "Checking output files:"
    ls -la bazel-bin/examples/react/frontend/dist/
    exit 0
else
    echo "❌ BUILD FAILED"
    echo
    echo "Last 20 lines of output:"
    tail -20 /tmp/build-output.txt
    exit 1
fi

