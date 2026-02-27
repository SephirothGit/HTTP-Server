package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type RouterDeps struct {
	OrderHandler http.HandlerFunc
}

func NewRouter(deps RouterDeps) http.Handler {

	r := chi.NewRouter()
	limiter := NewRateLimiter(10, time.Minute)

	// Middleware stack
	r.Use(RequestIDMiddleware)
	r.Use(LoggingMiddleware)
	r.Use(RecoveryMiddleware)
	r.Use(RequestSizeLimitMiddleware(1024 * 1024))
	r.Use(TimeoutMiddleware(5 * time.Second))
	r.Use(MetricsMiddleware)
	r.Use(limiter.Middleware)

	// System routes
	r.Get("/health", healthHandler)
	r.Get("/ready", readinessHandler)
	r.Handle("/metrics", MetricsHandler())

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		r.Put("/orders/{id}", deps.OrderHandler)
	})

	// JWT
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(JWTAuthMiddleware("supersecretkey"))
		r.Put("/orders/{id}", deps.OrderHandler)
	})

	// Custom error 404
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		writeJSONError(w, http.StatusNotFound, "route not found")
	})

	// Custom error 405
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	})

	return r
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func readinessHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

func writeJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}
