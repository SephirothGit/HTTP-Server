package server

import "net/http"

type RouterDeps struct {
	OrderHandler http.HandlerFunc
}

func NewRouter(deps RouterDeps) http.Handler {
	mux := http.NewServeMux()

// API v1
 v1 := http.NewServeMux()
 registerOrderRoutes(v1, deps)

 // Versioning
 mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

 // System routes
 mux.HandleFunc("/health", healthHandler)

 return mux
}

func registerOrderRoutes(mux *http.ServeMux, deps RouterDeps) {
	mux.HandleFunc("/orders/", deps.OrderHandler)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}