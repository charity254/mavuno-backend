# Mavuno API Documentation

Base URL: `http://localhost:8080` (development)

All protected routes require a JWT token in the Authorization header.
The token is obtained from the Login endpoint and expires after 72 hours.

**How to use:**
1. Call `POST /api/auth/login` to get your token
2. Copy the token value from the response
3. Attach it to every protected request like this:
```
Authorization: Bearer <your_token_here>
```

---

## Authentication (Public)

### Register
`POST /api/auth/register`

No token required.

**Request Body:**
```json
{
    "email": "farmer@mavuno.com",
    "password": "password123",
    "full_name": "John Kamau",
    "role": "farmer"
}
```

**Notes:**
- `role` must be either `farmer` or `buyer`
- `password` must be at least 8 characters

**Success Response — `201 Created`:**
```json
{
    "message": "account created successfully"
}
```

**Error Responses:**
| Status | Message |
|--------|---------|
| `400` | `password must be at least 8 characters` |
| `400` | `role must be either farmer or buyer` |
| `409` | `user with this email already exists` |

---

### Login
`POST /api/auth/login`

No token required.

**Request Body:**
```json
{
    "email": "farmer@mavuno.com",
    "password": "password123"
}
```

**Success Response — `200 OK`:**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "email": "farmer@mavuno.com",
        "full_name": "John Kamau",
        "role": "farmer"
    }
}
```

**Error Responses:**
| Status | Message |
|--------|---------|
| `401` | `invalid credentials` |
| `404` | `user not found` |

---

## Products (Farmer only)

### Create Product
`POST /api/products`

**Required:** JWT token + farmer role

**Request Body:**
```json
{
    "name": "Tomatoes",
    "variety": "Roma",
    "unit": "kg",
    "expected_yield": 500,
    "planting_date": "2026-01-15",
    "expected_harvest_date": "2026-03-15"
}
```

**Success Response — `201 Created`:**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Tomatoes",
    "variety": "Roma",
    "unit": "kg",
    "expected_yield": 500,
    "planting_date": "2026-01-15",
    "expected_harvest_date": "2026-03-15",
    "created_at": "2026-05-02T10:30:00Z"
}
```

---

### Get All Products
`GET /api/products`

**Required:** JWT token + farmer role

**Success Response — `200 OK`:**
```json
[
    {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Tomatoes",
        "variety": "Roma",
        "unit": "kg",
        "expected_yield": 500,
        "planting_date": "2026-01-15",
        "expected_harvest_date": "2026-03-15",
        "created_at": "2026-05-02T10:30:00Z"
    }
]
```

---

### Get Product by ID
`GET /api/products/{id}`

**Required:** JWT token + farmer role

**Success Response — `200 OK`:**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Tomatoes",
    "variety": "Roma",
    "unit": "kg",
    "expected_yield": 500,
    "planting_date": "2026-01-15",
    "expected_harvest_date": "2026-03-15",
    "created_at": "2026-05-02T10:30:00Z"
}
```

---

### Update Product
`PUT /api/products/{id}`

**Required:** JWT token + farmer role

**Request Body:**
```json
{
    "name": "Tomatoes",
    "variety": "Cherry",
    "unit": "kg",
    "expected_yield": 600,
    "expected_harvest_date": "2026-03-20"
}
```

**Success Response — `200 OK`:**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Tomatoes",
    "variety": "Cherry",
    "unit": "kg",
    "expected_yield": 600,
    "planting_date": "2026-01-15",
    "expected_harvest_date": "2026-03-20",
    "updated_at": "2026-05-02T11:00:00Z"
}
```

---

### Delete Product
`DELETE /api/products/{id}`

**Required:** JWT token + farmer role

**Success Response — `204 No Content`**

---

## Produce Entries (Farmer only)

### Create Entry
`POST /api/entries`

**Required:** JWT token + farmer role

