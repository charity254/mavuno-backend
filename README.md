#  Mavuno Backend

REST API for Mavuno: A comprehensive farm produce tracking and marketplace platform designed for Kenyan farmers. Features production tracking, supply chain management, analytics, and a peer-to-peer marketplace.

## Tech Stack
- **Language:** Go 1.24.3
- **Router:** Gorilla Mux
- **Database:** PostgreSQL (Supabase)
- **Authentication:** JWT (72-hour expiration)
- **Password Hashing:** bcrypt
- **Environment:** godotenv
- **CORS:** rs/cors
- **Rate Limiting:** Custom middleware

## Getting Started

### Prerequisites
- Go 1.24.3+
- PostgreSQL database
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
DB_URL=postgresql_connection_string
JWT_SECRET=jwt_secret
PORT=8080
```

Generate a secure JWT secret:
```bash
openssl rand -base64 32
```

### Run Migrations
Run each migration in order against your database:
```bash
psql "db_url" -f migrations/001_create_users.sql
psql "db_url" -f migrations/002_create_products.sql
psql "db_url" -f migrations/003_create_produce_entries.sql
psql "db_url" -f migrations/004_create_supply_locations.sql
psql "db_url" -f migrations/005_create_supply_agreements.sql
psql "db_url" -f migrations/006_create_listings.sql
psql "db_url" -f migrations/007_create_messages.sql
psql "db_url" -f migrations/008_create_articles.sql
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
├── Dockerfile          — multi-stage Docker build
└── .env.example        — environment variable template


## User Roles

| Role | Permissions |
|------|-------------|
| **Farmer** | Produce tracking, supply management, analytics, reports, marketplace listings, articles |
| **Buyer** | Browse marketplace listings, message farmers, read articles |


## API Documentation

## API Modules

| Module | Description |
|--------|-------------|
| Auth | Register and login |
| Products | Farmer product catalog |
| Produce Entries | Daily produce movement tracking |
| Supply Locations | Delivery destination management |
| Supply Agreements | Recurring supply commitments |
| Analytics | Revenue, stock, rejection and planned vs actual charts |
| Reports | JSON view and CSV download of produce records |
| Marketplace | Listings and buyer-farmer messaging |
| Articles | Educational farming articles |
For detailed API endpoint documentation, including request/response examples, authentication details, and error handling, see [API.md](API.md).

Key points:
- All endpoints except `/api/auth/register` and `/api/auth/login` require a JWT token
- Tokens expire after 72 hours
- Role-based access control: Farmers have full access, Buyers have limited access

---

## Security Features
- JWT authentication on all protected routes
- Role-based access control — farmer vs buyer
- bcrypt password hashing
- Rate limiting — 100 requests per minute per IP
- Request body size limit — 1MB max
- CORS configured for frontend domain
- Soft deletes — data is never permanently removed
- Version conflict detection on all updates

---

## Deployment

### Docker
```bash
docker build -t mavuno-backend .
docker run -p 8080:8080 --env-file .env mavuno-backend
```

### Render
The backend is deployed on Render using the included Dockerfile. Environment variables are configured in the Render dashboard.

---
## Database Migrations

| File | Description |
|------|-------------|
| `001_create_users.sql` | Users table with role and active status |
| `002_create_products.sql` | Products with soft delete and version control |
| `003_create_produce_entries.sql` | Daily produce entries with stock constraints |
| `004_create_supply_locations.sql` | Supply delivery locations |
| `005_create_supply_agreements.sql` | Recurring supply agreements with delivery days |
| `006_create_listings.sql` | Marketplace listings |
| `007_create_messages.sql` | Buyer-farmer messaging |
| `008_create_articles.sql` | Educational articles with category filtering |

---
## License
MIT