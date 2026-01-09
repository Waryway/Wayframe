#!/bin/bash
cd /c/Users/kawie/GolandProjects/Wayframe

echo "========================================" > /tmp/build-test.log
echo "Build Test - $(date)" >> /tmp/build-test.log
echo "========================================" >> /tmp/build-test.log
echo >> /tmp/build-test.log

echo "Running: bazel build //examples/react/frontend:build_react" >> /tmp/build-test.log
echo >> /tmp/build-test.log

bazel build //examples/react/frontend:build_react >> /tmp/build-test.log 2>&1
EXIT_CODE=$?

echo >> /tmp/build-test.log
echo "========================================" >> /tmp/build-test.log
echo "Exit code: $EXIT_CODE" >> /tmp/build-test.log
echo "========================================" >> /tmp/build-test.log

if [ $EXIT_CODE -eq 0 ]; then
    echo >> /tmp/build-test.log
    echo "✅ BUILD SUCCESS" >> /tmp/build-test.log
    echo >> /tmp/build-test.log
    echo "Output files:" >> /tmp/build-test.log
    ls -la bazel-bin/examples/react/frontend/dist/ >> /tmp/build-test.log 2>&1
else
    echo >> /tmp/build-test.log
    echo "❌ BUILD FAILED" >> /tmp/build-test.log
fi

cat /tmp/build-test.log

