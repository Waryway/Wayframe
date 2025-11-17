package react

import (
	"bytes"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Handler serves a React application with SPA routing support.
type Handler struct {
	config        Config
	fileServer    http.Handler
	indexHTML     []byte
	indexModified bool
}

// Config holds the configuration for serving a React application.
type Config struct {
	// BuildDir is the directory containing the React build output (e.g., "./build")
	BuildDir string

	// BasePath is the base URL path where the React app is served (e.g., "/" or "/app")
	BasePath string

	// EnvVars are environment variables to inject into the React app.
	// These replace placeholders in the build or are made available at runtime.
	// Keys should follow React's REACT_APP_ convention.
	EnvVars map[string]string

	// IndexFile is the name of the index file (default: "index.html")
	IndexFile string

	// CacheMaxAge sets the Cache-Control max-age for static assets (default: 1 year)
	// Index.html is never cached.
	CacheMaxAge time.Duration

	// NotFoundHandler is called when a file is not found (optional)
	// If nil, serves index.html for SPA routing
	NotFoundHandler http.Handler

	// FileSystem allows using an embedded filesystem (optional)
	// If nil, uses os.DirFS(BuildDir)
	FileSystem fs.FS
}

// NewHandler creates a new React application handler.
func NewHandler(cfg Config) (*Handler, error) {
	// Set defaults
	if cfg.IndexFile == "" {
		cfg.IndexFile = "index.html"
	}
	if cfg.BasePath == "" {
		cfg.BasePath = "/"
	}
	if cfg.CacheMaxAge == 0 {
		cfg.CacheMaxAge = 365 * 24 * time.Hour // 1 year
	}
	if cfg.EnvVars == nil {
		cfg.EnvVars = make(map[string]string)
	}

	// Determine filesystem
	var fsys fs.FS
	if cfg.FileSystem != nil {
		fsys = cfg.FileSystem
	} else {
		if cfg.BuildDir == "" {
			return nil, fmt.Errorf("BuildDir is required when FileSystem is not provided")
		}
		fsys = os.DirFS(cfg.BuildDir)
	}

	// Read and process index.html
	indexHTML, err := fs.ReadFile(fsys, cfg.IndexFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", cfg.IndexFile, err)
	}

	// Inject environment variables into index.html
	indexModified := false
	if len(cfg.EnvVars) > 0 {
		indexHTML, indexModified = injectEnvVars(indexHTML, cfg.EnvVars)
	}

	// Create file server
	fileServer := http.FileServer(http.FS(fsys))

	// Strip base path if needed
	if cfg.BasePath != "/" {
		fileServer = http.StripPrefix(cfg.BasePath, fileServer)
	}

	return &Handler{
		config:        cfg,
		fileServer:    fileServer,
		indexHTML:     indexHTML,
		indexModified: indexModified,
	}, nil
}

// ServeHTTP implements http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Clean the path
	urlPath := r.URL.Path
	if h.config.BasePath != "/" {
		urlPath = strings.TrimPrefix(urlPath, strings.TrimSuffix(h.config.BasePath, "/"))
	}
	if urlPath == "" {
		urlPath = "/"
	}

	// Serve index.html for root
	if urlPath == "/" || urlPath == "/index.html" {
		h.serveIndex(w, r)
		return
	}

	// Check if file exists
	cleanPath := strings.TrimPrefix(urlPath, "/")
	var fsys fs.FS
	if h.config.FileSystem != nil {
		fsys = h.config.FileSystem
	} else {
		fsys = os.DirFS(h.config.BuildDir)
	}

	fileInfo, err := fs.Stat(fsys, cleanPath)

	if err == nil && !fileInfo.IsDir() {
		// File exists - serve it with caching headers
		h.serveStaticFile(w, r, cleanPath)
		return
	}

	// File doesn't exist - check if it's an asset request
	if isAssetPath(urlPath) {
		// This looks like an asset request but file doesn't exist - 404
		if h.config.NotFoundHandler != nil {
			h.config.NotFoundHandler.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
		return
	}

	// Not an asset path - serve index.html for SPA routing
	h.serveIndex(w, r)
}

// serveIndex serves the index.html file with no caching.
func (h *Handler) serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.indexHTML)
}

