package main

// gazelle:ignore

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Waryway/Wayframe/internal/web"
	"github.com/Waryway/Wayframe/internal/web/stdlib"
	"github.com/Waryway/Wayframe/pkg/database"
)

// Note: `staticFiles` is generated at build time by the
// `generate_embedded_dist` genrule (zz_embedded_dist.go) and provides
// an fs.FS-like value containing the frontend assets.
// The genrule produces a variable named `staticFiles` (type memFS) so
// we don't use go:embed here.

// Migrations are read from disk at runtime so they don't need to be
// listed in go_library srcs for compilation.

// Item represents an item in the database
type Item struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateItemRequest represents the request body for creating an item
type CreateItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateItemRequest represents the request body for updating an item
type UpdateItemRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data    interface{} `json:"data"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
	Total   int         `json:"total"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Database  string    `json:"database"`
	Timestamp time.Time `json:"timestamp"`
}

// App holds the application dependencies
type App struct {
	db     *database.DB
	logger *slog.Logger
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Get database configuration from environment
	dbDriver := getEnv("DB_DRIVER", "sqlite")
	dbDSN := getEnv("DB_DSN", "./wayframe.db")

	// Initialize database
	dbConfig := database.Config{
		Driver:          dbDriver,
		DSN:             dbDSN,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
	}

	db, err := database.New(dbConfig)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("Connected to database", "driver", dbDriver)

	// Run migrations
	if err := runMigrations(db); err != nil {
		logger.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	app := &App{
		db:     db,
		logger: logger,
	}

	// Create web server
	serverAddr := getEnv("SERVER_ADDR", ":8080")
	serverConfig := web.Config{
		Addr:         serverAddr,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	server := stdlib.New(serverConfig)

	// Configure CORS
	corsOrigins := strings.Split(getEnv("CORS_ORIGINS", "*"), ",")
	corsConfig := stdlib.CORSConfig{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: false,
		MaxAge:           3600,
	}

	// Add CORS middleware
	server.Use(stdlib.CORS(corsConfig))

	// Add logging middleware
	server.Use(loggingMiddleware(logger))

	// Add recovery middleware
	server.Use(recoveryMiddleware(logger))

	// Register API routes
	server.HandleFunc("GET /api/health", app.handleHealth)
	server.HandleFunc("GET /api/items", app.handleGetItems)
	server.HandleFunc("GET /api/items/{id}", app.handleGetItem)
	server.HandleFunc("POST /api/items", app.handleCreateItem)
	server.HandleFunc("PUT /api/items/{id}", app.handleUpdateItem)
	server.HandleFunc("DELETE /api/items/{id}", app.handleDeleteItem)

	// Serve static files. We resolve the asset FS dynamically per-request
	// so that a developer can rebuild the frontend into the `examples/react/dist`
	// directory and the running Go server will start serving the new files
	// without needing a restart. The resolver prefers WAYFRAME_ASSETS_DIR,
	// then project-local examples/react/dist, then falls back to the embedded
	// `staticFiles` generated at build time.
	resolveDistFS := func() (fs.FS, error) {
		assetsDir := os.Getenv("WAYFRAME_ASSETS_DIR")
		if assetsDir != "" {
			if _, err := os.Stat(assetsDir); err == nil {
				return os.DirFS(assetsDir), nil
			}
		}
		if _, err := os.Stat("examples/react/dist"); err == nil {
			return os.DirFS("examples/react/dist"), nil
		}
		// Fall back to embedded files
		return fs.Sub(staticFiles, "dist")
	}

	// Use a handler that resolves the FS on each request.
	server.Handle("/", spaHandler(func() (fs.FS, error) { return resolveDistFS() }))

	logger.Info("Starting server", "addr", serverAddr)

	// Start server
	if err := server.Start(30 * time.Second); err != nil {
		logger.Error("Server error", "error", err)
		os.Exit(1)
	}
}

// runMigrations runs database migrations
func runMigrations(db *database.DB) error {
	migrations := []database.Migration{
		{
			Version:     1,
			Description: "Create items table",
			Up:          mustReadMigration("migrations/001_create_items_table.sql"),
		},
		{
			Version:     2,
			Description: "Seed data",
			Up:          mustReadMigration("migrations/002_seed_data.sql"),
		},
	}

	return db.RunMigrations(context.Background(), migrations)
}

// mustReadMigration reads a migration file or panics
func mustReadMigration(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to read migration %s: %v", path, err))
	}
	return string(data)
}

// spaHandler handles SPA routing by serving index.html for non-API routes
func spaHandler(resolveFS func() (fs.FS, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't serve index.html for API routes
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Try to serve the file
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Resolve the filesystem for the request
		distFS, err := resolveFS()
		if err != nil {
			http.Error(w, "Failed to resolve filesystem", http.StatusInternalServerError)
			return
		}

		// Check if file exists
		if _, err := fs.Stat(distFS, path); err == nil {
			http.FileServer(http.FS(distFS)).ServeHTTP(w, r)
			return
		}

		// Serve index.html for client-side routing
		r.URL.Path = "/"
		http.FileServer(http.FS(distFS)).ServeHTTP(w, r)
	})
}

