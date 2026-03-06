# Backend Order Service

Production-style Go backend service implementing order status management with clean architecture principles.

## Features

* REST API
* Clean Architecture (handler → service → domain → repository)
* Graceful shutdown
* JWT authentication
* Request logging (zap)
* Request ID middleware
* Prometheus metrics
* Rate limiting
* Timeout handling
* Request size middleware
* Custom error responses
* Health & readiness probes
* API versioning
* Optimistic locking (Version field)
* Domain events
* Unit tests

---

## Architecture

```
Client
   ↓
Router (chi)
   ↓
Middleware
   ↓
Handler
   ↓
Service
   ↓
Domain
   ↓
Repository
```

Project structure:

```
cmd/
   main.go

internal/
   domain/
   handler/
   repository/
   server/
   service/
```

---

## Run

### 1. Install dependencies

```
go mod tidy
```

### 2. Run server

```
go run ./cmd
```

Server starts on:

```
http://localhost:8080
```

---

## Tests

Unit tests cover the domain business logic.

Run tests:

```
go test./
```

Domain tests verify:

- Status transition rules
- Idempotent updates
- Version increment
- Domain events creation
- Event clearing behavior

Example:

```
go test ./internal/domain -v
```

## API

### Update order status

```
PUT /api/v1/orders/{id}
```

Body:

```
{
  "status": "paid"
}
```

Headers:

```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

Example:

```
curl -X PUT http://localhost:8080/api/v1/orders/123 \
-H "Authorization: Bearer TOKEN" \
-H "Content-Type: application/json" \
-d '{"status":"paid"}'
```

---

## JWT Token

Example token generation:

```
go run ./cmd/token
```

Or generate manually:

```
user_id = 123
secret = supersecretkey
```

---

## Health Checks

### Liveness probe

```
GET /health
```

Response:

```
ok
```

### Readiness probe

```
GET /ready
```

Response:

```
ready
```

---

## Middleware

### Request Logging

Logs include:

* request_id
* method
* path
* status
* duration

Example:

```
INFO request completed
method=PUT
path=/api/v1/orders/123
status=204
duration=5ms
request_id=abc123
```

---

## Metrics

Prometheus metrics endpoint:

GET /metrics

Example metrics:

http_requests_total

http_errors_total

http_request_duration_seconds

Used for monitoring:

- request count
- error rate
- request latency

---

### Timeout

Default timeout:

```
5 seconds
```

Long requests return:

```
504 Gateway Timeout
```

---

### Rate Limiting

Default:

```
10 requests per minute per IP
```

Response:

```
429 Too Many Requests
```

---

### Authentication

JWT authentication required for API endpoints.

Header:

```
Authorization: Bearer <token>
```

Response on failure:

```
401 Unauthorized
```

---

## Error Model

Example error:

```
{
  "error": "order not found"
}
```

HTTP Codes:

| Code | Meaning            |
| ---- | ------------------ |
| 400  | Bad request        |
| 401  | Unauthorized       |
| 404  | Not found          |
| 405  | Method not allowed |
| 409  | Conflict           |
| 429  | Rate limit         |
| 500  | Internal error     |
| 504  | Timeout            |

---

## Domain Model

Order:

```
ID
Status
Version
```

Statuses:

```
created
paid
shipped
canceled
```

Valid transitions:

```
created → paid
created → canceled
paid → shipped
paid → canceled
```

---

## Graceful Shutdown

The server handles:

* SIGINT
* SIGTERM

Shutdown timeout:

```
5 seconds
```

---

## Production Practices

This project demonstrates:

* Dependency Injection
* Domain-driven design basics
* Middleware architecture
* Structured logging
* Context propagation
* API versioning
* Concurrency-safe repository
* Idempotent operations

---

## Future Improvements

* PostgreSQL repository
* Transactions
* Outbox pattern
* OpenAPI documentation
* Docker support
* Integration tests

---

## Author

SephirothGit
