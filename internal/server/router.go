package server

import "net/http"

type RouterDeps struct {
	OrderHandler http.HandlerFunc
}

func NewRouter(deps RouterDeps) http.Handler {
	mux := http.NewServeMux()

	// API v1
	v1 := http.NewServeMux()
	v1.Handle("/orders/", deps.OrderHandler)
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	// System routes
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/ready", readinessHandler)

	return mux
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func readinessHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}