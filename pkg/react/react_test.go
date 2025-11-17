package react

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestNewHandler(t *testing.T) {
	// Create a test filesystem
	fsys := fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data: []byte(`<!DOCTYPE html><html><head></head><body><div id="root"></div></body></html>`),
		},
		"static/js/main.js": &fstest.MapFile{
			Data: []byte(`console.log("test");`),
		},
	}

	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config with filesystem",
			config: Config{
				FileSystem: fsys,
				BasePath:   "/",
			},
			wantErr: false,
		},
		{
			name: "valid config with custom index",
			config: Config{
				FileSystem: fsys,
				BasePath:   "/app",
				IndexFile:  "index.html",
			},
			wantErr: false,
		},
		{
			name: "valid config with env vars",
			config: Config{
				FileSystem: fsys,
				BasePath:   "/",
				EnvVars: map[string]string{
					"REACT_APP_API_URL": "http://localhost:8080",
				},
			},
			wantErr: false,
		},
		{
			name: "missing index file",
			config: Config{
				FileSystem: fstest.MapFS{},
				BasePath:   "/",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, err := NewHandler(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && handler == nil {
				t.Error("NewHandler() returned nil handler")
			}
		})
	}
}

func TestHandler_ServeHTTP(t *testing.T) {
	// Create a test filesystem
	fsys := fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data: []byte(`<!DOCTYPE html><html><head></head><body><div id="root"></div></body></html>`),
		},
		"static/js/main.abc123.js": &fstest.MapFile{
			Data: []byte(`console.log("main");`),
		},
		"static/css/style.css": &fstest.MapFile{
			Data: []byte(`body { margin: 0; }`),
		},
		"favicon.ico": &fstest.MapFile{
			Data: []byte{0x00, 0x00, 0x01, 0x00},
		},
	}

	handler, err := NewHandler(Config{
		FileSystem: fsys,
		BasePath:   "/",
	})
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	tests := []struct {
		name             string
		path             string
		wantStatus       int
		wantContentType  string
		wantCacheControl string
	}{
		{
			name:             "index.html",
			path:             "/",
			wantStatus:       http.StatusOK,
			wantContentType:  "text/html",
			wantCacheControl: "no-cache",
		},
		{
			name:             "explicit index.html",
			path:             "/index.html",
			wantStatus:       http.StatusOK,
			wantContentType:  "text/html",
			wantCacheControl: "no-cache",
		},
		{
			name:             "static js with hash",
			path:             "/static/js/main.abc123.js",
			wantStatus:       http.StatusOK,
			wantContentType:  "text/javascript",
			wantCacheControl: "public, max-age",
		},
		{
			name:             "static css",
			path:             "/static/css/style.css",
			wantStatus:       http.StatusOK,
			wantContentType:  "text/css",
			wantCacheControl: "public, max-age",
		},
		{
			name:             "favicon",
			path:             "/favicon.ico",
			wantStatus:       http.StatusOK,
			wantCacheControl: "public, max-age",
		},
		{
			name:             "SPA route - about",
			path:             "/about",
			wantStatus:       http.StatusOK,
			wantContentType:  "text/html",
			wantCacheControl: "no-cache",
		},
		{
			name:             "SPA route - user/123",
			path:             "/user/123",
			wantStatus:       http.StatusOK,
			wantContentType:  "text/html",
			wantCacheControl: "no-cache",
		},
		{
			name:       "missing asset",
			path:       "/static/missing.js",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Status = %v, want %v", rec.Code, tt.wantStatus)
			}

			if tt.wantContentType != "" {
				contentType := rec.Header().Get("Content-Type")
				if !strings.Contains(contentType, tt.wantContentType) {
					t.Errorf("Content-Type = %v, want to contain %v", contentType, tt.wantContentType)
				}
			}

			if tt.wantCacheControl != "" {
				cacheControl := rec.Header().Get("Cache-Control")
				if !strings.Contains(cacheControl, tt.wantCacheControl) {
					t.Errorf("Cache-Control = %v, want to contain %v", cacheControl, tt.wantCacheControl)
				}
			}
		})
	}
}