**Request Body:**
```json
{
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "quantity_produced": 450,
    "quantity_rejected": 50,
    "entry_date": "2026-03-15",
    "notes": "Good quality harvest"
}
```

**Success Response — `201 Created`:**
```json
{
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "quantity_produced": 450,
    "quantity_rejected": 50,
    "entry_date": "2026-03-15",
    "notes": "Good quality harvest",
    "created_at": "2026-05-02T10:30:00Z"
}
```

---

### Get All Entries
`GET /api/entries`

**Required:** JWT token + farmer role

**Query Parameters (optional):**
- `product_id` — Filter by product ID
- `start_date` — Filter from date (YYYY-MM-DD)
- `end_date` — Filter to date (YYYY-MM-DD)

**Success Response — `200 OK`:**
```json
[
    {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "product_id": "550e8400-e29b-41d4-a716-446655440000",
        "quantity_produced": 450,
        "quantity_rejected": 50,
        "entry_date": "2026-03-15",
        "notes": "Good quality harvest",
        "created_at": "2026-05-02T10:30:00Z"
    }
]
```

---

### Get Entry by ID
`GET /api/entries/{id}`

**Required:** JWT token + farmer role

**Success Response — `200 OK`:**
```json
{
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "quantity_produced": 450,
    "quantity_rejected": 50,
    "entry_date": "2026-03-15",
    "notes": "Good quality harvest",
    "created_at": "2026-05-02T10:30:00Z"
}
```

---

### Update Entry
`PUT /api/entries/{id}`

**Required:** JWT token + farmer role

**Request Body:**
```json
{
    "quantity_produced": 460,
    "quantity_rejected": 40,
    "notes": "Updated quality harvest"
}
```

**Success Response — `200 OK`:**
```json
{
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "quantity_produced": 460,
    "quantity_rejected": 40,
    "entry_date": "2026-03-15",
    "notes": "Updated quality harvest",
    "updated_at": "2026-05-02T11:00:00Z"
}
```

---

### Delete Entry
`DELETE /api/entries/{id}`

**Required:** JWT token + farmer role

**Success Response — `204 No Content`**

---

## Supply Locations (Farmer only)

### Create Supply Location
`POST /api/supply-locations`

**Required:** JWT token + farmer role

**Request Body:**
```json
{
    "name": "Main Storage",
    "address": "123 Farm Road, Nairobi",
    "latitude": -1.2864,
    "longitude": 36.8172,
    "capacity": 1000,
    "type": "warehouse"
}
```

**Success Response — `201 Created`:**
```json
{
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "name": "Main Storage",
    "address": "123 Farm Road, Nairobi",
    "latitude": -1.2864,
    "longitude": 36.8172,
    "capacity": 1000,
    "type": "warehouse",
    "created_at": "2026-05-02T10:30:00Z"
}
```

---

### Get All Supply Locations
`GET /api/supply-locations`

**Required:** JWT token + farmer role

**Success Response — `200 OK`:**
```json
[
    {
        "id": "770e8400-e29b-41d4-a716-446655440002",
        "name": "Main Storage",
        "address": "123 Farm Road, Nairobi",
        "latitude": -1.2864,
        "longitude": 36.8172,
        "capacity": 1000,
        "type": "warehouse",
        "created_at": "2026-05-02T10:30:00Z"
    }
]
```

---

### Get Supply Location by ID
`GET /api/supply-locations/{id}`

**Required:** JWT token + farmer role

**Success Response — `200 OK`:**
```json
{
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "name": "Main Storage",
    "address": "123 Farm Road, Nairobi",
    "latitude": -1.2864,
    "longitude": 36.8172,
    "capacity": 1000,
    "type": "warehouse",
    "created_at": "2026-05-02T10:30:00Z"
}
```

---

### Update Supply Location
`PUT /api/supply-locations/{id}`

**Required:** JWT token + farmer role

