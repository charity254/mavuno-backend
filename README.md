#  Mavuno Backend

REST API for Mavuno: A comprehensive farm produce tracking and marketplace platform designed for Kenyan farmers. Features production tracking, supply chain management, analytics, and a peer-to-peer marketplace.

## Tech Stack
- **Language:** Go 1.24.3
- **Router:** Gorilla Mux
- **Database:** PostgreSQL (Supabase)
- **Authentication:** JWT (72-hour expiration)
- **Password Hashing:** bcrypt
- **Environment:** godotenv

## Getting Started

### Prerequisites
- Go 1.24.3+
- PostgreSQL database (we use Supabase)
- Git

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

The server will start with: `🌾 Mavuno server starting on port 8080...`

### Health Check
```
GET http://localhost:8080/health
```

## Project Structure
- `cmd/server` — application entry point
- `internal/api` — HTTP route handlers and router setup
- `internal/config` — configuration management
- `internal/middleware` — JWT authentication and role-based access control
- `internal/models` — data structures for domain entities
- `internal/services` — business logic layer
- `internal/storage` — database repositories and queries
- `internal/utils` — shared utility functions
- `migrations` — SQL database migration scripts
- `docs` — API documentation

## API Documentation

For detailed API endpoint documentation, including request/response examples, authentication details, and error handling, see [API.md](API.md).

Key points:
- All endpoints except `/api/auth/register` and `/api/auth/login` require a JWT token
- Tokens expire after 72 hours
- Role-based access control: Farmers have full access, Buyers have limited access

## User Roles
- **Farmer:** Can track produce, manage supply agreements, access analytics, and post marketplace listings
- **Buyer:** Can browse marketplace listings and send messages to farmers

## Database Migrations
The following migrations are included:
- `001_create_users.sql` — User accounts with role-based access
- `002_create_products.sql` — Product catalog
- `003_create_produce_entries.sql` — Production records
- `004_create_supply_locations.sql` — Storage and delivery locations
- `005_create_supply_agreements.sql` — Supply contracts
- `006_create_listings.sql` — Marketplace listings
- `007_create_messages.sql` — Marketplace messaging

## License
MIT