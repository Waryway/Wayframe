#!/bin/bash
echo "Testing React build rule fix..."
echo

cd /c/Users/kawie/GolandProjects/Wayframe

echo "Building frontend..."
bazel build //examples/react/frontend:build_react

if [ $? -eq 0 ]; then
    echo "✅ Frontend build SUCCESS"
    echo
    echo "Checking output files..."
    ls -la bazel-bin/examples/react/frontend/dist/
else
    echo "❌ Frontend build FAILED"
    exit 1
fi

