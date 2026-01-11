# Full-Stack React + Go Example

This is a complete, production-ready example of a full-stack application using React for the frontend and Go with the Wayframe framework for the backend. It demonstrates a CRUD API with database integration, CORS support, and proper separation of concerns.

## Features

- **React 18** frontend with TypeScript and React Router
- **Go 1.25** backend with RESTful API
- **sqlx** for database operations (SQLite by default, supports PostgreSQL and MySQL)
- **CORS** middleware for cross-origin requests
- **Database migrations** with automatic execution
- **Embedded static files** for single-binary deployment
- **SPA routing** with fallback to index.html
- **Comprehensive error handling** and validation
- **Pagination** support for list endpoints
- **Health check** endpoint for monitoring

## Prerequisites

- **Bazel 8.0.0+** - For building the application
- **Go 1.25+** - For development
- **Node.js 18+** - For frontend development
- **npm or pnpm** - For managing frontend dependencies

## Project Structure

```
examples/react/
├── main.go                    # Go backend server
├── BUILD.bazel                # Bazel build configuration
├── config.example.json        # Example configuration file
├── frontend/                  # React frontend
│   ├── src/
│   │   ├── index.tsx         # React app entry point
│   │   ├── App.tsx           # Main application component
│   │   ├── components/       # React components
│   │   │   ├── ItemList.tsx
│   │   │   ├── ItemDetail.tsx
│   │   │   └── CreateItem.tsx
│   │   ├── api/
│   │   │   └── client.ts     # API client with fetch
│   │   ├── types/
│   │   │   └── api.ts        # TypeScript type definitions
│   │   └── styles.css        # Application styles
│   ├── public/
│   │   └── index.html        # HTML template
│   ├── package.json
│   ├── tsconfig.json
│   └── webpack.config.js
├── dist/                      # Built frontend (embedded in binary)
└── migrations/                # Database migrations
    ├── 001_create_items_table.sql
    └── 002_seed_data.sql
```

## Getting Started

### Database Setup

The example uses **pure Go database drivers** - no CGO required! This means:
- ✅ Fast compilation
- ✅ Easy cross-compilation
- ✅ No C compiler needed
- ✅ Works with Bazel out of the box

**Supported databases:**
- **SQLite** (pure Go via modernc.org/sqlite) - Default for development
- **PostgreSQL** - Recommended for production
- **MySQL** - Alternative for production

### Option 1: Quick Start with SQLite (Default)

```bash
# Build frontend
cd examples/react/frontend
npm install && npm run build

# Run server (SQLite by default)
cd ..
go run main.go
```

The server will start on `http://localhost:8080` and create a `wayframe.db` SQLite file automatically.

### Option 2: Using PostgreSQL with Docker (Production)

```bash
# Start PostgreSQL in Docker
docker run --name wayframe-postgres -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=wayframe -p 5432:5432 -d postgres:16

# Set environment variables
export DB_DRIVER=postgres
export DB_DSN="host=localhost port=5432 user=postgres password=secret dbname=wayframe sslmode=disable"

# Build and run
cd examples/react/frontend
npm install && npm run build
cd ..
go run main.go
```


### Option 3: Using MySQL with Docker

```bash
# Start MySQL in Docker  
docker run --name wayframe-mysql -e MYSQL_ROOT_PASSWORD=secret -e MYSQL_DATABASE=wayframe -p 3306:3306 -d mysql:8

# Set environment variables
export DB_DRIVER=mysql
export DB_DSN="root:secret@tcp(localhost:3306)/wayframe?parseTime=true"

# Build and run
go run main.go
```

### 1. Build the Frontend

```bash
cd examples/react/frontend
npm install
npm run build
```

This creates the production build in `examples/react/dist/`.

### 2. Build and Run with Bazel

```bash
# Build everything
bazel build //examples/react:react

# Run the server
bazel run //examples/react:react
```

The server will start on `http://localhost:8080` by default.

### 3. Build and Run with Go

```bash
# From the project root
cd examples/react
go run main.go
```

## Development Workflow

### Frontend Development

For faster frontend development with hot reload:

```bash
cd examples/react/frontend
npm run dev
```

This starts the webpack dev server on `http://localhost:3000` with:
- Hot module replacement
- Proxy to backend API on `http://localhost:8080`

Make sure the backend server is running separately.

### Backend Development

Run the Go server directly for faster iteration:

```bash
cd examples/react
go run main.go
```

## Configuration

