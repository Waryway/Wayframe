// Package react provides utilities for serving React applications from Go servers.
// It handles static file serving, SPA routing, and React build-time environment variables.
//
// The package integrates seamlessly with Wayframe's server package and supports:
//   - Static asset serving with proper caching headers
//   - SPA fallback routing (all non-asset routes serve index.html)
//   - Environment variable injection into React builds
//   - Gzip/Brotli compression support
//   - Development and production modes
//   - Bootstrapping new React applications (CRA, Vite, Next.js)
//
// Basic usage:
//
//	reactHandler := react.NewHandler(react.Config{
//	    BuildDir: "./build",
//	    BasePath: "/",
//	})
//
//	srv := server.New(server.Config{Addr: ":8080"})
//	srv.Handle("/", reactHandler)
//	srv.Start(30 * time.Second)
//
// With environment variable injection:
//
//	reactHandler := react.NewHandler(react.Config{
//	    BuildDir: "./build",
//	    BasePath: "/",
//	    EnvVars: map[string]string{
//	        "REACT_APP_API_URL": "https://api.example.com",
//	        "REACT_APP_VERSION": "1.0.0",
//	    },
//	})
//
// Bootstrap a new React application:
//
//	err := react.Bootstrap(react.BootstrapConfig{
//	    AppName:    "my-app",
//	    Template:   "vite",
//	    TypeScript: true,
//	})
package react
