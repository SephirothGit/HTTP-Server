package server

import "net/http"

// Returns error 504
func Timeout504Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		done := make(chan struct{})

		go func() {
			next.ServeHTTP(w, r)
			close(done)
		}()

		select {
		case <-r.Context().Done():
			w.WriteHeader(http.StatusGatewayTimeout)
			w.Write([]byte("gateway timeout"))
		case <-done:
			// All good
		}
	})
}
