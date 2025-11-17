# React Package Guide

The `pkg/react` package provides comprehensive support for serving React applications from Go servers with the Wayframe framework.

## Features

### 1. React Application Serving
- ✅ Static file serving with intelligent caching
- ✅ SPA (Single Page Application) routing support
- ✅ Runtime environment variable injection
- ✅ Content-hash based cache optimization
- ✅ Support for embedded filesystems
- ✅ Compatible with any React build tool (CRA, Vite, Next.js, etc.)

### 2. React Application Bootstrapping
- ✅ Create new React apps with multiple templates
- ✅ Support for Create React App, Vite, and Next.js
- ✅ TypeScript support
- ✅ Interactive and programmatic modes
- ✅ Build and dependency management helpers

## Package Structure

```
pkg/react/
├── doc.go           # Package documentation
├── react.go         # Core serving functionality
├── react_test.go    # Tests for serving
├── bootstrap.go     # React app creation tools
├── bootstrap_test.go # Tests for bootstrap
└── BUILD.bazel      # Bazel build configuration
```

## Usage

### Serving a React Application

```go
package main

import (
    "time"
    "github.com/Waryway/Wayframe/pkg/react"
    "github.com/Waryway/Wayframe/pkg/server"
)

func main() {
    // Create React handler
    reactHandler, err := react.NewHandler(react.Config{
        BuildDir: "./build",           // React build output directory
        BasePath: "/",                 // Base URL path
        EnvVars: map[string]string{    // Runtime environment variables
            "REACT_APP_API_URL": "http://localhost:8080/api",
            "REACT_APP_VERSION": "1.0.0",
            "REACT_APP_ENV":     "production",
        },
        CacheMaxAge: 365 * 24 * time.Hour, // Cache duration for hashed assets
    })
    if err != nil {
        panic(err)
    }

    // Create server and register handler
    srv := server.New(server.Config{Addr: ":8080"})
    
    // Register API routes before React handler
    srv.HandleFunc("/api/health", healthHandler)
    
    // React handler catches all remaining routes
    srv.Handle("/", reactHandler)
    
    srv.Start(30 * time.Second)
}
```

### Environment Variable Injection

The package automatically injects environment variables into your React app by adding a script tag to `index.html`:

```html
<script>
window.__REACT_ENV__ = {
  "REACT_APP_API_URL": "http://localhost:8080/api",
  "REACT_APP_VERSION": "1.0.0",
  "REACT_APP_ENV": "production"
};
</script>
```

Access in your React code:

```javascript
const apiUrl = window.__REACT_ENV__.REACT_APP_API_URL;
const version = window.__REACT_ENV__.REACT_APP_VERSION;
```

### SPA Routing

The handler automatically serves `index.html` for all non-asset routes:

| Request | Response |
|---------|----------|
| `/` | `index.html` |
| `/about` | `index.html` (for client-side routing) |
| `/user/123` | `index.html` (for client-side routing) |
| `/static/js/main.js` | Static JS file with caching |
| `/static/css/style.css` | Static CSS file with caching |
| `/favicon.ico` | Static icon file |
| `/static/missing.js` | 404 error |

### Smart Caching

Files with content hashes (e.g., `main.abc123.js`) get aggressive caching:
```
Cache-Control: public, max-age=31536000, immutable
```

The `index.html` file is never cached:
```
Cache-Control: no-cache, no-store, must-revalidate
```

### Bootstrapping New React Apps

#### Programmatic API

```go
package main

import "github.com/Waryway/Wayframe/pkg/react"

func main() {
    err := react.Bootstrap(react.BootstrapConfig{
        AppName:     "my-app",
        Directory:   "./apps",
        Template:    "vite",        // "cra", "vite", or "next"
        TypeScript:  true,
        SkipInstall: false,
        SkipGit:     false,
        SetupBazel:  true,          // Setup Bazel build files
    })
    if err != nil {
        panic(err)
    }
}
```

#### Building with Bazel

```go
package main

import "github.com/Waryway/Wayframe/pkg/react"

func main() {
    // Build React app with Bazel
    err := react.BuildWithBazel(react.BazelBuildConfig{
        WorkspaceDir: ".",
        Target:       "//apps/my-app:build",
        Config:       "opt",         // Optimized build
        OutputDir:    "./dist",      // Copy output here
    })
    if err != nil {
        panic(err)
    }
}
```

#### Setup Bazel Workspace

```go
// Generate BUILD.bazel and webpack config for existing React app
err := react.SetupReactBazelWorkspace(
    "./my-react-app",  // App directory
    "my-react-app",    // App name
    true,              // TypeScript
)
```

#### Interactive Mode

```go
package main

import "github.com/Waryway/Wayframe/pkg/react"

func main() {
    if err := react.BootstrapInteractive(); err != nil {
        panic(err)
    }
}
```

#### CLI Tool

See `examples/react-bootstrap` for a ready-to-use CLI tool:

```bash
# Interactive mode
bazel run //examples/react-bootstrap -- -interactive

# Create a Vite app with TypeScript
bazel run //examples/react-bootstrap -- -name my-app -template vite -typescript

# Create with Bazel build support
bazel run //examples/react-bootstrap -- -name my-app -template vite -typescript -bazel

# Build an existing app
bazel run //examples/react-bootstrap -- -build ./my-app

# Build with Bazel
bazel run //examples/react-bootstrap -- -build-bazel //apps/my-app:build
```

## Configuration Options

### Handler Configuration

