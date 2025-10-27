// Package gorilla provides a Gorilla Mux router server implementation.
package gorilla

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/Waryway/Wayframe/internal/web"
)

// Server wraps Gorilla Mux with the web.Server interface.
type Server struct {
	httpServer *http.Server
	router     *mux.Router
	middleware []mux.MiddlewareFunc
	addr       string
}

// New creates a new Gorilla Mux server with the given configuration.
func New(cfg web.Config) web.Server {
	router := mux.NewRouter()
	
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.Addr,
			Handler:      router,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		router:     router,
		middleware: make([]mux.MiddlewareFunc, 0),
		addr:       cfg.Addr,
	}
}

// Use adds middleware to the server.
func (s *Server) Use(middleware ...interface{}) {
	for _, mw := range middleware {
		if m, ok := mw.(mux.MiddlewareFunc); ok {
			s.router.Use(m)
		} else if m, ok := mw.(func(http.Handler) http.Handler); ok {
			s.router.Use(mux.MiddlewareFunc(m))
		}
	}
}

// Handle registers a handler for the given pattern.
func (s *Server) Handle(pattern string, handler interface{}) {
	if h, ok := handler.(http.Handler); ok {
		s.router.Handle(pattern, h)
	} else if h, ok := handler.(http.HandlerFunc); ok {
		s.router.Handle(pattern, h)
	} else if h, ok := handler.(func(http.ResponseWriter, *http.Request)); ok {
		s.router.HandleFunc(pattern, h)
	} else {
		panic(fmt.Sprintf("unsupported handler type: %T", handler))
	}
}

// HandleFunc registers a handler function for the given pattern.
func (s *Server) HandleFunc(pattern string, handlerFunc interface{}) {
	if h, ok := handlerFunc.(func(http.ResponseWriter, *http.Request)); ok {
		s.router.HandleFunc(pattern, h)
	} else if h, ok := handlerFunc.(http.HandlerFunc); ok {
		s.router.Handle(pattern, h)
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
func LoggingMiddleware(logger interface{ Infof(string, ...interface{}) }) mux.MiddlewareFunc {
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
func RecoveryMiddleware(logger interface{ Errorf(string, ...interface{}) }) mux.MiddlewareFunc {
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