**Request Body:**
```json
{
    "name": "Main Storage Updated",
    "capacity": 1200
}
```

**Success Response — `200 OK`:**
```json
{
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "name": "Main Storage Updated",
    "address": "123 Farm Road, Nairobi",
    "latitude": -1.2864,
    "longitude": 36.8172,
    "capacity": 1200,
    "type": "warehouse",
    "updated_at": "2026-05-02T11:00:00Z"
}
```

---

### Delete Supply Location
`DELETE /api/supply-locations/{id}`

**Required:** JWT token + farmer role

**Success Response — `204 No Content`**

---

## Supply Agreements (Farmer only)

### Create Supply Agreement
`POST /api/supply-agreements`

**Required:** JWT token + farmer role

**Request Body:**
```json
{
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "supply_location_id": "770e8400-e29b-41d4-a716-446655440002",
    "buyer_name": "Fresh Foods Ltd",
    "quantity_committed": 500,
    "unit": "kg",
    "price_per_unit": 50,
    "start_date": "2026-05-01",
    "end_date": "2026-08-31",
    "delivery_frequency": "weekly"
}
```

**Success Response — `201 Created`:**
```json
{
    "id": "880e8400-e29b-41d4-a716-446655440003",
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "supply_location_id": "770e8400-e29b-41d4-a716-446655440002",
    "buyer_name": "Fresh Foods Ltd",
    "quantity_committed": 500,
    "unit": "kg",
    "price_per_unit": 50,
    "start_date": "2026-05-01",
    "end_date": "2026-08-31",
    "delivery_frequency": "weekly",
    "status": "active",
    "created_at": "2026-05-02T10:30:00Z"
}
```

---

### Get All Supply Agreements
`GET /api/supply-agreements`

**Required:** JWT token + farmer role

**Success Response — `200 OK`:**
```json
[
    {
        "id": "880e8400-e29b-41d4-a716-446655440003",
        "product_id": "550e8400-e29b-41d4-a716-446655440000",
        "supply_location_id": "770e8400-e29b-41d4-a716-446655440002",
        "buyer_name": "Fresh Foods Ltd",
        "quantity_committed": 500,
        "unit": "kg",
        "price_per_unit": 50,
        "start_date": "2026-05-01",
        "end_date": "2026-08-31",
        "delivery_frequency": "weekly",
        "status": "active",
        "created_at": "2026-05-02T10:30:00Z"
    }
]
```

---

### Get Active Supply Agreements
`GET /api/supply-agreements/active`

**Required:** JWT token + farmer role

Returns only agreements with status "active" and within the date range.

**Success Response — `200 OK`:**
```json
[
    {
        "id": "880e8400-e29b-41d4-a716-446655440003",
        "product_id": "550e8400-e29b-41d4-a716-446655440000",
        "supply_location_id": "770e8400-e29b-41d4-a716-446655440002",
        "buyer_name": "Fresh Foods Ltd",
        "quantity_committed": 500,
        "unit": "kg",
        "price_per_unit": 50,
        "start_date": "2026-05-01",
        "end_date": "2026-08-31",
        "delivery_frequency": "weekly",
        "status": "active",
        "next_delivery": "2026-05-08",
        "created_at": "2026-05-02T10:30:00Z"
    }
]
```

---

### Get Supply Agreement by ID
`GET /api/supply-agreements/{id}`

**Required:** JWT token + farmer role

**Success Response — `200 OK`:**
```json
{
    "id": "880e8400-e29b-41d4-a716-446655440003",
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "supply_location_id": "770e8400-e29b-41d4-a716-446655440002",
    "buyer_name": "Fresh Foods Ltd",
    "quantity_committed": 500,
    "unit": "kg",
    "price_per_unit": 50,
    "start_date": "2026-05-01",
    "end_date": "2026-08-31",
    "delivery_frequency": "weekly",
    "status": "active",
    "created_at": "2026-05-02T10:30:00Z"
}
```

---

