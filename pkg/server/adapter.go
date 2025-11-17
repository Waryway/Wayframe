// Package server defines a common adapter interface for pluggable HTTP servers.
package server

import (
	"context"
	"time"
)

// Adapter is a generic interface implemented by server backends (stdlib, fiber, gorilla).
// It provides a minimal, flexible surface for examples/apps to register routes and middleware
// without importing a specific HTTP framework in their code.
type Adapter interface {
	// Use adds middleware to the server. The accepted types depend on the adapter
	// (e.g., func(http.Handler) http.Handler for stdlib/gorilla, fiber.Handler for fiber).
	Use(middleware ...interface{})

	// Handle registers a handler for the given pattern.
	Handle(pattern string, handler interface{})

	// HandleFunc registers a handler function for the given pattern.
	HandleFunc(pattern string, handlerFunc interface{})

	// Start starts the server and blocks until shutdown.
	Start(shutdownTimeout time.Duration) error

	// Shutdown gracefully shuts down the server.
	Shutdown(ctx context.Context) error

	// Addr returns the server address.
	Addr() string
}
