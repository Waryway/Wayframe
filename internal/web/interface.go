// Package web provides web server abstractions for Wayframe applications.
// It defines a common interface for different web server implementations.
package web

import (
	"context"
	"net/http"
	"time"
)

// Server represents a web server interface that all implementations must satisfy.
type Server interface {
	// Use adds middleware to the server
	Use(middleware ...interface{})
	
	// Handle registers a handler for the given pattern
	Handle(pattern string, handler interface{})
	
	// HandleFunc registers a handler function for the given pattern
	HandleFunc(pattern string, handlerFunc interface{})
	
	// Start starts the server and blocks until shutdown
	Start(shutdownTimeout time.Duration) error
	
	// Shutdown gracefully shuts down the server
	Shutdown(ctx context.Context) error
	
	// Addr returns the server address
	Addr() string
}

// Config holds common configuration for web servers.
type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Middleware is a generic middleware function type.
type Middleware func(http.Handler) http.Handler
