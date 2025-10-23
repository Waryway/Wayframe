// Package main demonstrates using Wayframe with Gorilla Mux router.
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Waryway/Wayframe/internal/env"
	"github.com/Waryway/Wayframe/internal/web"
	gorillaserver "github.com/Waryway/Wayframe/internal/web/gorilla"
	"github.com/Waryway/Wayframe/pkg/logger"
)

// AppConfig defines the application configuration structure.
type AppConfig struct {
	Port            int           `config:"port" env:"APP_PORT" default:"8080"`
	LogLevel        string        `config:"log_level" env:"LOG_LEVEL" default:"INFO"`
	ShutdownTimeout time.Duration `config:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" default:"30s"`
}

func main() {
	// Initialize environment
	e := env.New("APP")
	
	// Load configuration
	var cfg AppConfig
	if err := e.LoadConfig(&cfg); err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	
	// Set up logger
	level := logger.InfoLevel
	if cfg.LogLevel == "DEBUG" {
		level = logger.DebugLevel
	}
	e.InitLogger(level)
	
	log := e.GetLogger()
	log.Info("Starting Wayframe Gorilla Mux example")
	log.WithField("port", cfg.Port).Info("Configuration loaded")
	
	// Create server using web interface
	srv := gorillaserver.New(web.Config{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})
	
	// Add middleware
	srv.Use(gorillaserver.LoggingMiddleware(log))
	srv.Use(gorillaserver.RecoveryMiddleware(log))
	
	// Register routes
	srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.WithField("path", r.URL.Path).Debug("Handling request")
		fmt.Fprintf(w, "Welcome to Wayframe with Gorilla Mux!\n")
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
