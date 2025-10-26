// Package stdlib provides a standard library HTTP server implementation.
package stdlib

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Waryway/Wayframe/internal/web"
)

// Server wraps http.Server with graceful shutdown capabilities.
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	middleware []web.Middleware
	addr       string
}

// New creates a new stdlib Server with the given configuration.
func New(cfg web.Config) web.Server {
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
		middleware: make([]web.Middleware, 0),
		addr:       cfg.Addr,
	}
}

// Use adds middleware to the server.
func (s *Server) Use(middleware ...interface{}) {
	for _, mw := range middleware {
		if m, ok := mw.(web.Middleware); ok {
			s.middleware = append(s.middleware, m)
		} else if m, ok := mw.(func(http.Handler) http.Handler); ok {
			s.middleware = append(s.middleware, web.Middleware(m))
		}
	}
}

// Handle registers a handler for the given pattern.
func (s *Server) Handle(pattern string, handler interface{}) {
	var h http.Handler
	if hh, ok := handler.(http.Handler); ok {
		h = hh
	} else if hf, ok := handler.(http.HandlerFunc); ok {
		h = hf
	} else if hf, ok := handler.(func(http.ResponseWriter, *http.Request)); ok {
		h = http.HandlerFunc(hf)
	} else {
		panic(fmt.Sprintf("unsupported handler type: %T", handler))
	}
	
	// Apply middleware in reverse order
	for i := len(s.middleware) - 1; i >= 0; i-- {
		h = s.middleware[i](h)
	}
	s.mux.Handle(pattern, h)
}

// HandleFunc registers a handler function for the given pattern.
func (s *Server) HandleFunc(pattern string, handlerFunc interface{}) {
	if hf, ok := handlerFunc.(func(http.ResponseWriter, *http.Request)); ok {
		s.Handle(pattern, http.HandlerFunc(hf))
	} else if hf, ok := handlerFunc.(http.HandlerFunc); ok {
		s.Handle(pattern, hf)
	} else {
		panic(fmt.Sprintf("unsupported handler function type: %T", handlerFunc))
	}
}

// Start starts the HTTP server and blocks until a shutdown signal is received.
func (s *Server) Start(shutdownTimeout time.Duration) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	errChan := make(chan error, 1)
	
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()
	
	select {
	case err := <-errChan:
		return err
	case sig := <-quit:
		fmt.Printf("Received signal: %v, shutting down gracefully...\n", sig)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}
	
	fmt.Println("Server exited gracefully")
	return nil
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// Addr returns the server address.
func (s *Server) Addr() string {
	return s.addr
}

// LoggingMiddleware logs each HTTP request.
func LoggingMiddleware(logger interface{ Infof(string, ...interface{}) }) web.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			logger.Infof("%s %s - %v", r.Method, r.URL.Path, duration)
		})
	}
}

// RecoveryMiddleware recovers from panics.
func RecoveryMiddleware(logger interface{ Errorf(string, ...interface{}) }) web.Middleware {
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
