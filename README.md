#  Mavuno Backend

REST API for Mavuno: A farm produce tracking and marketplace platform built for Kenyan farmers.

## Tech Stack
- **Language:** Go
- **Router:** Gorilla Mux
- **Database:** PostgreSQL (Supabase)
- **Authentication:** JWT
- **Password Hashing:** bcrypt

## Getting Started

### Prerequisites
- Go 1.22+
- PostgreSQL database (we use Supabase)

### Installation
```bash
git clone https://github.com/charity254/mavuno-backend.git
cd mavuno-backend
go mod download
```

### Environment Variables
Create a `.env` file in the root directory:
```
DB_URL=your_supabase_connection_string
JWT_SECRET=your_jwt_secret
PORT=8080
```

### Run the Server
```bash
go run cmd/server/main.go
```

### Health Check
```
GET http://localhost:8080/health
```

## Project Structure
- `cmd/server` — application entry point
- `internal/api` — route handlers
- `internal/middleware` — authentication and logging
- `internal/models` — data structures
- `internal/services` — business logic
- `internal/storage` — database queries
- `internal/utils` — shared helpers
- `migrations` — SQL database migrations

## License
MIT