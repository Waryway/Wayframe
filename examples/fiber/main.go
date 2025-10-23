// Package main demonstrates using Wayframe with Fiber web framework.
package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/Waryway/Wayframe/internal/env"
	"github.com/Waryway/Wayframe/internal/web"
	fiberserver "github.com/Waryway/Wayframe/internal/web/fiber"
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
	log.Info("Starting Wayframe Fiber example")
	log.WithField("port", cfg.Port).Info("Configuration loaded")
	
	// Create server using web interface
	srv := fiberserver.New(web.Config{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})
	
	// Add middleware
	srv.Use(fiberserver.LoggingMiddleware(log))
	srv.Use(fiberserver.RecoveryMiddleware(log))
	
	// Register routes
	srv.HandleFunc("/", func(c *fiber.Ctx) error {
		log.WithField("path", c.Path()).Debug("Handling request")
		return c.SendString("Welcome to Wayframe with Fiber!\n")
	})
	
	srv.HandleFunc("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK\n")
	})
	
	srv.HandleFunc("/hello", func(c *fiber.Ctx) error {
		name := c.Query("name", "World")
		log.WithField("name", name).Info("Greeting user")
		return c.SendString(fmt.Sprintf("Hello, %s!\n", name))
	})
	
	// Start server
	log.Infof("Server listening on %s", srv.Addr())
	if err := srv.Start(cfg.ShutdownTimeout); err != nil {
		log.Errorf("Server error: %v", err)
	}
}
