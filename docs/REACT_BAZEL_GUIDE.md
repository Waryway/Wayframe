# React with Bazel Example

This example demonstrates building and serving a React application using Bazel instead of npm/webpack.

## Features

- ✅ Build React app with Bazel
- ✅ TypeScript support
- ✅ Webpack bundling via Bazel
- ✅ Deterministic builds
- ✅ Fast incremental builds
- ✅ Integration with Go server

## Prerequisites

- Bazel 8.0.0+
- Node.js 20+ (for Bazel's Node toolchain)
- Go 1.25+

## Project Structure

```
react-bazel-example/
├── BUILD.bazel          # Bazel build file for React app
├── webpack.config.js    # Webpack configuration
├── package.json         # npm dependencies
├── .bazelrc            # Bazel configuration
├── src/
│   ├── index.tsx       # Entry point
│   ├── App.tsx         # Main component
│   └── index.css       # Styles
├── public/
│   ├── index.html      # HTML template
│   └── favicon.ico     # Favicon
└── server/
    ├── main.go         # Go server
    └── BUILD.bazel     # Go server build file
```

## Quick Start

### 1. Create React App with Bazel Support

```bash
# Using the bootstrap tool
bazel run //examples/react-bootstrap -- \
  -name my-bazel-app \
  -template vite \
  -typescript \
  -bazel

cd my-bazel-app
```

### 2. Build with Bazel

```bash
# Development build
bazel build //:my-bazel-app_build

# Production/optimized build
bazel build --config=opt //:my-bazel-app_build

# Watch mode (with ibazel)
ibazel build //:my-bazel-app_build
```

### 3. Serve with Go

```go
package main

import (
    "github.com/Waryway/Wayframe/pkg/react"
    "github.com/Waryway/Wayframe/pkg/server"
    "time"
)

func main() {
    // Build React app with Bazel
    react.BuildWithBazel(react.BazelBuildConfig{
        Target: "//:my-bazel-app_build",
        Config: "opt",
        OutputDir: "./dist",
    })
    
    // Serve with Wayframe
    reactHandler, _ := react.NewHandler(react.Config{
        BuildDir: "./dist",
        EnvVars: map[string]string{
            "REACT_APP_API_URL": "http://localhost:8080/api",
        },
    })
    
    srv := server.New(server.Config{Addr: ":8080"})
    srv.Handle("/", reactHandler)
    srv.Start(30 * time.Second)
}
```

## Bazel Configuration

### BUILD.bazel

The generated BUILD.bazel file uses Aspect Rules for JavaScript:

```python
load("@aspect_rules_js//js:defs.bzl", "js_library")
load("@aspect_rules_webpack//webpack:defs.bzl", "webpack_bundle")
load("@aspect_rules_ts//ts:defs.bzl", "ts_project")

# TypeScript compilation
ts_project(
    name = "ts",
    srcs = glob(["src/**/*.ts", "src/**/*.tsx"]),
    declaration = True,
    tsconfig = "tsconfig.json",
)

# Webpack bundle
webpack_bundle(
    name = "build",
    srcs = [
        ":ts",
        "package.json",
        "webpack.config.js",
    ] + glob(["public/**/*"]),
    output_dir = True,
    webpack_config = "webpack.config.js",
)

# Export for Go server
filegroup(
    name = "my-bazel-app_build",
    srcs = [":build"],
    visibility = ["//visibility:public"],
)
```

### .bazelrc

```bash
# React build configuration
build --strategy=Webpack=worker
build --worker_max_instances=4

# Production build
build:opt --compilation_mode=opt
build:opt --define=NODE_ENV=production

# Development build
build:dev --compilation_mode=fastbuild
build:dev --define=NODE_ENV=development
```

## Benefits of Bazel for React

### 1. Deterministic Builds
- Same inputs always produce same outputs
- No "works on my machine" issues
- Reproducible across environments

### 2. Fast Incremental Builds
- Only rebuilds what changed
- Shared cache across team
- Remote build execution support

### 3. Monorepo Support
- Build multiple React apps efficiently
- Share components across apps
- Consistent tooling

### 4. Integration with Go
- Build React and Go in same workflow
- Single build system
- Easy CI/CD setup

## Advanced Usage

### Custom Webpack Config

Modify `webpack.config.js` for your needs:

```javascript
module.exports = {
  entry: './src/index.tsx',
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'static/js/[name].[contenthash:8].js',
  },
  // ... your custom config
};
```

### Multiple React Apps

```python
# apps/admin/BUILD.bazel
webpack_bundle(
    name = "admin_build",
    srcs = [...],
)

# apps/client/BUILD.bazel
webpack_bundle(
    name = "client_build",
    srcs = [...],
)

# Serve both in Go
```

### Shared Components

```python
# components/BUILD.bazel
ts_project(
    name = "common",
    srcs = glob(["*.tsx"]),
    visibility = ["//apps:__subpackages__"],
)

# apps/my-app/BUILD.bazel
webpack_bundle(
    name = "build",
    deps = [
        "//components:common",
    ],
)
```

### Remote Caching

Enable remote build cache in `.bazelrc`:

```bash
build --remote_cache=https://your-cache-server
build --remote_upload_local_results=true
```

## Build Commands

```bash
# Build React app
bazel build //:my-app_build

# Build optimized
bazel build --config=opt //:my-app_build

# Build with remote cache
bazel build --remote_cache=grpc://localhost:9092 //:my-app_build

# Clean
bazel clean

# Test
bazel test //...

# Build everything
bazel build //...
```

## Integration with CI/CD

### GitHub Actions

```yaml
name: Build React with Bazel

on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: bazel-contrib/setup-bazel@v1
        with:
          bazelisk-version: '1.x'
      
      - name: Build React App
        run: bazel build --config=opt //:my-app_build
      
      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: react-build
          path: bazel-bin/build
```

### With Remote Cache

```yaml
      - name: Setup Bazel cache
        uses: actions/cache@v3
        with:
          path: ~/.cache/bazel
          key: bazel-${{ hashFiles('.bazelversion', 'MODULE.bazel') }}
      
      - name: Build with cache
        run: |
          bazel build \
            --remote_cache=https://cache.example.com \
            --config=opt \
            //:my-app_build
```

## Troubleshooting

### "Cannot find webpack"

Ensure webpack is listed in `deps` in BUILD.bazel:

```python
deps = [
    "//:node_modules/webpack",
    "//:node_modules/webpack-cli",
]
```

### Slow builds

1. Enable worker strategy: `--strategy=Webpack=worker`
2. Increase workers: `--worker_max_instances=8`
3. Use remote cache
4. Use `ibazel` for watch mode

### TypeScript errors

Check your `tsconfig.json` is valid and referenced in BUILD.bazel:

```python
ts_project(
    name = "ts",
    tsconfig = "tsconfig.json",
)
```

## Comparison: npm vs Bazel

| Feature | npm/webpack | Bazel |
|---------|-------------|-------|
| Build speed | Moderate | Fast (incremental) |
| Caching | Local only | Local + Remote |
| Reproducibility | Variable | Deterministic |
| Monorepo | Complex | Native |
| Learning curve | Low | Medium |
| Tooling | Mature | Growing |

## When to Use Bazel

**Use Bazel when:**
- ✅ You have a monorepo with multiple apps
- ✅ You need deterministic, reproducible builds
- ✅ You want to share build cache across team
- ✅ You're already using Bazel for other parts of your stack
- ✅ You need to build React + Go together

**Stick with npm when:**
- ❌ You have a single, standalone React app
- ❌ Your team isn't familiar with Bazel
- ❌ You rely on npm-specific tooling

## Resources

- [Aspect Rules JS](https://github.com/aspect-build/rules_js)
- [Bazel Documentation](https://bazel.build/docs)
- [Rules Webpack](https://github.com/aspect-build/rules_webpack)

## License

See the main project LICENSE file.

