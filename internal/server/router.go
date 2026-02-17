package server

import (
	"encoding/json"
	"net/http"
)

type RouterDeps struct {
	OrderHandler http.HandlerFunc
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func NewRouter(deps RouterDeps) http.Handler {
	mux := http.NewServeMux()

	v1 := http.NewServeMux()

	// Orders route with method guard
	v1.Handle("/orders/", methodGuard(
		http.MethodPut,
		deps.OrderHandler,
	))

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	// System routes
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/ready", readinessHandler)

	return notFound(mux)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func readinessHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

func notFound(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r * http.Request){
		rw := &statusRecorder{ResponseWriter: w, status: 200}

		next.ServeHTTP(rw, r)

		if rw.status == http.StatusNotFound {
			writeJSON404(w)
		}
	})
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func writeJSON404(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": "route not found",
	})
}

func methodGuard(method string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)

			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "method not allowed",
			})
			return
		}

		h.ServeHTTP(w, r)
	}
}