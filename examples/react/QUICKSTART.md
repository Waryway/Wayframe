# Quick Start Guide - Full-Stack React + Go Example

## ✅ Everything is Ready!

This guide will get you running the full-stack React + Go application in under 5 minutes.

**Note:** This example uses **pure Go database drivers** (no CGO required), making it easy to build and cross-compile!

## Option 1: SQLite (Quickest - Default)

### Step 1: Run the Server
```bash
cd examples/react
go run main.go
```

That's it! The server creates a `wayframe.db` SQLite database automatically.

### Step 2: Open Browser
Visit: **http://localhost:8080**

You should see the React application with 5 example items loaded from the database!

---

## Option 2: PostgreSQL (Production-Ready)

### Step 1: Start PostgreSQL Database
```bash
docker run --name wayframe-db \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=wayframe \
  -p 5432:5432 \
  -d postgres:16
```

### Step 2: Set Environment Variables
```bash
export DB_DRIVER=postgres
export DB_DSN="host=localhost port=5432 user=postgres password=secret dbname=wayframe sslmode=disable"
```

### Step 3: Run the Server
```bash
cd examples/react
go run main.go
```

### Step 4: Open Browser
Visit: **http://localhost:8080**

You should see the React application with 5 example items loaded from the database!

---

## Option 3: MySQL (Alternative Production)

```bash
# Start MySQL
docker run --name wayframe-mysql -e MYSQL_ROOT_PASSWORD=secret -e MYSQL_DATABASE=wayframe -p 3306:3306 -d mysql:8

# Configure
export DB_DRIVER=mysql
export DB_DSN="root:secret@tcp(localhost:3306)/wayframe?parseTime=true"

# Run
cd examples/react && go run main.go
```

---

## What You'll See

The application includes:

1. **Item List Page** (`/`)
   - Shows 5 example items
   - Pagination controls
   - Delete button for each item
   - "Create New Item" button

2. **Item Detail Page** (`/items/:id`)
   - View item details
   - Edit inline
   - Delete item
   - Timestamps

3. **Create Item Page** (`/create`)
   - Form with name and description
   - Validation
   - Creates new item

## Testing the API

### Health Check
```bash
curl http://localhost:8080/api/health
```

Expected response:
```json
{
  "status": "ok",
  "database": "ok",
  "timestamp": "2024-12-23T..."
}
```

### List Items
```bash
curl http://localhost:8080/api/items
```

### Get Single Item
```bash
curl http://localhost:8080/api/items/1
```

### Create Item
```bash
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{"name":"My Item","description":"Item description"}'
```

### Update Item
```bash
curl -X PUT http://localhost:8080/api/items/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Name"}'
```

### Delete Item
```bash
curl -X DELETE http://localhost:8080/api/items/1
```

---

## Development Mode

### Frontend with Hot Reload
```bash
cd examples/react/frontend
npm install
npm run dev
```

This starts the webpack dev server on **http://localhost:3000** with:
- Hot module replacement
- Fast refresh
- Proxy to backend on port 8080

### Backend
```bash
cd examples/react
go run main.go
```

Now you can edit React components and see changes instantly!

---

## Building for Production

### Step 1: Build Frontend
```bash
cd examples/react/frontend
npm run build
```

This creates optimized bundles in `dist/` directory.

### Step 2: Build with Bazel
```bash
bazel build //examples/react:react
```

### Step 3: Run the Binary
```bash
bazel run //examples/react:react
```

Or directly:
```bash
./bazel-bin/examples/react/react_/react.exe
```

The binary includes everything:
- ✅ Go server code
- ✅ React frontend (embedded)
- ✅ Database migrations (embedded)

---

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_ADDR` | Server address | `:8080` |
| `DB_DRIVER` | Database driver | `sqlite` |
| `DB_DSN` | Connection string | `./wayframe.db` |
| `CORS_ORIGINS` | Allowed origins | `*` |

### Example: MySQL
```bash
export DB_DRIVER=mysql
export DB_DSN="root:password@tcp(localhost:3306)/wayframe?parseTime=true"
go run main.go
```

---

## Troubleshooting

### Port Already in Use
```bash
export SERVER_ADDR=:3000
go run main.go
```

### Database Connection Error
1. Check `examples/react/dist/` exists and contains files
2. Rebuild frontend: `cd frontend && npm run build`
3. Restart the server

### Database Connection Error
1. Verify database is running: `docker ps`
2. Check connection string format
3. Test connection: `psql -h localhost -U postgres -d wayframe`

---

## Next Steps

Now that your application is running:

1. **Explore the Code**
   - Check `examples/react/main.go` for API implementation
   - Review `frontend/src/` for React components
   - Look at `pkg/database/` for database utilities

2. **Customize**
   - Add new fields to the Item model
   - Create new API endpoints
   - Add authentication
   - Style the UI

3. **Deploy**
   - Build single binary with Bazel
   - Deploy to your server
   - Configure with environment variables

4. **Read Documentation**
   - See `examples/react/README.md` for comprehensive docs
   - Check `examples/react/IMPLEMENTATION_SUMMARY.md` for technical details

---

## Success!

Your full-stack React + Go application is now running! 🎉

- **Frontend**: Modern React 18 + TypeScript
- **Backend**: Go 1.25 REST API
- **Database**: PostgreSQL/MySQL/SQLite with migrations
- **CORS**: Properly configured
- **Production Ready**: Single binary deployment

**Happy coding!** 🚀

