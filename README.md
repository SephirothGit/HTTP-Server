# Backend Service (Order Management)

Simple HTTP service for managing order status with clean architecture principles.

## Features

- Layered architecture (domain / service / repository / handler)
- In-memory repository implementation
- Status transition validation
- Graceful shutdown
- Logging middleware
- Clean and extensible structure

## Architecture

cmd/app – application entrypoint
internal/domain – business logic & rules
internal/service – use cases
internal/repository – data persistence
internal/handler – HTTP transport layer
internal/server – server, router, middleware


## Getting Started

### Run locally

```bash
go run ./cmd/app
```

Server starts on:
http://localhost:8080

## API

Update order status
PUT /orders/{id}
Body:

{
  "status": "paid"
}
## Possible statuses:

created

paid

shipped

canceled

## Design Principles

Separation of concerns

Dependency inversion

Explicit interfaces

Domain-driven mindset