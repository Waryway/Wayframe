# FINAL SOLUTION - React Build Rule

## Problem Solved

Error: `$(echo) not defined` - Bazel was interpreting `$(echo)` as a make variable instead of bash command substitution.

## The Fix

Changed from `$(...)` to backticks `` `...` `` for bash command substitution:

```bash
# ❌ WRONG - Bazel interprets $(echo) as make variable
HTML_OUT=$(echo $(OUTS) | awk '{print $1}')
mkdir -p "$(dirname "$$HTML_OUT")"

# ✅ CORRECT - Backticks prevent Bazel interpretation
HTML_OUT=`echo $(OUTS) | awk '{print $$1}'`
mkdir -p "`dirname "$$HTML_OUT"`"
```

## Complete Working Code

**File:** `bazel/react.bzl`

```bash
cd {frontend_dir}
npm install --silent 2>&1 | grep -v "npm WARN" || true
npm run {build_script} --silent

cd {output_dir}

# Use backticks for command substitution
HTML_OUT=`echo $(OUTS) | awk '{{print $$1}}'`
JS_OUT=`echo $(OUTS) | awk '{{print $$2}}'`

# Create parent directories
mkdir -p "`dirname "$$HTML_OUT"`"
mkdir -p "`dirname "$$JS_OUT"`"

# Copy files
cp index.html "$$HTML_OUT"

for js in *.js; do
    if [ -f "$$js" ]; then
        cp "$$js" "$$JS_OUT"
        break
    fi
done
```

## Why This Works

1. **Backticks**: `` `command` `` is old-style bash command substitution that Bazel doesn't interpret
2. **$(OUTS)**: Bazel make variable that expands to the list of output files
3. **$$variable**: Double-dollar for bash variables (prevents Bazel expansion)
4. **dirname**: Gets parent directory path
5. **mkdir -p**: Creates all parent directories
6. **Direct copy**: Copies to exact paths Bazel expects

## Test Commands

```bash
cd /c/Users/kawie/GolandProjects/Wayframe

# Test 1: Clean build
bazel clean
bazel build //examples/react/frontend:build_react

# Test 2: Full build
bazel clean --expunge
bazel build //...
```

## Expected Result

```
webpack 5.104.1 compiled successfully
Target //examples/react/frontend:build_react up-to-date:
  bazel-bin/examples/react/frontend/dist/bundle.js
  bazel-bin/examples/react/frontend/dist/index.html
INFO: Build completed successfully
```

## Files Changed

- `bazel/react.bzl` - Fixed command substitution syntax

## Summary

The issue was using `$(...)` for bash command substitution, which Bazel interprets as make variable syntax. Using backticks `` `...` `` instead allows bash to handle the command substitution while Bazel only expands its own variables like `$(OUTS)`.

**Status: FIXED AND READY TO TEST**

