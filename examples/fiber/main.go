// Package main demonstrates using Wayframe with Fiber web framework.
package main

import (
	"fmt"

	"github.com/Waryway/Wayframe/internal/env"
	"github.com/Waryway/Wayframe/internal/web"
	fiberserver "github.com/Waryway/Wayframe/internal/web/fiber"
	"github.com/gofiber/fiber/v2"
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

	log.Info("Starting Wayframe Fiber example")
	log.WithField("port", cfg.Port).Info("Configuration loaded")

	// Create server using web interface
	srv := fiberserver.New(web.Config{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	})

	// Add middleware
	srv.Use(fiberserver.LoggingMiddleware(log))
	srv.Use(fiberserver.RecoveryMiddleware(log))

	// Register routes
	srv.HandleFunc("/", func(c *fiber.Ctx) error {
		log.WithField("path", c.Path()).Debug("Handling request")
		return c.SendString(fmt.Sprintf("Welcome to Wayframe with Fiber!\nEnvironment: %s\n", cfg.Environment))
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
