# Order Service

HTTP backend на Go для управления заказами.

Проект для практики:
- layered architecture (handler / service / repository / domain)
- REST handlers
- middleware
- context usage
- error mapping

```bash
go mod tidy
go run .

Сервер стартует на :8080

Endpoints:

Update order status
PUT /orders/{id}

Responses:

204 — success
400 — invalid input
404 — order not found
409 — invalid status transition
500 — internal error

Architecture:

handler   → HTTP layer
service   → business logic
repository→ data storage
domain    → entities + rules
server    → router + middleware + http.Server