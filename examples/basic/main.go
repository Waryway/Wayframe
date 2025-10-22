// Package main demonstrates basic usage of the Wayframe framework.
// It shows how to use config, logger, and server packages together.
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Waryway/Wayframe/pkg/config"
	"github.com/Waryway/Wayframe/pkg/logger"
	"github.com/Waryway/Wayframe/pkg/server"
)

func main() {
	// Load configuration
	cfg := config.New("APP")
	port := cfg.String("PORT", "8080")
	logLevel := cfg.String("LOG_LEVEL", "INFO")
	shutdownTimeout := cfg.Duration("SHUTDOWN_TIMEOUT", 30*time.Second)
	
	// Setup logger
	level := logger.InfoLevel
	if logLevel == "DEBUG" {
		level = logger.DebugLevel
	}
	log := logger.New(level)
	
	log.Info("Starting Wayframe example application")
	log.WithField("port", port).Info("Configuration loaded")
	
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
	
	// Register routes
	srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("path", r.URL.Path).Debug("Handling request")
		fmt.Fprintf(w, "Welcome to Wayframe!\n")
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
	
	// Start server with graceful shutdown
	log.Infof("Server listening on :%s", port)
	if err := srv.Start(shutdownTimeout); err != nil {
		log.Errorf("Server error: %v", err)
	}
}