### Update Supply Agreement
`PUT /api/supply-agreements/{id}`

**Required:** JWT token + farmer role

**Request Body:**
```json
{
    "quantity_committed": 600,
    "price_per_unit": 55,
    "status": "completed"
}
```

**Success Response — `200 OK`:**
```json
{
    "id": "880e8400-e29b-41d4-a716-446655440003",
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "supply_location_id": "770e8400-e29b-41d4-a716-446655440002",
    "buyer_name": "Fresh Foods Ltd",
    "quantity_committed": 600,
    "unit": "kg",
    "price_per_unit": 55,
    "start_date": "2026-05-01",
    "end_date": "2026-08-31",
    "delivery_frequency": "weekly",
    "status": "completed",
    "updated_at": "2026-05-02T11:00:00Z"
}
```

---

### Delete Supply Agreement
`DELETE /api/supply-agreements/{id}`

**Required:** JWT token + farmer role

**Success Response — `204 No Content`**

---

## Analytics (Farmer only)

### Get Today's Summary
`GET /api/analytics/today`

**Required:** JWT token + farmer role

**Success Response — `200 OK`:**
```json
{
    "total_produced": 100,
    "total_rejected": 5,
    "acceptance_rate": 95,
    "total_entries": 2
}
```

---

### Get Revenue Trend
`GET /api/analytics/revenue`

**Required:** JWT token + farmer role

**Query Parameters (optional):**
- `days` — Number of days to analyze (default: 30)

**Success Response — `200 OK`:**
```json
{
    "period": "30 days",
    "total_revenue": 25000,
    "average_daily": 833.33,
    "trend_data": [
        {
            "date": "2026-04-02",
            "revenue": 500
        },
        {
            "date": "2026-04-03",
            "revenue": 750
        }
    ]
}
```

---

### Get Stock Trend
`GET /api/analytics/stock`

**Required:** JWT token + farmer role

**Query Parameters (optional):**
- `days` — Number of days to analyze (default: 30)

**Success Response — `200 OK`:**
```json
{
    "period": "30 days",
    "total_stock": 500,
    "trend_data": [
        {
            "date": "2026-04-02",
            "quantity": 100
        },
        {
            "date": "2026-04-03",
            "quantity": 150
        }
    ]
}
```

---

### Get Rejection Trend
`GET /api/analytics/rejected`

**Required:** JWT token + farmer role

**Query Parameters (optional):**
- `days` — Number of days to analyze (default: 30)

**Success Response — `200 OK`:**
```json
{
    "period": "30 days",
    "total_rejected": 50,
    "rejection_rate": 5.5,
    "trend_data": [
        {
            "date": "2026-04-02",
            "rejected": 5
        },
        {
            "date": "2026-04-03",
            "rejected": 8
        }
    ]
}
```

---

### Get Product Summary
`GET /api/analytics/product-summary`

**Required:** JWT token + farmer role

**Success Response — `200 OK`:**
```json
[
    {
        "product_id": "550e8400-e29b-41d4-a716-446655440000",
        "product_name": "Tomatoes",
        "total_produced": 500,
        "total_rejected": 25,
        "acceptance_rate": 95,
        "entries_count": 5
    },
    {
        "product_id": "550e8400-e29b-41d4-a716-446655440001",
        "product_name": "Peppers",
        "total_produced": 300,
        "total_rejected": 10,
        "acceptance_rate": 96.7,
        "entries_count": 3
    }
]
```

---

### Get Planned vs Actual
`GET /api/analytics/planned-vs-actual`

**Required:** JWT token + farmer role

**Success Response — `200 OK`:**
```json
[
    {
        "product_id": "550e8400-e29b-41d4-a716-446655440000",
        "product_name": "Tomatoes",
        "planned_yield": 1000,
        "actual_yield": 800,
        "variance": -200,
        "variance_percentage": -20
    },
    {
        "product_id": "550e8400-e29b-41d4-a716-446655440001",
        "product_name": "Peppers",
        "planned_yield": 600,
        "actual_yield": 650,
        "variance": 50,
        "variance_percentage": 8.33
    }
]
```

