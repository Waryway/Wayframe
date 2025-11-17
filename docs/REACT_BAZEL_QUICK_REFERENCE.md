# Quick Reference: React with Bazel in Wayframe

## Setup New React App with Bazel

### Interactive
```bash
bazel run //examples/react-bootstrap -- -interactive
# Answer "y" to "Setup Bazel build files?"
```

### Command Line
```bash
bazel run //examples/react-bootstrap -- \
  -name my-app \
  -template vite \
  -typescript \
  -bazel
```

### Programmatic
```go
import "github.com/Waryway/Wayframe/pkg/react"

react.Bootstrap(react.BootstrapConfig{
    AppName:    "my-app",
    Template:   "vite",
    TypeScript: true,
    SetupBazel: true,
})
```

## Add Bazel to Existing React App

```go
react.SetupReactBazelWorkspace(
    "./my-react-app",
    "my-react-app",
    true, // TypeScript
)
```

This creates:
- `BUILD.bazel` - Bazel build configuration
- `webpack.config.js` - Webpack configuration
- `.bazelrc` - Bazel settings

## Build React App with Bazel

### From Command Line
```bash
# Development build
bazel build //apps/my-app:build

# Production build
bazel build --config=opt //apps/my-app:build

# Specific output directory
bazel run //examples/react-bootstrap -- \
  -build-bazel //apps/my-app:build
```

### From Go Code
```go
react.BuildWithBazel(react.BazelBuildConfig{
    WorkspaceDir: ".",
    Target:       "//apps/my-app:build",
    Config:       "opt",
    OutputDir:    "./dist",
})
```

## Serve Bazel-Built React App

```go
package main

import (
    "time"
    "github.com/Waryway/Wayframe/pkg/react"
    "github.com/Waryway/Wayframe/pkg/server"
)

func main() {
    // Build with Bazel
    react.BuildWithBazel(react.BazelBuildConfig{
        Target:    "//apps/my-app:build",
        Config:    "opt",
        OutputDir: "./dist",
    })
    
    // Serve
    handler, _ := react.NewHandler(react.Config{
        BuildDir: "./dist",
        EnvVars: map[string]string{
            "REACT_APP_API_URL": "http://localhost:8080/api",
        },
    })
    
    srv := server.New(server.Config{Addr: ":8080"})
    srv.Handle("/", handler)
    srv.Start(30 * time.Second)
}
```

## Common Bazel Commands

```bash
# Build
bazel build //apps/my-app:build

# Build optimized
bazel build --config=opt //apps/my-app:build

# Build for dev
bazel build --config=dev //apps/my-app:build

# Clean
bazel clean

# Clean everything
bazel clean --expunge

# Test
bazel test //...

# Run
bazel run //apps/my-app:serve
```

## Bazel Configuration Files

### MODULE.bazel (Workspace Root)
```python
bazel_dep(name = "aspect_rules_js", version = "2.1.0")
bazel_dep(name = "aspect_rules_ts", version = "3.2.1")
bazel_dep(name = "aspect_rules_webpack", version = "0.18.0")
```

### BUILD.bazel (App Directory)
```python
load("@aspect_rules_ts//ts:defs.bzl", "ts_project")
load("@aspect_rules_webpack//webpack:defs.bzl", "webpack_bundle")

ts_project(
    name = "ts",
    srcs = glob(["src/**/*.ts", "src/**/*.tsx"]),
    tsconfig = "tsconfig.json",
)

webpack_bundle(
    name = "build",
    srcs = [":ts"],
    output_dir = True,
    webpack_config = "webpack.config.js",
)
```

### .bazelrc (App Directory)
```bash
build --strategy=Webpack=worker
build --worker_max_instances=4

build:opt --compilation_mode=opt
build:opt --define=NODE_ENV=production

build:dev --compilation_mode=fastbuild
build:dev --define=NODE_ENV=development
```

## Directory Structure

```
my-app/
├── BUILD.bazel              # Generated
├── .bazelrc                 # Generated
├── webpack.config.js        # Generated
├── package.json            # From template
├── tsconfig.json           # From template (if TypeScript)
├── src/
│   ├── index.tsx
│   ├── App.tsx
│   └── ...
└── public/
    ├── index.html
    └── favicon.ico
```

## Build Output

```
bazel-bin/
└── apps/
    └── my-app/
        └── build/           # React build output
            ├── index.html
            └── static/
                ├── js/
                └── css/
```

## Comparison

| Feature | npm Build | Bazel Build |
|---------|-----------|-------------|
| Command | `npm run build` | `bazel build //apps/my-app:build` |
| Output | `build/` or `dist/` | `bazel-bin/apps/my-app/build/` |
| Caching | Local | Local + Remote |
| Speed | Moderate | Fast (incremental) |
| Reproducible | No | Yes |
| Monorepo | Complex | Native |

## Tips

1. **Use --config=opt for production**
   ```bash
   bazel build --config=opt //apps/my-app:build
   ```

2. **Enable remote cache for team builds**
   ```bash
   build --remote_cache=grpc://cache.example.com
   ```

3. **Use ibazel for watch mode**
   ```bash
   ibazel build //apps/my-app:build
   ```

4. **Copy output for serving**
   ```bash
   cp -r bazel-bin/apps/my-app/build ./dist
   ```

5. **Clean when switching branches**
   ```bash
   bazel clean
   ```

## Troubleshooting

### Build fails with "Cannot find module"
Check `deps` in BUILD.bazel includes all npm packages.

### Slow builds
- Use `--worker_max_instances=8`
- Enable remote cache
- Use `ibazel` for incremental builds

### Output not found
Bazel output is in `bazel-bin/path/to/target/build/`

## Resources

- Full Guide: `docs/REACT_BAZEL_GUIDE.md`
- Package Docs: `docs/REACT_PACKAGE_GUIDE.md`
- Examples: `examples/react-bootstrap/`

## Need Help?

```bash
# Show help
bazel help

# CLI tool help
bazel run //examples/react-bootstrap -- -h

# Check Bazel version
bazel version
```

