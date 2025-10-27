// Package server provides HTTP server utilities for Wayframe applications.
//
// The server package wraps Go's standard http.Server with convenient
// features like graceful shutdown, middleware support, and common patterns.
//
// # Basic Usage
//
//	srv := server.New(server.Config{
//	    Addr:         ":8080",
//	    ReadTimeout:  10 * time.Second,
//	    WriteTimeout: 10 * time.Second,
//	})
//
//	srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//	    fmt.Fprintf(w, "Hello, World!")
//	})
//
//	srv.Start(30 * time.Second)
//
// # Middleware
//
// Add middleware to process requests:
//
//	srv.Use(server.LoggingMiddleware(log))
//	srv.Use(server.RecoveryMiddleware(log))
//
// Middleware is applied in the order it's added. The first middleware
// added is the outermost wrapper.
//
// # Built-in Middleware
//
// The package includes common middleware:
//   - LoggingMiddleware: Logs each request with method, path, and duration
//   - RecoveryMiddleware: Recovers from panics and returns 500 errors
//
// # Graceful Shutdown
//
// The Start method handles graceful shutdown automatically:
//   - Listens for SIGINT and SIGTERM signals
//   - Stops accepting new connections
//   - Waits for existing requests to complete (up to timeout)
//   - Returns when shutdown is complete
//
// # Example
//
//	srv := server.New(server.Config{Addr: ":8080"})
//	srv.Use(server.LoggingMiddleware(log))
//
//	srv.HandleFunc("/", indexHandler)
//	srv.HandleFunc("/api/users", usersHandler)
//
//	if err := srv.Start(30 * time.Second); err != nil {
//	    log.Fatalf("Server error: %v", err)
//	}
package server