---

## Reports (Farmer only)

### Get Report
`GET /api/reports/export`

**Required:** JWT token + farmer role

**Query Parameters (optional):**
- `start_date` — Report start date (YYYY-MM-DD)
- `end_date` — Report end date (YYYY-MM-DD)

**Success Response — `200 OK`:**
```json
{
    "period": {
        "start": "2026-04-01",
        "end": "2026-05-02"
    },
    "summary": {
        "total_produced": 1200,
        "total_rejected": 60,
        "acceptance_rate": 95,
        "total_revenue": 60000
    },
    "products": [
        {
            "name": "Tomatoes",
            "produced": 500,
            "rejected": 25,
            "entries": 5
        }
    ],
    "entries": [
        {
            "date": "2026-04-15",
            "product": "Tomatoes",
            "quantity": 100,
            "rejected": 5
        }
    ]
}
```

---

### Download CSV
`GET /api/reports/download`

**Required:** JWT token + farmer role

**Query Parameters (optional):**
- `start_date` — Report start date (YYYY-MM-DD)
- `end_date` — Report end date (YYYY-MM-DD)

**Success Response — `200 OK` (CSV file)**

Returns a CSV file with headers:
```
Date,Product,Produced,Rejected,Quality Rate
2026-04-15,Tomatoes,100,5,95%
2026-04-16,Peppers,75,2,97%
```

---

## Marketplace - Listings

### Create Listing (Farmer only)
`POST /api/listings`

**Required:** JWT token + farmer role

**Request Body:**
```json
{
    "product_name": "Fresh Tomatoes",
    "quantity": 100,
    "unit": "kg",
    "price": 50,
    "description": "Roma tomatoes, freshly harvested",
    "location": "Nairobi, Kenya",
    "image_url": "https://example.com/image.jpg"
}
```

**Success Response — `201 Created`:**
```json
{
    "id": "990e8400-e29b-41d4-a716-446655440004",
    "farmer_id": "550e8400-e29b-41d4-a716-446655440000",
    "product_name": "Fresh Tomatoes",
    "quantity": 100,
    "unit": "kg",
    "price": 50,
    "description": "Roma tomatoes, freshly harvested",
    "location": "Nairobi, Kenya",
    "image_url": "https://example.com/image.jpg",
    "status": "active",
    "created_at": "2026-05-02T10:30:00Z"
}
```

---

### Get All Listings
`GET /api/listings`

**Required:** JWT token (any authenticated user)

**Query Parameters (optional):**
- `location` — Filter by location
- `product` — Filter by product name
- `min_price` — Minimum price
- `max_price` — Maximum price

**Success Response — `200 OK`:**
```json
[
    {
        "id": "990e8400-e29b-41d4-a716-446655440004",
        "farmer_id": "550e8400-e29b-41d4-a716-446655440000",
        "farmer_name": "John Kamau",
        "product_name": "Fresh Tomatoes",
        "quantity": 100,
        "unit": "kg",
        "price": 50,
        "description": "Roma tomatoes, freshly harvested",
        "location": "Nairobi, Kenya",
        "image_url": "https://example.com/image.jpg",
        "status": "active",
        "created_at": "2026-05-02T10:30:00Z"
    }
]
```

---

### Get Listing by ID
`GET /api/listings/{id}`

**Required:** JWT token (any authenticated user)

**Success Response — `200 OK`:**
```json
{
    "id": "990e8400-e29b-41d4-a716-446655440004",
    "farmer_id": "550e8400-e29b-41d4-a716-446655440000",
    "farmer_name": "John Kamau",
    "farmer_phone": "+254712345678",
    "farmer_email": "john@mavuno.com",
    "product_name": "Fresh Tomatoes",
    "quantity": 100,
    "unit": "kg",
    "price": 50,
    "description": "Roma tomatoes, freshly harvested",
    "location": "Nairobi, Kenya",
    "image_url": "https://example.com/image.jpg",
    "status": "active",
    "created_at": "2026-05-02T10:30:00Z"
}
```

