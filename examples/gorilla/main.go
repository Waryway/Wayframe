// Package main demonstrates using Wayframe with Gorilla Mux router.
package main

import (
	"fmt"
	"net/http"

	"github.com/Waryway/Wayframe/internal/env"
	"github.com/Waryway/Wayframe/internal/web"
	gorillaserver "github.com/Waryway/Wayframe/internal/web/gorilla"
)

func main() {
	// Initialize environment
	e := env.New("APP")

	// Load standard configuration
	if err := e.LoadStandardConfig(); err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// Get the standard config object
	cfg := e.GetAppConfig()
	log := e.GetLogger()

	log.Info("Starting Wayframe Gorilla Mux example")
	log.WithField("port", cfg.Port).Info("Configuration loaded")

	// Create server using web interface
	srv := gorillaserver.New(web.Config{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	})

	// Add middleware
	srv.Use(gorillaserver.LoggingMiddleware(log))
	srv.Use(gorillaserver.RecoveryMiddleware(log))

	// Register routes
	srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("path", r.URL.Path).Debug("Handling request")
		fmt.Fprintf(w, "Welcome to Wayframe with Gorilla Mux!\n")
		fmt.Fprintf(w, "Environment: %s\n", cfg.Environment)
	})

	srv.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK\n")
	})

	srv.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}
		log.WithField("name", name).Info("Greeting user")
		fmt.Fprintf(w, "Hello, %s!\n", name)
	})

	// Start server
	log.Infof("Server listening on %s", srv.Addr())
	if err := srv.Start(cfg.ShutdownTimeout); err != nil {
		log.Errorf("Server error: %v", err)
	}
}
