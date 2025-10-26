// Package fiber provides a Fiber v2 web framework server implementation.
package fiber

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/Waryway/Wayframe/internal/web"
)

// Server wraps Fiber app with the web.Server interface.
type Server struct {
	app  *fiber.App
	addr string
}

// New creates a new Fiber server with the given configuration.
func New(cfg web.Config) web.Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	})
	
	return &Server{
		app:  app,
		addr: cfg.Addr,
	}
}

// Use adds middleware to the server.
func (s *Server) Use(middleware ...interface{}) {
	for _, mw := range middleware {
		if m, ok := mw.(func(*fiber.Ctx) error); ok {
			s.app.Use(m)
		} else if m, ok := mw.(fiber.Handler); ok {
			s.app.Use(m)
		}
	}
}

// Handle registers a handler for the given pattern.
func (s *Server) Handle(pattern string, handler interface{}) {
	if h, ok := handler.(func(*fiber.Ctx) error); ok {
		s.app.All(pattern, h)
	} else if h, ok := handler.(fiber.Handler); ok {
		s.app.All(pattern, h)
	} else {
		panic(fmt.Sprintf("unsupported handler type: %T", handler))
	}
}

// HandleFunc registers a handler function for the given pattern.
func (s *Server) HandleFunc(pattern string, handlerFunc interface{}) {
	s.Handle(pattern, handlerFunc)
}

// Start starts the Fiber server and blocks until shutdown.
func (s *Server) Start(shutdownTimeout time.Duration) error {
	errChan := make(chan error, 1)
	
	go func() {
		if err := s.app.Listen(s.addr); err != nil {
			errChan <- err
		}
	}()
	
	// Wait for error
	if err := <-errChan; err != nil {
		return err
	}
	
	return nil
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

// Addr returns the server address.
func (s *Server) Addr() string {
	return s.addr
}

// LoggingMiddleware logs each HTTP request.
func LoggingMiddleware(logger interface{ Infof(string, ...interface{}) }) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		logger.Infof("%s %s - %v", c.Method(), c.Path(), duration)
		return err
	}
}

// RecoveryMiddleware recovers from panics.
func RecoveryMiddleware(logger interface{ Errorf(string, ...interface{}) }) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("panic recovered: %v", err)
				c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}
		}()
		return c.Next()
	}
}
