// Package main demonstrates using Wayframe with Fiber web framework.
package main

import (
	"fmt"
	"time"

	"github.com/Waryway/Wayframe/pkg/config"
	"github.com/Waryway/Wayframe/pkg/logger"
	"github.com/Waryway/Wayframe/pkg/server"
	fiberserver "github.com/Waryway/Wayframe/pkg/server/fiber"
	"github.com/gofiber/fiber/v2"
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

	log.Info("Starting Wayframe Fiber example")
	log.WithField("port", port).Info("Configuration loaded")

	// Create server using pkg/server fiber adapter
	srv := fiberserver.New(server.Config{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	})

	// Add middleware
	srv.Use(fiberserver.LoggingMiddleware(log))
	srv.Use(fiberserver.RecoveryMiddleware(log))

	// Register routes
	srv.HandleFunc("/", func(c *fiber.Ctx) error {
		log.WithField("path", c.Path()).Debug("Handling request")
		return c.SendString(fmt.Sprintf("Welcome to Wayframe with Fiber!\nEnvironment: %s\n", appEnv))
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
	if err := srv.Start(shutdownTimeout); err != nil {
		log.Errorf("Server error: %v", err)
	}
}
