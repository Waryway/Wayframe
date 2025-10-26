// Package server provides HTTP server utilities for Wayframe applications.
// It offers graceful shutdown, middleware support, and common HTTP patterns.
package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server wraps http.Server with graceful shutdown capabilities.
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	middleware []Middleware
}

// Middleware is a function that wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// Config holds the configuration for creating a new Server.
type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// New creates a new Server with the given configuration.
func New(cfg Config) *Server {
	mux := http.NewServeMux()
	
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.Addr,
			Handler:      mux,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		mux:        mux,
		middleware: make([]Middleware, 0),
	}
}

// Use adds middleware to the server. Middleware is applied in the order it's added.
func (s *Server) Use(mw Middleware) {
	s.middleware = append(s.middleware, mw)
}

// Handle registers a handler for the given pattern.
// Middleware is applied to the handler.
func (s *Server) Handle(pattern string, handler http.Handler) {
	// Apply middleware in reverse order so first added is outermost
	for i := len(s.middleware) - 1; i >= 0; i-- {
		handler = s.middleware[i](handler)
	}
	s.mux.Handle(pattern, handler)
}

// HandleFunc registers a handler function for the given pattern.
// Middleware is applied to the handler.
func (s *Server) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	s.Handle(pattern, handlerFunc)
}

// Start starts the HTTP server and blocks until a shutdown signal is received.
// It performs graceful shutdown with a timeout.
func (s *Server) Start(shutdownTimeout time.Duration) error {
	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	// Channel to receive server errors
	errChan := make(chan error, 1)
	
	// Start server in a goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()
	
	// Wait for interrupt signal or error
	select {
	case err := <-errChan:
		return err
	case sig := <-quit:
		fmt.Printf("Received signal: %v, shutting down gracefully...\n", sig)
	}
	
	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	
	// Attempt graceful shutdown
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}
	
	fmt.Println("Server exited gracefully")
	return nil
}

// Shutdown gracefully shuts down the server with the given context.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// LoggingMiddleware logs each HTTP request with method, path, and duration.
func LoggingMiddleware(logger interface{ Infof(string, ...interface{}) }) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			logger.Infof("%s %s - %v", r.Method, r.URL.Path, duration)
		})
	}
}

// RecoveryMiddleware recovers from panics and returns a 500 Internal Server Error.
func RecoveryMiddleware(logger interface{ Errorf(string, ...interface{}) }) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Errorf("panic recovered: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
