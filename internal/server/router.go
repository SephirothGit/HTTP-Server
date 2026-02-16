package server

import (
	"context"
	"net/http"
	"time"
)

type RouterDeps struct {
	OrderHandler http.HandlerFunc
}

func NewRouter(deps RouterDeps) http.Handler {
	mux := http.NewServeMux()

	// API v1
	v1 := http.NewServeMux()
	v1.Handle("/orders/", applyContext(deps.OrderHandler))
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1/orders/", v1))

	// System routes
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/ready", readinessHandler)

	return mux
}

// Makes context with timeout for each request
func applyContext(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func readinessHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}