// serveStaticFile serves a static file with appropriate caching headers.
func (h *Handler) serveStaticFile(w http.ResponseWriter, r *http.Request, filePath string) {
	// Set content type
	ext := filepath.Ext(filePath)
	if contentType := mime.TypeByExtension(ext); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	// Set caching headers for static assets
	// Files with hashes in their names can be cached forever
	if hasContentHash(filePath) {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d, immutable", int(h.config.CacheMaxAge.Seconds())))
	} else {
		// Other static files get shorter cache
		w.Header().Set("Cache-Control", "public, max-age=3600")
	}

	// Let the file server handle the actual serving
	h.fileServer.ServeHTTP(w, r)
}

// isAssetPath determines if a path looks like a static asset.
func isAssetPath(urlPath string) bool {
	// Check for file extensions that indicate static assets
	ext := path.Ext(urlPath)
	staticExts := map[string]bool{
		".js":    true,
		".css":   true,
		".json":  true,
		".png":   true,
		".jpg":   true,
		".jpeg":  true,
		".gif":   true,
		".svg":   true,
		".ico":   true,
		".woff":  true,
		".woff2": true,
		".ttf":   true,
		".eot":   true,
		".map":   true,
		".txt":   true,
		".xml":   true,
	}

	return staticExts[ext]
}

// hasContentHash checks if a filename contains a content hash (e.g., main.abc123.js).
func hasContentHash(filePath string) bool {
	base := filepath.Base(filePath)
	// React/Webpack typically generate files like: main.abc123def.js or main.abc123def.chunk.js
	// Look for patterns with hex strings
	parts := strings.Split(base, ".")
	if len(parts) < 3 {
		return false
	}

	// Check if any part looks like a hash (8+ hex characters)
	for _, part := range parts[:len(parts)-1] {
		if len(part) >= 8 && isHexString(part) {
			return true
		}
	}

	return false
}

// isHexString checks if a string contains only hexadecimal characters.
func isHexString(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// injectEnvVars injects environment variables into the HTML content.
// Returns the modified HTML and a boolean indicating if modifications were made.
func injectEnvVars(html []byte, envVars map[string]string) ([]byte, bool) {
	if len(envVars) == 0 {
		return html, false
	}

	// Build the environment script
	var scriptBuilder strings.Builder
	scriptBuilder.WriteString("<script>window.__REACT_ENV__ = {")

	first := true
	for key, value := range envVars {
		if !first {
			scriptBuilder.WriteString(",")
		}
		first = false
		// Escape the value for JavaScript
		escapedValue := strings.ReplaceAll(value, "\\", "\\\\")
		escapedValue = strings.ReplaceAll(escapedValue, "\"", "\\\"")
		escapedValue = strings.ReplaceAll(escapedValue, "\n", "\\n")
		scriptBuilder.WriteString(fmt.Sprintf("%q:%q", key, escapedValue))
	}

	scriptBuilder.WriteString("};</script>")
	envScript := scriptBuilder.String()

	// Try to inject before </head>
	headTag := []byte("</head>")
	if idx := bytes.Index(html, headTag); idx != -1 {
		var result bytes.Buffer
		result.Write(html[:idx])
		result.WriteString(envScript)
		result.Write(html[idx:])
		return result.Bytes(), true
	}

	// Fallback: inject before </body>
	bodyTag := []byte("</body>")
	if idx := bytes.Index(html, bodyTag); idx != -1 {
		var result bytes.Buffer
		result.Write(html[:idx])
		result.WriteString(envScript)
		result.Write(html[idx:])
		return result.Bytes(), true
	}

	// If no head or body tags, prepend to the content
	var result bytes.Buffer
	result.WriteString(envScript)
	result.Write(html)
	return result.Bytes(), true
}

// Middleware returns a middleware that can be used with the server package.
// It serves the React app at the specified base path.
func (h *Handler) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the request path matches our base path
			if strings.HasPrefix(r.URL.Path, h.config.BasePath) {
				h.ServeHTTP(w, r)
				return
			}
			// Otherwise, pass to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// GetConfig returns a copy of the handler's configuration.
func (h *Handler) GetConfig() Config {
	return h.config
}