// handleHealth handles health check requests
func (app *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	dbStatus := "ok"
	if err := app.db.HealthCheck(r.Context()); err != nil {
		dbStatus = "error: " + err.Error()
	}

	response := HealthResponse{
		Status:    "ok",
		Database:  dbStatus,
		Timestamp: time.Now(),
	}

	respondJSON(w, http.StatusOK, response)
}

// handleGetItems handles GET /api/items
func (app *App) handleGetItems(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// Get total count
	var total int
	err := app.db.Get(&total, "SELECT COUNT(*) FROM items")
	if err != nil {
		app.logger.Error("Failed to count items", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to fetch items")
		return
	}

	// Get items
	var items []Item
	err = app.db.Select(&items, "SELECT * FROM items ORDER BY created_at DESC LIMIT ? OFFSET ?", perPage, offset)
	if err != nil {
		app.logger.Error("Failed to fetch items", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to fetch items")
		return
	}

	if items == nil {
		items = []Item{}
	}

	response := PaginatedResponse{
		Data:    items,
		Page:    page,
		PerPage: perPage,
		Total:   total,
	}

	respondJSON(w, http.StatusOK, response)
}

// handleGetItem handles GET /api/items/{id}
func (app *App) handleGetItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var item Item
	err = app.db.Get(&item, "SELECT * FROM items WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "Item not found")
			return
		}
		app.logger.Error("Failed to fetch item", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "Failed to fetch item")
		return
	}

	respondJSON(w, http.StatusOK, item)
}

// handleCreateItem handles POST /api/items
func (app *App) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	var req CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.Description == "" {
		respondError(w, http.StatusBadRequest, "Name and description are required")
		return
	}

	result, err := app.db.Exec(
		"INSERT INTO items (name, description) VALUES (?, ?)",
		req.Name, req.Description,
	)
	if err != nil {
		app.logger.Error("Failed to create item", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	id, _ := result.LastInsertId()

	var item Item
	err = app.db.Get(&item, "SELECT * FROM items WHERE id = ?", id)
	if err != nil {
		app.logger.Error("Failed to fetch created item", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "Failed to fetch created item")
		return
	}

	respondJSON(w, http.StatusCreated, item)
}

// handleUpdateItem handles PUT /api/items/{id}
func (app *App) handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}

	if req.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *req.Description)
	}

	if len(updates) == 0 {
		respondError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	query := fmt.Sprintf("UPDATE items SET %s WHERE id = ?", strings.Join(updates, ", "))

	_, err = app.db.Exec(query, args...)
	if err != nil {
		app.logger.Error("Failed to update item", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "Failed to update item")
		return
	}

	var item Item
	err = app.db.Get(&item, "SELECT * FROM items WHERE id = ?", id)
	if err != nil {
		app.logger.Error("Failed to fetch updated item", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "Failed to fetch updated item")
		return
	}

	respondJSON(w, http.StatusOK, item)
}

// handleDeleteItem handles DELETE /api/items/{id}
func (app *App) handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	result, err := app.db.Exec("DELETE FROM items WHERE id = ?", id)
	if err != nil {
		app.logger.Error("Failed to delete item", "error", err, "id", id)
		respondError(w, http.StatusInternalServerError, "Failed to delete item")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondError(w, http.StatusNotFound, "Item not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("Request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
			)
		})
	}
}

// recoveryMiddleware recovers from panics
func recoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered", "error", err)
					respondError(w, http.StatusInternalServerError, "Internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
