// Package main demonstrates serving a React application with Wayframe.
// It shows how to use the react package with the server package to serve
// a React SPA with environment variable injection and proper routing.
package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	env "github.com/Waryway/Wayframe/examples/react/env"
	"github.com/Waryway/Wayframe/pkg/react"
	"github.com/Waryway/Wayframe/pkg/server"
)

func findReactBuildDir(defaultDir string) string {
	// 1. Use REACT_BUILD_DIR env var if set (Bazel can set this)
	if dir := os.Getenv("REACT_BUILD_DIR"); dir != "" {
		return dir
	}
	// 2. Try to find in Bazel runfiles (if running under Bazel)
	exe, err := os.Executable()
	if err == nil {
		runfilesDir := filepath.Join(filepath.Dir(exe), "../react/build")
		if stat, err := os.Stat(runfilesDir); err == nil && stat.IsDir() {
			return runfilesDir
		}
	}
	// 3. Fallback to default (./build)
	return defaultDir
}

func main() {
	// Load environment for the React example
	e := env.New("APP")
	e.LoadFiles("config.json")
	log := e.Logger()

	port := e.String("PORT", "8080")
	shutdownTimeout := e.Duration("SHUTDOWN_TIMEOUT", 30*time.Second)

	// React app configuration
	defaultBuildDir := e.String("REACT_BUILD_DIR", "./build")
	reactBuildDir := findReactBuildDir(defaultBuildDir)
	reactBasePath := e.String("REACT_BASE_PATH", "/")
	apiURL := e.String("API_URL", "http://localhost:8080")
	appVersion := e.String("APP_VERSION", "1.0.0")
	appEnv := e.String("APP_ENV", "production")

	log.Info("Starting Wayframe React example")
	log.WithField("port", port).Info("Configuration loaded")
	log.WithField("build_dir", reactBuildDir).Info("React build directory")

	// Create React handler with environment variables
	reactHandler, err := react.NewHandler(react.Config{
		BuildDir: reactBuildDir,
		BasePath: reactBasePath,
		EnvVars: map[string]string{
			"REACT_APP_API_URL": apiURL,
			"REACT_APP_VERSION": appVersion,
			"REACT_APP_ENV":     appEnv,
		},
		CacheMaxAge: 365 * 24 * time.Hour, // 1 year for hashed assets
	})
	if err != nil {
		log.WithField("path", reactBuildDir).WithField("err", err).Error("Failed to create React handler")
		return
	}

	log.Info("React handler initialized successfully")

	// Create server
	srv := server.New(server.Config{
		Addr:         fmt.Sprintf(":%s", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	// Add middleware
	srv.Use(server.LoggingMiddleware(log))
	srv.Use(server.RecoveryMiddleware(log))

	// Register API routes (before React handler)
	srv.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"status":"ok","version":"%s"}`, appVersion)
	})

	srv.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{
			"api_url": "%s",
			"version": "%s",
			"environment": "%s",
			"timestamp": "%s"
		}`, apiURL, appVersion, appEnv, time.Now().Format(time.RFC3339))
	})

	// Register React handler (catches all remaining routes)
	srv.Handle("/", reactHandler)

	// Start server with graceful shutdown
	log.Infof("Server listening on :%s", port)
	log.Info("Open your browser and navigate to http://localhost:" + port)
	log.Info("API endpoints available at /api/health and /api/info")

	if err := srv.Start(shutdownTimeout); err != nil {
		log.Errorf("Server error: %v", err)
	}
}
