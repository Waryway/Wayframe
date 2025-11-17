// Package main demonstrates using Wayframe with Gorilla Mux router.
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Waryway/Wayframe/pkg/config"
	"github.com/Waryway/Wayframe/pkg/logger"
	"github.com/Waryway/Wayframe/pkg/server"
	gorillaserver "github.com/Waryway/Wayframe/pkg/server/gorilla"
)

func main() {
	// Load configuration (prefix APP_*)
	cfg := config.New("APP")
	_ = cfg.LoadFile("config.json") // optional

	port := cfg.Int("PORT", 8080)
	host := cfg.String("HOST", "0.0.0.0")
	readTimeout := cfg.Duration("READ_TIMEOUT", 10*time.Second)
	writeTimeout := cfg.Duration("WRITE_TIMEOUT", 10*time.Second)
	idleTimeout := cfg.Duration("IDLE_TIMEOUT", 120*time.Second)
	shutdownTimeout := cfg.Duration("SHUTDOWN_TIMEOUT", 30*time.Second)
	appEnv := cfg.String("ENVIRONMENT", "development")
	logLevel := cfg.String("LOG_LEVEL", "INFO")

	// Setup logger
	level := logger.InfoLevel
	if logLevel == "DEBUG" {
		level = logger.DebugLevel
	} else if logLevel == "WARN" {
		level = logger.WarnLevel
	} else if logLevel == "ERROR" {
		level = logger.ErrorLevel
	}
	log := logger.New(level)

	log.Info("Starting Wayframe Gorilla Mux example")
	log.WithField("port", port).Info("Configuration loaded")

	// Create server using pkg/server gorilla adapter
	srv := gorillaserver.New(server.Config{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	})

	// Add middleware
	srv.Use(gorillaserver.LoggingMiddleware(log))
	srv.Use(gorillaserver.RecoveryMiddleware(log))

	// Register routes
	srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("path", r.URL.Path).Debug("Handling request")
		_, _ = fmt.Fprintf(w, "Welcome to Wayframe with Gorilla Mux!\n")
		_, _ = fmt.Fprintf(w, "Environment: %s\n", appEnv)
	})

	srv.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "OK\n")
	})

	srv.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}
		log.WithField("name", name).Info("Greeting user")
		_, _ = fmt.Fprintf(w, "Hello, %s!\n", name)
	})

	// Start server
	log.Infof("Server listening on %s", srv.Addr())
	if err := srv.Start(shutdownTimeout); err != nil {
		log.Errorf("Server error: %v", err)
	}
}