```go
type Config struct {
    BuildDir        string            // React build directory (required if FileSystem not provided)
    BasePath        string            // Base URL path (default: "/")
    EnvVars         map[string]string // Environment variables to inject
    IndexFile       string            // Index file name (default: "index.html")
    CacheMaxAge     time.Duration     // Cache duration (default: 1 year)
    NotFoundHandler http.Handler      // Custom 404 handler (optional)
    FileSystem      fs.FS             // Custom filesystem (optional, for embedding)
}
```

### Bootstrap Configuration

```go
type BootstrapConfig struct {
    AppName     string // React app name (required)
    Directory   string // Creation directory (default: ".")
    Template    string // "cra", "vite", or "next" (default: "cra")
    TypeScript  bool   // Use TypeScript (default: false)
    SkipInstall bool   // Skip npm install (default: false)
    SkipGit     bool   // Skip git init (default: false)
}
```

## Advanced Usage

### Embedded Filesystem

Use Go's `embed` directive for self-contained binaries:

```go
package main

import (
    "embed"
    "github.com/Waryway/Wayframe/pkg/react"
)

//go:embed build/*
var buildFS embed.FS

func main() {
    reactHandler, _ := react.NewHandler(react.Config{
        FileSystem: buildFS,
        BasePath:   "/",
    })
    // ... rest of setup
}
```

### Custom Base Path

Serve the React app at a different path:

```go
reactHandler, _ := react.NewHandler(react.Config{
    BuildDir: "./build",
    BasePath: "/app",  // App served at /app/*
})
```

### Multiple React Apps

Serve different React apps on different paths:

```go
adminHandler, _ := react.NewHandler(react.Config{
    BuildDir: "./admin-build",
    BasePath: "/admin",
})

clientHandler, _ := react.NewHandler(react.Config{
    BuildDir: "./client-build",
    BasePath: "/",
})

srv.Handle("/admin/", adminHandler)
srv.Handle("/", clientHandler)
```

### Custom 404 Handler

```go
reactHandler, _ := react.NewHandler(react.Config{
    BuildDir: "./build",
    NotFoundHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("Custom 404 page"))
    }),
})
```

## Integration with Other Wayframe Packages

### With Config Package

```go
cfg := config.New("APP")
cfg.LoadFile("config.json")

reactHandler, _ := react.NewHandler(react.Config{
    BuildDir: cfg.String("REACT_BUILD_DIR", "./build"),
    BasePath: cfg.String("REACT_BASE_PATH", "/"),
    EnvVars: map[string]string{
        "REACT_APP_API_URL": cfg.String("API_URL", "http://localhost:8080"),
        "REACT_APP_VERSION": cfg.String("APP_VERSION", "1.0.0"),
    },
})
```

### With Logger Package

```go
log := logger.New(logger.InfoLevel)

srv := server.New(server.Config{Addr: ":8080"})
srv.Use(server.LoggingMiddleware(log))

reactHandler, _ := react.NewHandler(react.Config{
    BuildDir: "./build",
})

srv.Handle("/", reactHandler)
```

## Examples

1. **Basic Example**: See `examples/react` for a complete working example
2. **Bootstrap CLI**: See `examples/react-bootstrap` for a CLI tool to create React apps

## Build Tool Compatibility

### Create React App (CRA)
```bash
# Build directory: build/
npm run build
```

### Vite
```bash
# Build directory: dist/
npm run build
```

### Next.js
```bash
# Build directory: out/ (for static export)
# Add to next.config.js: output: 'export'
npm run build
```

### Bazel
```bash
# Build directory: bazel-bin/path/to/target
bazel build --config=opt //apps/my-app:build

# Or use the helper
react.BuildWithBazel(react.BazelBuildConfig{
    Target: "//apps/my-app:build",
    Config: "opt",
})
```

**Why Bazel?**
- Deterministic, reproducible builds
- Fast incremental builds with caching
- Great for monorepos
- Remote build execution support
- Integration with Go builds

See [REACT_BAZEL_GUIDE.md](REACT_BAZEL_GUIDE.md) for detailed Bazel documentation.

## Best Practices

1. **Environment Variables**: Always prefix with `REACT_APP_` for convention
2. **Content Hashing**: Use build tools that generate content-hashed filenames
3. **API Routes**: Register API routes before the React handler
4. **Base Path**: Use base paths when serving multiple apps or in subpaths
5. **Caching**: Let the package handle caching - it's optimized for React builds
6. **Build Directory**: Use relative paths or absolute paths consistently
7. **Error Handling**: Always check errors from `NewHandler`

## Performance Tips

1. **Content Hashing**: Ensure your build tool generates hashed filenames for optimal caching
2. **Compression**: Use reverse proxy (nginx) for gzip/brotli compression
3. **CDN**: Consider serving static assets from a CDN
4. **Embedded FS**: Use `embed.FS` for faster startup and simpler deployment

## Troubleshooting

### "Failed to read index.html"
- Check that BuildDir points to the correct React build output
- Verify the build was successful
- Check file permissions

### "Node.js is required but not found"
- Install Node.js from nodejs.org
- Ensure `node` and `npm` are in your PATH
- Only needed for bootstrap functionality, not for serving

### Environment Variables Not Showing
- Check that variables are in the `EnvVars` map
- Verify the `window.__REACT_ENV__` object in browser console
- Ensure no CSP blocking inline scripts

### SPA Routes Not Working
- Make sure the React handler is registered last
- Check that routes don't conflict with API routes
- Verify client-side routing is configured in your React app

### Static Files Not Loading
- Check browser network tab for 404s
- Verify file paths in React build
- Check that BuildDir is correct
- Ensure files exist in the build directory

## Testing

Run tests with:

```bash
# Go
go test ./pkg/react/...

# Bazel
bazel test //pkg/react:react_test
```

## License

See the main project LICENSE file.