func TestInjectEnvVars(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		envVars  map[string]string
		want     string
		modified bool
	}{
		{
			name: "inject before head",
			html: `<!DOCTYPE html><html><head></head><body></body></html>`,
			envVars: map[string]string{
				"REACT_APP_API_URL": "http://localhost:8080",
			},
			want:     `<script>window.__REACT_ENV__ = {"REACT_APP_API_URL":"http://localhost:8080"};</script></head>`,
			modified: true,
		},
		{
			name: "inject multiple vars",
			html: `<!DOCTYPE html><html><head></head><body></body></html>`,
			envVars: map[string]string{
				"REACT_APP_API_URL": "http://localhost:8080",
				"REACT_APP_VERSION": "1.0.0",
			},
			modified: true,
		},
		{
			name:     "no env vars",
			html:     `<!DOCTYPE html><html><head></head><body></body></html>`,
			envVars:  map[string]string{},
			want:     `<!DOCTYPE html><html><head></head><body></body></html>`,
			modified: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, modified := injectEnvVars([]byte(tt.html), tt.envVars)

			if modified != tt.modified {
				t.Errorf("modified = %v, want %v", modified, tt.modified)
			}

			if tt.want != "" && !strings.Contains(string(got), tt.want) {
				t.Errorf("result doesn't contain expected string:\ngot: %s\nwant to contain: %s", got, tt.want)
			}

			// Verify all env vars are present
			for key, value := range tt.envVars {
				if !strings.Contains(string(got), key) {
					t.Errorf("result doesn't contain key: %s", key)
				}
				if !strings.Contains(string(got), value) {
					t.Errorf("result doesn't contain value: %s", value)
				}
			}
		})
	}
}

func TestIsAssetPath(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/static/js/main.js", true},
		{"/static/css/style.css", true},
		{"/images/logo.png", true},
		{"/favicon.ico", true},
		{"/manifest.json", true},
		{"/about", false},
		{"/user/123", false},
		{"/", false},
		{"/static/bundle.map", true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isAssetPath(tt.path); got != tt.want {
				t.Errorf("isAssetPath(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestHasContentHash(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"main.abc123def.js", true},
		{"main.abc123def456.chunk.js", true},
		{"style.1a2b3c4d.css", true},
		{"main.js", false},
		{"style.css", false},
		{"app.min.js", false},
		{"bundle.12345678.js", true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := hasContentHash(tt.path); got != tt.want {
				t.Errorf("hasContentHash(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestHandler_WithBasePath(t *testing.T) {
	fsys := fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data: []byte(`<!DOCTYPE html><html><head></head><body><div id="root"></div></body></html>`),
		},
		"static/js/main.js": &fstest.MapFile{
			Data: []byte(`console.log("test");`),
		},
	}

	handler, err := NewHandler(Config{
		FileSystem: fsys,
		BasePath:   "/app",
	})
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "base path index",
			path:       "/app",
			wantStatus: http.StatusOK,
		},
		{
			name:       "base path with slash",
			path:       "/app/",
			wantStatus: http.StatusOK,
		},
		{
			name:       "base path static asset",
			path:       "/app/static/js/main.js",
			wantStatus: http.StatusOK,
		},
		{
			name:       "base path SPA route",
			path:       "/app/about",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

// Benchmark tests
func BenchmarkHandler_ServeIndex(b *testing.B) {
	fsys := fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data: []byte(strings.Repeat(`<!DOCTYPE html><html><head></head><body><div id="root"></div></body></html>`, 10)),
		},
	}

	handler, _ := NewHandler(Config{
		FileSystem: fsys,
		BasePath:   "/",
	})

	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

func BenchmarkHandler_ServeStatic(b *testing.B) {
	fsys := fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data: []byte(`<!DOCTYPE html><html><head></head><body><div id="root"></div></body></html>`),
		},
		"static/js/main.js": &fstest.MapFile{
			Data: []byte(strings.Repeat(`console.log("test");`, 100)),
		},
	}

	handler, _ := NewHandler(Config{
		FileSystem: fsys,
		BasePath:   "/",
	})

	req := httptest.NewRequest("GET", "/static/js/main.js", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

// Integration test with actual filesystem (optional, requires test data)
func TestHandler_RealFilesystem(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	indexHTML := `<!DOCTYPE html><html><head></head><body><div id="root"></div></body></html>`
	if err := os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(indexHTML), 0644); err != nil {
		t.Fatalf("Failed to write index.html: %v", err)
	}

	staticDir := filepath.Join(tmpDir, "static", "js")
	if err := os.MkdirAll(staticDir, 0755); err != nil {
		t.Fatalf("Failed to create static dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(staticDir, "main.js"), []byte(`console.log("test");`), 0644); err != nil {
		t.Fatalf("Failed to write main.js: %v", err)
	}

	handler, err := NewHandler(Config{
		BuildDir: tmpDir,
		BasePath: "/",
	})
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	// Test serving index
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status = %v, want %v", rec.Code, http.StatusOK)
	}

	// Test serving static file
	req = httptest.NewRequest("GET", "/static/js/main.js", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status = %v, want %v", rec.Code, http.StatusOK)
	}
}
