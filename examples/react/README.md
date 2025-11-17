# React Example

This example demonstrates how to serve a React application using Wayframe's `pkg/react` package.

## Features

- ✅ Serve React SPA from Go server
- ✅ Environment variable injection at runtime
- ✅ Smart caching for static assets
- ✅ SPA routing support (all non-asset routes serve index.html)
- ✅ API endpoints alongside React app
- ✅ Graceful shutdown
- ✅ Structured logging

## Quick Start

### Using Go

```bash
# From the examples/react directory
go run main.go
```

### Using Bazel

```bash
# From the project root
bazel run //examples/react

# Or build and run separately
bazel build //examples/react
./bazel-bin/examples/react/react_/react
```

## Configuration

The example can be configured via environment variables or a JSON config file:

```json
{
  "PORT": "8080",
  "LOG_LEVEL": "INFO",
  "REACT_BUILD_DIR": "./build",
  "REACT_BASE_PATH": "/",
  "API_URL": "http://localhost:8080",
  "APP_VERSION": "1.0.0",
  "APP_ENV": "development"
}
```

### Configuration Options

- `PORT` - Server port (default: 8080)
- `LOG_LEVEL` - Log level: DEBUG, INFO (default: INFO)
- `REACT_BUILD_DIR` - Path to React build directory (default: ./build)
- `REACT_BASE_PATH` - Base URL path for React app (default: /)
- `API_URL` - API URL to inject into React app
- `APP_VERSION` - Application version
- `APP_ENV` - Environment: development, production

## How It Works

### 1. React Handler Setup

```go
reactHandler, err := react.NewHandler(react.Config{
    BuildDir: "./build",
    BasePath: "/",
    EnvVars: map[string]string{
        "REACT_APP_API_URL": "http://localhost:8080",
        "REACT_APP_VERSION": "1.0.0",
        "REACT_APP_ENV":     "production",
    },
})
```

### 2. Environment Variable Injection

The Go server injects environment variables into the React app by adding a script tag to `index.html`:

```html
<script>
window.__REACT_ENV__ = {
  "REACT_APP_API_URL": "http://localhost:8080",
  "REACT_APP_VERSION": "1.0.0",
  "REACT_APP_ENV": "production"
};
</script>
```

Your React code can access these variables:

```javascript
const apiUrl = window.__REACT_ENV__.REACT_APP_API_URL;
const version = window.__REACT_ENV__.REACT_APP_VERSION;
```

### 3. SPA Routing

The handler automatically serves `index.html` for all non-asset routes:

- `/` → index.html
- `/about` → index.html (for client-side routing)
- `/user/123` → index.html (for client-side routing)
- `/static/js/main.js` → static file
- `/static/css/style.css` → static file

### 4. Smart Caching

Files with content hashes (e.g., `main.abc123.js`) are cached with:
```
Cache-Control: public, max-age=31536000, immutable
```

Index.html is never cached:
```
Cache-Control: no-cache, no-store, must-revalidate
```

### 5. API Routes

API routes are registered before the React handler:

```go
srv.HandleFunc("/api/health", healthHandler)
srv.HandleFunc("/api/info", infoHandler)
srv.Handle("/", reactHandler)  // Catches all remaining routes
```

## Project Structure

```
examples/react/
├── main.go                    # Go server
├── config.example.json        # Example configuration
├── BUILD.bazel               # Bazel build file
├── README.md                 # This file
└── build/                    # React build output
    ├── index.html           # Main HTML file
    ├── manifest.json        # PWA manifest
    └── static/
        ├── js/
        │   └── main.a1b2c3d4.js
        └── css/
            └── main.5e6f7g8h.css
```

## Endpoints

- `http://localhost:8080/` - React application
- `http://localhost:8080/about` - React route (SPA)
- `http://localhost:8080/features` - React route (SPA)
- `http://localhost:8080/api/health` - Health check API
- `http://localhost:8080/api/info` - Application info API

## Using with a Real React Build

To use with a real React application:

1. Build your React app:
   ```bash
   cd your-react-app
   npm run build
   ```

2. Copy the build directory to `examples/react/build/`

3. Update the configuration if needed

4. Run the server:
   ```bash
   go run main.go
   ```

## Integration with Wayframe

This example integrates multiple Wayframe packages:

- `pkg/config` - Configuration management
- `pkg/logger` - Structured logging  
- `pkg/server` - HTTP server with middleware
- `pkg/react` - React application serving

## Advanced Usage

### Custom Base Path

Serve the React app at a different path:

```go
reactHandler, err := react.NewHandler(react.Config{
    BuildDir: "./build",
    BasePath: "/app",  // Serve at /app/*
})
```

### Embedded Filesystem

Use Go's embed directive for self-contained binaries:

```go
//go:embed build/*
var buildFS embed.FS

reactHandler, err := react.NewHandler(react.Config{
    FileSystem: buildFS,
    BasePath: "/",
})
```

### Custom 404 Handler

Provide a custom handler for missing files:

```go
reactHandler, err := react.NewHandler(react.Config{
    BuildDir: "./build",
    NotFoundHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "Custom 404", http.StatusNotFound)
    }),
})
```

## Building for Production

### With Go

```bash
go build -o react-server main.go
./react-server
```

### With Bazel

```bash
bazel build --compilation_mode=opt //examples/react
```

## Tips

1. **Development Mode**: Set `APP_ENV=development` and `LOG_LEVEL=DEBUG` for verbose logging
2. **Production Mode**: Use `APP_ENV=production` for optimized settings
3. **Environment Variables**: Prefix React env vars with `REACT_APP_` by convention
4. **API Proxy**: Register API routes before the React handler to avoid conflicts
5. **Caching**: Use content hashes in filenames for optimal caching

## License

See the main project LICENSE file.

