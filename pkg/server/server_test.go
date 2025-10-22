package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type mockLogger struct {
	messages []string
}

func (m *mockLogger) Infof(format string, args ...interface{}) {
	m.messages = append(m.messages, fmt.Sprintf(format, args...))
}

func (m *mockLogger) Errorf(format string, args ...interface{}) {
	m.messages = append(m.messages, fmt.Sprintf(format, args...))
}

func TestNew(t *testing.T) {
	srv := New(Config{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})
	
	if srv == nil {
		t.Fatal("expected server to be created")
	}
	if srv.httpServer.Addr != ":8080" {
		t.Errorf("expected addr :8080, got %s", srv.httpServer.Addr)
	}
}

func TestHandle(t *testing.T) {
	srv := New(Config{Addr: ":0"})
	
	srv.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "test response")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	srv.mux.ServeHTTP(w, req)
	
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	
	if string(body) != "test response" {
		t.Errorf("expected 'test response', got '%s'", string(body))
	}
}

func TestLoggingMiddleware(t *testing.T) {
	mockLog := &mockLogger{}
	srv := New(Config{Addr: ":0"})
	
	srv.Use(LoggingMiddleware(mockLog))
	srv.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	srv.mux.ServeHTTP(w, req)
	
	if len(mockLog.messages) != 1 {
		t.Errorf("expected 1 log message, got %d", len(mockLog.messages))
	}
	if !strings.Contains(mockLog.messages[0], "GET") {
		t.Error("log should contain HTTP method")
	}
	if !strings.Contains(mockLog.messages[0], "/test") {
		t.Error("log should contain path")
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	mockLog := &mockLogger{}
	srv := New(Config{Addr: ":0"})
	
	srv.Use(RecoveryMiddleware(mockLog))
	srv.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})
	
	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	
	// Should not panic
	srv.mux.ServeHTTP(w, req)
	
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
	
	if len(mockLog.messages) != 1 {
		t.Errorf("expected 1 error log, got %d", len(mockLog.messages))
	}
	if !strings.Contains(mockLog.messages[0], "panic recovered") {
		t.Error("log should contain panic recovery message")
	}
}

func TestShutdown(t *testing.T) {
	srv := New(Config{Addr: ":0"})
	
	// Start server in background
	go func() {
		// Use a test server since we can't easily test Start()
		testServer := httptest.NewServer(srv.mux)
		defer testServer.Close()
	}()
	
	// Give it time to start
	time.Sleep(100 * time.Millisecond)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err := srv.Shutdown(ctx)
	if err != nil {
		t.Errorf("shutdown should not error: %v", err)
	}
}

func TestMiddlewareOrder(t *testing.T) {
	order := []string{}
	
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw1-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw1-after")
		})
	}
	
	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw2-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw2-after")
		})
	}
	
	srv := New(Config{Addr: ":0"})
	srv.Use(middleware1)
	srv.Use(middleware2)
	
	srv.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	srv.mux.ServeHTTP(w, req)
	
	expected := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}
	if len(order) != len(expected) {
		t.Errorf("expected %d elements, got %d", len(expected), len(order))
	}
	for i := range expected {
		if i >= len(order) || order[i] != expected[i] {
			t.Errorf("at position %d: expected %s, got %s", i, expected[i], order[i])
		}
	}
}
