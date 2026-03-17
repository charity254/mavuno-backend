# Mavuno API Documentation

Base URL: `http://localhost:8080` (development)

All protected routes require a JWT token in the Authorization header.
The token is obtained from the Login endpoint and expires after 72hours.

**How to use:**
1. Call `POST /api/auth/login` to get your token
2. Copy the token value from the response
3. Attach it to every protected request like this:
```
Authorization: Bearer <your_token_here>
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
    "token": ""
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

## General Error Responses

These can be returned by any protected route:

| Status | Message |
|--------|---------|
| `401` | `authorization header is required` |
| `401` | `invalid or expired token` |
| `403` | `you do not have permission to access this resource` |