The application can be configured via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_ADDR` | Server address and port | `:8080` |
| `DB_DRIVER` | Database driver (sqlite, postgres, mysql) | `sqlite` |
| `DB_DSN` | Database connection string | `./wayframe.db` |
| `CORS_ORIGINS` | Comma-separated allowed origins | `*` |

### Example Configurations

**SQLite (default) - Pure Go, no CGO**
```bash
export DB_DRIVER=sqlite
export DB_DSN=./wayframe.db
```

**PostgreSQL**
```bash
export DB_DRIVER=postgres
export DB_DSN="host=localhost port=5432 user=postgres password=secret dbname=wayframe sslmode=disable"
```

**MySQL**
```bash
export DB_DRIVER=mysql
export DB_DSN="user:password@tcp(localhost:3306)/wayframe?parseTime=true"
```

**CORS for specific origins**
```bash
export CORS_ORIGINS="http://localhost:3000,https://example.com"
```

## API Endpoints

### Health Check

```bash
GET /api/health
```

**Response:**
```json
{
  "status": "ok",
  "database": "ok",
  "timestamp": "2024-12-23T10:30:00Z"
}
```

### List Items

```bash
GET /api/items?page=1&per_page=10
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Example Item",
      "description": "Item description",
      "created_at": "2024-12-23T10:00:00Z",
      "updated_at": "2024-12-23T10:00:00Z"
    }
  ],
  "page": 1,
  "per_page": 10,
  "total": 5
}
```

### Get Single Item

```bash
GET /api/items/{id}
```

**Response:**
```json
{
  "id": 1,
  "name": "Example Item",
  "description": "Item description",
  "created_at": "2024-12-23T10:00:00Z",
  "updated_at": "2024-12-23T10:00:00Z"
}
```

### Create Item

```bash
POST /api/items
Content-Type: application/json

{
  "name": "New Item",
  "description": "Item description"
}
```

**Response:** `201 Created` with the created item

### Update Item

```bash
PUT /api/items/{id}
Content-Type: application/json

{
  "name": "Updated Name",
  "description": "Updated description"
}
```

**Response:** `200 OK` with the updated item

### Delete Item

```bash
DELETE /api/items/{id}
```

**Response:** `204 No Content`

## Testing the API

### Using curl

```bash
# Health check
curl http://localhost:8080/api/health

# List items
curl http://localhost:8080/api/items

# Get single item
curl http://localhost:8080/api/items/1

# Create item
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Item","description":"Test description"}'

# Update item
curl -X PUT http://localhost:8080/api/items/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Name"}'

# Delete item
curl -X DELETE http://localhost:8080/api/items/1
```

## Database Migrations

Migrations are automatically run on startup. They are embedded in the binary and stored in the `migrations/` directory.

### Migration Files

- `001_create_items_table.sql` - Creates the items table with auto-incrementing ID
- `002_seed_data.sql` - Adds example data for testing

The migration system:
- Tracks applied migrations in a `migrations` table
- Runs migrations in order by version number
- Skips already-applied migrations
- Runs each migration in a transaction
- Automatically creates/updates the `updated_at` timestamp

### Adding New Migrations

1. Create a new SQL file: `migrations/003_your_migration.sql`
2. Add the migration to `main.go`:

```go
migrations := []database.Migration{
    // ...existing migrations...
    {
        Version:     3,
        Description: "Your migration description",
        Up:          mustReadMigration("migrations/003_your_migration.sql"),
    },
}
```

## Production Build

### Single Binary

The application embeds all static assets and can be deployed as a single binary:

```bash
bazel build //examples/react:react --compilation_mode=opt
```

The binary will include:
- Go backend code
- React frontend assets (from `dist/`)
- Database migrations

### Docker

Example Dockerfile:

```dockerfile
FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN go build -o server examples/react/main.go

FROM debian:bookworm-slim
COPY --from=builder /app/server /server
EXPOSE 8080
CMD ["/server"]
```

## Troubleshooting

### Frontend not loading

1. Ensure you've built the frontend: `cd frontend && npm run build`
2. Check that `dist/` directory exists and contains `index.html`
3. Verify the embed directives in `main.go`

### Database errors

1. Check database permissions
2. Verify DSN format for your database driver
3. Ensure database server is running (for PostgreSQL/MySQL)
4. Check migration SQL syntax for your database

### CORS errors

1. Set `CORS_ORIGINS` to include your frontend URL
2. Check browser console for specific CORS errors
3. Verify preflight OPTIONS requests are handled

### Port already in use

Change the server port:
```bash
export SERVER_ADDR=:3000
```

## Architecture

### Backend (Go)

- **Framework**: Wayframe with stdlib HTTP server
- **Database**: sqlx with support for SQLite, PostgreSQL, MySQL
- **Routing**: Standard library with pattern matching (Go 1.22+)
- **Middleware**: CORS, logging, recovery
- **Static Files**: Embedded with `embed.FS`

### Frontend (React)

- **Framework**: React 18 with TypeScript
- **Routing**: React Router v6
- **Build Tool**: Webpack 5
- **State Management**: React hooks (useState, useEffect)
- **Styling**: Plain CSS with modern features

## Contributing

This example demonstrates best practices for:
- Clean separation of concerns
- Proper error handling
- Type safety (Go and TypeScript)
- Database transactions
- Middleware patterns
- RESTful API design
- SPA routing
- Production-ready builds

Feel free to use this as a starting point for your own applications!

## License

This example is part of the Wayframe project and follows the same license.

