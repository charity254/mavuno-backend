# Mavuno API Documentation

Base URL (development): `http://localhost:8080`
Base URL (production): `https://mavuno-backend.onrender.com`

---

All protected routes require a JWT token in the Authorization header.
The token is obtained from the Login endpoint and expires after 72 hours.

**How to use:**
1. Call `POST /api/auth/login` to get your token
2. Copy the token value from the response
3. In every protected request add a header called `Authorization` with the value `Bearer` followed by your token

**Note:** A new token is generated on every login. Store it on the frontend and reuse it until it expires.

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

## Authentication

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
| `400` | `email is required` |
| `400` | `full name is required` |
| `409` | `email already exists` |

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
    "token": "<your_token_here>"
}
```

**Notes:**
- Store this token on the frontend
- Attach it to every future request in the Authorization header
- Token expires after 72 hours — user must login again

**Error Responses:**
| Status | Message |
|--------|---------|
| `401` | `invalid email or password` |

---

## Products
Farmer only.

### Create Product
`POST /api/products`

**Request Body:**
```json
{
    "name": "Eggs",
    "unit_type": "crate",
    "description": "Fresh farm eggs"
}
```

**Success Response — `201 Created`:**
```json
{
    "id": "uuid",
    "farmer_id": "uuid",
    "name": "Eggs",
    "unit_type": "crate",
    "description": "Fresh farm eggs",
    "version": 1,
    "created_at": "2026-04-09T00:00:00Z",
    "updated_at": "2026-04-09T00:00:00Z"
}
```

---

### Get All Products
`GET /api/products`

Returns all active products belonging to the authenticated farmer.

**Success Response — `200 OK`:**
```json
[
    {
        "id": "uuid",
        "farmer_id": "uuid",
        "name": "Eggs",
        "unit_type": "crate",
        "description": "Fresh farm eggs",
        "version": 1,
        "created_at": "2026-04-09T00:00:00Z",
        "updated_at": "2026-04-09T00:00:00Z"
    }
]
```

---

### Get Single Product
`GET /api/products/{id}`

**Success Response — `200 OK`**

**Error Responses:**
| Status | Message |
|--------|---------|
| `404` | `product not found` |
| `403` | `you do not have permission to view this product` |

---

### Update Product
`PUT /api/products/{id}`

**Request Body:**
```json
{
    "name": "Eggs",
    "unit_type": "tray",
    "description": "Fresh farm eggs",
    "version": 1
}
```

**Notes:**
- `version` must match the current version in the database
- Version increments by 1 on every successful update

**Success Response — `200 OK`**

**Error Responses:**
| Status | Message |
|--------|---------|
| `409` | `conflict: product was updated by another session` |
| `403` | `you do not have permission to update this product` |

---

### Delete Product
`DELETE /api/products/{id}`

**Success Response — `204 No Content`**

---

## Produce Entries
Farmer only.

### Create Entry
`POST /api/entries`

**Request Body:**
```json
{
    "product_id": "uuid",
    "entry_date": "2026-04-09",
    "opening_stock": 10,
    "added_stock": 5,
    "sold_quantity": 8,
    "rejected_quantity": 1,
    "price_per_unit": 5000,
    "notes": "Good day for sales"
}
```

**Notes:**
- `entry_date` format: `YYYY-MM-DD`
- `price_per_unit` is in KES cents — e.g. `5000` = KES 50.00
- `sold_quantity + rejected_quantity` cannot exceed `opening_stock + added_stock`
- All values must be non negative

**Success Response — `201 Created`:**
```json
{
    "id": "uuid",
    "farmer_id": "uuid",
    "product_id": "uuid",
    "entry_date": "2026-04-09T00:00:00Z",
    "opening_stock": 10,
    "added_stock": 5,
    "sold_quantity": 8,
    "rejected_quantity": 1,
    "price_per_unit": 5000,
    "notes": "Good day for sales",
    "version": 1,
    "created_at": "2026-04-09T00:00:00Z",
    "updated_at": "2026-04-09T00:00:00Z",
    "total_available": 15,
    "remaining_stock": 6,
    "revenue_generated": 40000
}
```

**Error Responses:**
| Status | Message |
|--------|---------|
| `400` | `sold and rejected quantities cannot exceed total available stock of X` |
| `400` | `product not found` |
| `403` | `you do not have permission to use this product` |

---

### Get All Entries
`GET /api/entries`

**Optional query params:**
- `?start=2026-04-01` — filter by start date
- `?end=2026-04-30` — filter by end date
- `?product_id=uuid` — filter by product

**Success Response — `200 OK`**

---

### Get Single Entry
`GET /api/entries/{id}`

**Success Response — `200 OK`**

**Error Responses:**
| Status | Message |
|--------|---------|
| `404` | `entry not found` |

---

### Update Entry
`PUT /api/entries/{id}`

**Request Body:**
```json
{
    "product_id": "uuid",
    "entry_date": "2026-04-09",
    "opening_stock": 10,
    "added_stock": 5,
    "sold_quantity": 10,
    "rejected_quantity": 1,
    "price_per_unit": 5000,
    "notes": "Updated entry",
    "version": 1
}
```

**Error Responses:**
| Status | Message |
|--------|---------|
| `409` | `conflict: entry was updated by another session` |

---

### Delete Entry
`DELETE /api/entries/{id}`

**Success Response — `204 No Content`**

---

## Supply Locations
Farmer only.

### Create Supply Location
`POST /api/supply-locations`

**Request Body:**
```json
{
    "name": "Savannah Hotel",
    "contact_person": "Jane Wanjiku",
    "phone_number": "0712345678",
    "location_address": "Ngong Road, Nairobi",
    "notes": "Deliver before 8am"
}
```

**Notes:**
- `name`, `contact_person`, `phone_number`, `location_address` are all required

**Success Response — `201 Created`:**
```json
{
    "id": "uuid",
    "farmer_id": "uuid",
    "name": "Savannah Hotel",
    "contact_person": "Jane Wanjiku",
    "phone_number": "0712345678",
    "location_address": "Ngong Road, Nairobi",
    "notes": "Deliver before 8am",
    "version": 1,
    "created_at": "2026-04-09T00:00:00Z",
    "updated_at": "2026-04-09T00:00:00Z"
}
```

---

### Get All Supply Locations
`GET /api/supply-locations`

**Success Response — `200 OK`**

---

### Get Single Supply Location
`GET /api/supply-locations/{id}`

**Success Response — `200 OK`**

---

### Update Supply Location
`PUT /api/supply-locations/{id}`

**Request Body:**
```json
{
    "name": "Savannah Hotel",
    "contact_person": "Jane Wanjiku",
    "phone_number": "0712345678",
    "location_address": "Ngong Road, Nairobi",
    "notes": "Deliver before 7am",
    "version": 1
}
```

**Error Responses:**
| Status | Message |
|--------|---------|
| `409` | `conflict: supply location was updated by another session` |

---

### Delete Supply Location
`DELETE /api/supply-locations/{id}`

**Success Response — `204 No Content`**

---

## Supply Agreements
Farmer only.

### Create Supply Agreement
`POST /api/supply-agreements`

**Request Body:**
```json
{
    "product_id": "uuid",
    "supply_location_id": "uuid",
    "quantity_per_delivery": 10,
    "price_per_unit": 5000,
    "delivery_days": ["Monday", "Wednesday", "Friday"],
    "delivery_notes": "Deliver before 8am"
}
```

**Notes:**
- `delivery_days` must contain valid weekday names e.g. `Monday`, `Tuesday`
- At least one delivery day is required
- `price_per_unit` is in KES cents
- `product_id` must belong to the authenticated farmer
- `supply_location_id` must belong to the authenticated farmer

**Success Response — `201 Created`:**
```json
{
    "id": "uuid",
    "farmer_id": "uuid",
    "product_id": "uuid",
    "supply_location_id": "uuid",
    "quantity_per_delivery": 10,
    "price_per_unit": 5000,
    "delivery_days": ["Monday", "Wednesday", "Friday"],
    "delivery_notes": "Deliver before 8am",
    "active": true,
    "version": 1,
    "created_at": "2026-04-09T00:00:00Z",
    "updated_at": "2026-04-09T00:00:00Z"
}
```

---

### Get All Supply Agreements
`GET /api/supply-agreements`

**Success Response — `200 OK`**

---

### Get Active Supply Agreements
`GET /api/supply-agreements/active`

Returns only agreements where `active = true`.
Used by the frontend for dashboard reminders.

**Success Response — `200 OK`**

---

### Get Single Supply Agreement
`GET /api/supply-agreements/{id}`

**Success Response — `200 OK`**

---

### Update Supply Agreement
`PUT /api/supply-agreements/{id}`

**Request Body:**
```json
{
    "product_id": "uuid",
    "supply_location_id": "uuid",
    "quantity_per_delivery": 15,
    "price_per_unit": 5000,
    "delivery_days": ["Monday", "Wednesday", "Friday"],
    "delivery_notes": "Deliver before 7am",
    "active": true,
    "version": 1
}
```

**Error Responses:**
| Status | Message |
|--------|---------|
| `409` | `conflict: supply agreement was updated by another session` |

---

### Delete Supply Agreement
`DELETE /api/supply-agreements/{id}`

**Success Response — `204 No Content`**

---

## Analytics
Farmer only.

### Today's Summary
`GET /api/analytics/today`

Returns aggregated totals for today across all products.

**Success Response — `200 OK`:**
```json
{
    "total_available": 25,
    "total_sold": 10,
    "total_revenue": 50000,
    "total_rejected": 2
}
```

**Note:** `total_revenue` is in KES cents.

---

### Revenue Trend
`GET /api/analytics/revenue?start=2026-04-01&end=2026-04-30`

Returns daily revenue totals for the date range.

**Success Response — `200 OK`:**
```json
[
    {
        "date": "2026-04-09",
        "revenue": 50000
    }
]
```

---

### Stock Trend
`GET /api/analytics/stock?start=2026-04-01&end=2026-04-30&product_id=uuid`

Returns daily remaining stock for a specific product.

**Required query params:** `start`, `end`, `product_id`

**Success Response — `200 OK`:**
```json
[
    {
        "date": "2026-04-09",
        "remaining_stock": 13
    }
]
```

---

### Rejection Trend
`GET /api/analytics/rejected?start=2026-04-01&end=2026-04-30&product_id=uuid`

Returns daily rejected quantities for a specific product.

**Required query params:** `start`, `end`, `product_id`

**Success Response — `200 OK`:**
```json
[
    {
        "date": "2026-04-09",
        "rejected_quantity": 2
    }
]
```

---

### Product Summary
`GET /api/analytics/product-summary?start=2026-04-01&end=2026-04-30`

Returns aggregated totals per product.

**Required query params:** `start`, `end`

**Success Response — `200 OK`:**
```json
[
    {
        "product_name": "Milk",
        "total_sold": 10,
        "total_rejected": 2,
        "total_revenue": 50000
    }
]
```

---

### Planned vs Actual
`GET /api/analytics/planned-vs-actual?start=2026-04-01&end=2026-04-30&product_id=uuid`

Compares agreed delivery quantities against actual sold quantities per day.

**Required query params:** `start`, `end`, `product_id`

**Success Response — `200 OK`:**
```json
[
    {
        "date": "2026-04-09",
        "planned_quantity": 10,
        "actual_sold": 8,
        "difference": -2
    }
]
```

---

## Reports
Farmer only.

### View Report as JSON
`GET /api/reports/export?start=2026-04-01&end=2026-04-30`

Optional: `&product_id=uuid`

Returns structured report data for displaying in the web app.

**Success Response — `200 OK`:**
```json
{
    "rows": [
        {
            "date": "2026-04-09",
            "product_name": "Tomatoes",
            "opening_stock": 50,
            "added_stock": 20,
            "total_available": 70,
            "sold": 30,
            "rejected": 5,
            "remaining": 35,
            "price_per_unit_kes": 80.00,
            "revenue_kes": 2400.00
        }
    ],
    "total_sold": 30,
    "total_rejected": 5,
    "total_revenue_kes": 2400.00
}
```

---

### Download CSV Report
`GET /api/reports/download?start=2026-04-01&end=2026-04-30`

Optional: `&product_id=uuid`

Returns a CSV file for downloading.

**CSV Headers:**

**Success Response — `200 OK` (CSV file)**

---

## Marketplace

### Create Listing
`POST /api/listings` — Farmer only

**Request Body:**
```json
{
    "title": "Fresh Tomatoes Available",
    "description": "Freshly harvested tomatoes from our farm",
    "price": 8000,
    "quantity": 50,
    "unit_type": "kg",
    "location": "Ngong Road, Nairobi"
}
```

**Notes:**
- `price` is in KES cents
- `title`, `unit_type`, `location` are required

**Success Response — `201 Created`:**
```json
{
    "id": "uuid",
    "farmer_id": "uuid",
    "title": "Fresh Tomatoes Available",
    "description": "Freshly harvested tomatoes from our farm",
    "price": 8000,
    "quantity": 50,
    "unit_type": "kg",
    "location": "Ngong Road, Nairobi",
    "version": 1,
    "created_at": "2026-05-02T00:00:00Z",
    "updated_at": "2026-05-02T00:00:00Z"
}
```

---

### Get All Listings
`GET /api/listings` — All authenticated users

**Success Response — `200 OK`**

---

### Get Single Listing
`GET /api/listings/{id}` — All authenticated users

**Success Response — `200 OK`**

---

### Update Listing
`PUT /api/listings/{id}` — Farmer only

**Request Body:**
```json
{
    "title": "Fresh Tomatoes Available",
    "description": "Freshly harvested tomatoes",
    "price": 8000,
    "quantity": 40,
    "unit_type": "kg",
    "location": "Ngong Road, Nairobi",
    "version": 1
}
```

**Error Responses:**
| Status | Message |
|--------|---------|
| `409` | `conflict: listing was updated by another session` |
| `403` | `you do not have permission to update this listing` |

---

### Delete Listing
`DELETE /api/listings/{id}` — Farmer only

**Success Response — `204 No Content`**

---

### Send Message
`POST /api/listings/{id}/messages` — All authenticated users

**Request Body:**
```json
{
    "content": "I am interested in buying 10kg of tomatoes. Are they still available?"
}
```

**Notes:**
- `content` cannot be empty
- `content` cannot exceed 1000 characters
- Farmers cannot message themselves on their own listing
- The receiver is automatically set to the farmer who owns the listing

**Success Response — `201 Created`:**
```json
{
    "id": "uuid",
    "listing_id": "uuid",
    "sender_id": "uuid",
    "receiver_id": "uuid",
    "content": "I am interested in buying 10kg of tomatoes. Are they still available?",
    "created_at": "2026-05-02T00:00:00Z"
}
```

---

### Get Messages
`GET /api/listings/{id}/messages` — Participants only

Only the farmer who owns the listing and the buyer who sent the first message can read the conversation.

**Success Response — `200 OK`:**
```json
[
    {
        "id": "uuid",
        "listing_id": "uuid",
        "sender_id": "uuid",
        "receiver_id": "uuid",
        "content": "I am interested in buying 10kg of tomatoes.",
        "created_at": "2026-05-02T00:00:00Z"
    }
]
```

**Error Responses:**
| Status | Message |
|--------|---------|
| `403` | `you do not have access to these messages` |

---

## Articles

### Create Article
`POST /api/articles` — Farmer only

**Request Body:**
```json
{
    "title": "How to prevent tomato blight",
    "content": "Tomato blight is a common disease...",
    "category": "Disease Alerts"
}
```

**Notes:**
- `title`, `content`, `category` are all required

**Success Response — `201 Created`:**
```json
{
    "id": "uuid",
    "author_id": "uuid",
    "title": "How to prevent tomato blight",
    "content": "Tomato blight is a common disease...",
    "category": "Disease Alerts",
    "created_at": "2026-05-11T00:00:00Z",
    "updated_at": "2026-05-11T00:00:00Z"
}
```

---

### Get All Articles
`GET /api/articles` — All authenticated users

Optional: `?category=Disease Alerts`

**Success Response — `200 OK`**

---

### Get Single Article
`GET /api/articles/{id}` — All authenticated users

**Success Response — `200 OK`**

---

### Delete Article
`DELETE /api/articles/{id}` — Author only

**Success Response — `204 No Content`**

**Error Responses:**
| Status | Message |
|--------|---------|
| `403` | `you do not have permission to delete this article` |

---

## General Error Responses

These can be returned by any protected route:

| Status | Message |
|--------|---------|
| `400` | `invalid request body` |
| `401` | `authorization header is required` |
| `401` | `invalid or expired token` |
| `403` | `you do not have permission to access this resource` |
| `404` | `record not found` |
| `409` | `conflict: record was updated by another session` |
| `429` | `too many requests — please try again later` |
| `204` | successful delete — no body returned |

---

## User Roles

| Role | Access |
|------|--------|
| **Farmer** | Products, entries, supply locations, supply agreements, analytics, reports, marketplace listings, articles |
| **Buyer** | Browse marketplace listings, send and read messages, read articles |