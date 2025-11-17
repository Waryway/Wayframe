# React Bootstrap CLI

A command-line tool to bootstrap React applications for use with Wayframe.

## Features

- Create React apps with multiple templates (CRA, Vite, Next.js)
- TypeScript support
- Interactive mode
- Build and dependency management helpers

## Usage

### Interactive Mode

The easiest way to get started:

```bash
# Using Go
go run main.go -interactive

# Using Bazel
bazel run //examples/react-bootstrap -- -interactive
```

### Command-Line Mode

Create a React app with specific options:

```bash
# Create a Vite app with TypeScript
go run main.go -name my-app -template vite -typescript

# Create with Bazel build support
go run main.go -name my-app -template vite -typescript -bazel

# Create a Create React App
go run main.go -name my-app -template cra

# Create a Next.js app
go run main.go -name my-app -template next -typescript

# Create in a specific directory
go run main.go -name my-app -dir ./apps -template vite
```

### Build Existing App

Build an existing React application:

```bash
# Build with npm
go run main.go -build ./my-app

# Build with Bazel
go run main.go -build-bazel //apps/my-app:build
```

### Install Dependencies

Install dependencies for an existing React app:

```bash
go run main.go -install ./my-app
```

## Flags

- `-name` - Name of the React application to create (required for bootstrap)
- `-dir` - Directory where to create the app (default: current directory)
- `-template` - Template to use: `cra`, `vite`, or `next` (default: vite)
- `-typescript` - Use TypeScript (default: false)
- `-bazel` - Setup Bazel build files (default: false)
- `-skip-install` - Skip npm install step (default: false)
- `-skip-git` - Skip git initialization (default: false)
- `-interactive` - Run in interactive mode
- `-build` - Build an existing React app (provide app directory)
- `-build-bazel` - Build a React app with Bazel (provide Bazel target)
- `-install` - Install dependencies for an existing React app (provide app directory)

## Templates

### Vite (Recommended)

Fast, modern build tool with excellent developer experience.

```bash
go run main.go -name my-app -template vite
```

### Create React App (CRA)

The traditional Create React App.

```bash
go run main.go -name my-app -template cra
```

### Next.js

React framework with server-side rendering and routing.

```bash
go run main.go -name my-app -template next
```

## Examples

### Create a TypeScript Vite App

```bash
go run main.go -name my-awesome-app -template vite -typescript
cd my-awesome-app
npm run dev  # Start development server
npm run build  # Build for production
```

### Use with Wayframe

After creating and building your React app:

```go
package main

import (
    "github.com/Waryway/Wayframe/pkg/react"
    "github.com/Waryway/Wayframe/pkg/server"
    "time"
)

func main() {
    // Create React handler
    reactHandler, _ := react.NewHandler(react.Config{
        BuildDir: "./my-awesome-app/dist",  // or 'build' for CRA
        BasePath: "/",
        EnvVars: map[string]string{
            "REACT_APP_API_URL": "http://localhost:8080/api",
        },
    })
    
    // Create server
    srv := server.New(server.Config{Addr: ":8080"})
    srv.Handle("/", reactHandler)
    srv.Start(30 * time.Second)
}
```

## Requirements

- Node.js (v14 or higher)
- npm or yarn
- Internet connection (for downloading templates)

## Tips

1. **Vite is faster** than CRA for development and builds
2. Use `-skip-git` if you're creating the app in an existing git repository
3. The `-skip-install` flag is useful for CI/CD pipelines where you want to cache dependencies separately
4. For Vite apps, the build output is in `dist/`, for CRA it's in `build/`

## Troubleshooting

### "node not found"

Install Node.js from [nodejs.org](https://nodejs.org/)

### "npm not found"

npm should be installed with Node.js. Try reinstalling Node.js.

### Permission errors

On Linux/Mac, you may need to use `sudo` or fix npm permissions:
```bash
npm config set prefix ~/.npm-global
export PATH=~/.npm-global/bin:$PATH
```

## License

See the main project LICENSE file.

