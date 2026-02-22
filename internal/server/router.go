package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"encoding/json"
	"time"
	"context"
)

type RouterDeps struct {
	OrderHandler http.HandlerFunc
}

func NewRouter(deps RouterDeps) http.Handler {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(LoggingMiddleware)
	r.Use(TimeoutMiddleware(5 * time.Second))

	// System routes
	r.Get("/health", healthHandler)
	r.Get("/ready", readinessHandler)

	// API v1
	r.Route("/api/v1", func(r chi.Router){
		r.Put("/orders/{id}", deps.OrderHandler)
	})

	// Custom error 404
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request){
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	})

	return r
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))

	func readinessHandler(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	}
}

writeJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}

func TimeoutMiddleware(timeout time.duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r * http.Request) {
			ctx, cancel := context.WithTiemout(r.Context(), timeout)
			defer cancel
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}