---

### Update Listing (Farmer only)
`PUT /api/listings/{id}`

**Required:** JWT token + farmer role

**Request Body:**
```json
{
    "quantity": 80,
    "price": 55,
    "status": "sold_out"
}
```

**Success Response — `200 OK`:**
```json
{
    "id": "990e8400-e29b-41d4-a716-446655440004",
    "farmer_id": "550e8400-e29b-41d4-a716-446655440000",
    "product_name": "Fresh Tomatoes",
    "quantity": 80,
    "unit": "kg",
    "price": 55,
    "description": "Roma tomatoes, freshly harvested",
    "location": "Nairobi, Kenya",
    "image_url": "https://example.com/image.jpg",
    "status": "sold_out",
    "updated_at": "2026-05-02T11:00:00Z"
}
```

---

### Delete Listing (Farmer only)
`DELETE /api/listings/{id}`

**Required:** JWT token + farmer role

**Success Response — `204 No Content`**

---

## Marketplace - Messages

### Send Message
`POST /api/listings/{id}/messages`

**Required:** JWT token (any authenticated user)

**Request Body:**
```json
{
    "message": "Is this still available? Can you deliver to my location?"
}
```

**Success Response — `201 Created`:**
```json
{
    "id": "aaa0e8400-e29b-41d4-a716-446655440005",
    "listing_id": "990e8400-e29b-41d4-a716-446655440004",
    "sender_id": "550e8400-e29b-41d4-a716-446655440001",
    "sender_name": "Jane Buyer",
    "message": "Is this still available? Can you deliver to my location?",
    "created_at": "2026-05-02T10:45:00Z"
}
```

---

### Get Messages
`GET /api/listings/{id}/messages`

**Required:** JWT token (any authenticated user)

**Query Parameters (optional):**
- `limit` — Number of messages to return (default: 50)
- `offset` — Number of messages to skip (default: 0)

**Success Response — `200 OK`:**
```json
[
    {
        "id": "aaa0e8400-e29b-41d4-a716-446655440005",
        "listing_id": "990e8400-e29b-41d4-a716-446655440004",
        "sender_id": "550e8400-e29b-41d4-a716-446655440001",
        "sender_name": "Jane Buyer",
        "message": "Is this still available? Can you deliver to my location?",
        "created_at": "2026-05-02T10:45:00Z"
    },
    {
        "id": "aaa0e8400-e29b-41d4-a716-446655440006",
        "listing_id": "990e8400-e29b-41d4-a716-446655440004",
        "sender_id": "550e8400-e29b-41d4-a716-446655440000",
        "sender_name": "John Kamau",
        "message": "Yes, available. I can deliver this Friday.",
        "created_at": "2026-05-02T11:10:00Z"
    }
]
```

---

## Error Handling

All error responses follow this format:

```json
{
    "error": "Error message here",
    "status": 400
}
```

Common HTTP Status Codes:
- `200` — Success
- `201` — Created
- `204` — No Content
- `400` — Bad Request
- `401` — Unauthorized
- `403` — Forbidden
- `404` — Not Found
- `409` — Conflict
- `500` — Internal Server Error

---

## Health Check

### Health Status
`GET /health`

No token required.

**Success Response — `200 OK`:**
```json
{
    "status": "ok"
}
```

---

## User Roles

- **Farmer:** Can create and manage products, entries, supply agreements, analytics, reports, and marketplace listings. Can browse other listings and send messages.
- **Buyer:** Can browse marketplace listings and send messages to farmers. Cannot access farmer-specific endpoints (products, analytics, reports).
